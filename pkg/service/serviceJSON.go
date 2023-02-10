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

package service

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var jsonTemplates = []string{
	"json://",
	"json://{host}:{port}",
	"json://{user}@{host}:{port}",
	"json://{user}:{password}@{host}",
	"json://{user}:{password}@{host}:{port}",
}

func init() {
	specs = append(specs, jsonSpec())
}

// check if ServiceJSON implements Service interface
var _ Service = (*ServiceJSON)(nil)

type ServiceJSON struct {
	method   string
	port     int
	scheme   string
	url      url.URL
	isTLS    bool
	rawRoute string
	vars     Vars // this should probably be a Field struct
}

func jsonSpec() Spec {
	return Spec{
		Template: jsonTemplates,
		Init: func() Service {
			return defaultJsonService()
		},
	}
}

func defaultJsonService() *ServiceJSON {
	return &ServiceJSON{
		method: "POST",
		scheme: "json",
		port:   80,
		isTLS:  false,
	}
}

// Init initializes the service
func (s *ServiceJSON) Init(route string, vars Vars, opts ...Option) error {
	for _, opt := range opts {
		opt(s)
	}
	s.rawRoute = route

	fmt.Println(vars)
	return nil
}

// Send sends a JSON message.
func (s *ServiceJSON) Send(title, body string) (*http.Response, error) {
	fmt.Printf("JSON: %s - %s\n", title, body)

	// client := http.DefaultClient
	// req, err := http.NewRequestWithContext(context.Background(),
	// 	s.method, s.URL(), nil)
	// if err != nil {
	// 	return nil, err
	// }

	// client.Do(req)
	fmt.Printf("Making request to %s, Method: %s\n", s.URL(), s.method)
	return &http.Response{}, nil
}

// URL returns the URL for the http request.
func (s *ServiceJSON) URL() string {
	url, err := url.Parse(s.rawRoute)
	if err != nil {
		return ""
	}

	// intialize a string builder
	var sb strings.Builder
	// if isTLS is true, use https, otherwise use http
	if s.isTLS {
		sb.WriteString("https://")
	} else {
		sb.WriteString("http://")
	}

	if url.User != nil {
		sb.WriteString(url.User.String())
		sb.WriteString("@")
	}

	sb.WriteString(url.Hostname())
	if s.port != 80 {
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%d", s.port))
	}
	sb.WriteString(url.Path)

	return sb.String()
}
