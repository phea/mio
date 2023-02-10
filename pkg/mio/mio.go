/*
 * BSD 3-Clause License
 *
 * Copyright (c) 2023, Phea Duch <phea.duch@gmail.com>
 * All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license
 * that can be found in the LICENSE file.
 *
 */

package mio

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/phea/mio/internal/matcher"
	"github.com/phea/mio/pkg/service"
)

var (
	ErrServiceNotFound = fmt.Errorf("route does not match any service")
)

type serviceMatch struct {
	service service.Service
	matcher *matcher.Matcher
}

var serviceMatchers = []serviceMatch{}

func init() {
	// Register all services here.
	for _, spec := range service.Specs() {
		addMatcher(spec)
	}
}

// addMatcher takes an array of service specs and adds
// them to the service matchers array.
func addMatcher(spec service.Spec) {
	for _, t := range spec.Template {
		serviceMatchers = append(serviceMatchers, serviceMatch{
			service: spec.Init(),
			matcher: matcher.New(t),
		})
	}
}

// Notifier is responsible for sending messages.
type Notifier struct {
	svcs []service.Service
}

// Add matches the route to a service and adds it to the notifier.
// NOTE: This function assumes templates are ordered least to most specific.
func (n *Notifier) Add(route string, opts ...service.Option) error {
	var svc service.Service
	for _, m := range serviceMatchers {
		if m.matcher.IsMatch(route) {
			vars := m.matcher.Vars(route)
			err := m.service.Init(route, vars, opts...)
			if err != nil {
				return err
			}

			svc = m.service
		}
	}
	if svc == nil {
		return ErrServiceNotFound
	}

	n.svcs = append(n.svcs, svc)
	return nil
}

// Must is a helper function that calls Add and panics if an error occurs.
func (n *Notifier) Must(route string, opts ...service.Option) {
	if err := n.Add(route, opts...); err != nil {
		panic(err)
	}
}

type broadCastResult struct {
	serviceName string
	resp        http.Response // could be change to a buffer or []byte slice
	err         error
}

// Broadcast asynchronously sends a message to all registered services.
func (n *Notifier) Broadcast(title, body string) {
	// create a channel to receive results from goroutines
	results := make(chan broadCastResult, len(n.svcs))
	// create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	// add the number of goroutines to wait for
	wg.Add(len(n.svcs))
	// iterate over all registered services
	for _, svc := range n.svcs {
		// create a goroutine for each service
		go func(svc service.Service) {
			// send the result to the channel
			resp, err := svc.Send(title, body)
			results <- broadCastResult{
				// serviceName: svc.Route(),
				resp: *resp,
				err:  err,
			}
			// decrement the wait group
			wg.Done()
		}(svc)
	}

	// wait for all goroutines to finish
	wg.Wait()
	// close the channel
	close(results)

	// we can do something with the results here
	// for example, we can print the results
	successN, failN := 0, 0

	for result := range results {
		if result.err != nil {
			failN++
		} else {
			successN++
		}
	}

	log.Printf("Broadcast: %d success, %d failed\n", successN, failN)
}
