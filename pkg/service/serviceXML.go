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
)

var xmlTemplates = []string{
	"xml://",
	"xml://{host}:{port}",
	"xml://{user}@{host}:{port}",
	"xml://{user}:{password}@{host}",
	"xml://{user}:{password}@{host}:{port}",
}

func init() {
	specs = append(specs, xmlSpec())
}

// check if ServiceXML implements Service interface
var _ Service = (*ServiceXML)(nil)

type ServiceXML struct {
	method   string
	port     string
	scheme   string
	rawRoute string
	url      *url.URL
	isTLS    bool
	vars     Vars // this should probably be a Field struct
}

func xmlSpec() Spec {
	return Spec{
		Template: xmlTemplates,
		Init: func() Service {
			return defaultXMLService()
		},
	}
}

func defaultXMLService() *ServiceXML {
	return &ServiceXML{
		scheme: "xml",
		method: "POST",
		port:   "80",
		isTLS:  false,
	}
}

// Init initializes the service with the given Service Options.
func (s *ServiceXML) Init(route string, vars Vars, opts ...Option) error {
	for _, opt := range opts {
		opt(s)
	}
	s.rawRoute = route
	s.vars = vars
	return nil
}

// Parse parses the route and fills in the ServiceXML struct.
func (s *ServiceXML) Parse(route string) error {
	s.rawRoute = route
	u, err := url.Parse(route)
	if err != nil {
		return err
	}
	s.url = u
	s.scheme = u.Scheme
	if u.Port() != "" {
		s.port = u.Port()
	}
	if u.User != nil {
		s.vars["user"] = u.User.Username()
		s.vars["password"], _ = u.User.Password()
	}
	return nil
}

// Send sends a XML message.
func (s *ServiceXML) Send(title, body string) (*http.Response, error) {
	fmt.Printf("XML: %s - %s", title, body)
	return &http.Response{}, nil
}

// URL returns the URL for the service http request call.
func (s *ServiceXML) URL() string {
	return ""
}
