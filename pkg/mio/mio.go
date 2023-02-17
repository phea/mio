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
	err error
}

// Broadcast asynchronously sends a message to all registered services.
func (n *Notifier) Broadcast(title, body string) {
	results := make(chan broadCastResult, len(n.svcs))
	var wg sync.WaitGroup
	wg.Add(len(n.svcs))
	for _, svc := range n.svcs {
		go func(svc service.Service) {
			err := svc.Send(title, body)
			results <- broadCastResult{
				// serviceName: svc.Route(),
				err: err,
			}
			wg.Done()
		}(svc)
	}

	// wait for all goroutines to finish
	wg.Wait()
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
