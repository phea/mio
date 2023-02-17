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
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"net/url"
	"strings"
)

var xmlTemplates = []string{
	"xml://{host}",
	"xml://{host}:{port}",
	"xml://{user}@{host}",
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
	scheme   string
	rawRoute string
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
		isTLS:  true,
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

type xmlPayload struct {
	Title   string
	Body    string
	XMLName xml.Name `xml:"Payload"`
}

// Send sends a XML message.
func (s *ServiceXML) Send(title, body string) error {
	payload := &xmlPayload{Title: title, Body: body}
	data, err := xml.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(context.Background(),
		s.method, s.Endpoint(), reader)
	if err != nil {
		return err
	}

	// set xml headers
	req.Header.Set("Content-Type", "application/xml")
	_, err = client.Do(req)
	return err
}

// Endpoint returns the endpoint URL for the http request.
func (s *ServiceXML) Endpoint() string {
	url, err := url.Parse(s.rawRoute)
	if err != nil {
		return ""
	}

	var sb strings.Builder
	// if isTLS is true, use https, otherwise use http
	if s.isTLS {
		sb.WriteString("https://")
	} else {
		sb.WriteString("http://")
	}

	// if user is not empty, append it to the endpoint
	if url.User != nil {
		sb.WriteString(url.User.String())
		sb.WriteString("@")
	}

	sb.WriteString(url.Hostname())
	if url.Port() != "" {
		sb.WriteString(":")
		sb.WriteString(url.Port())
	}
	sb.WriteString(url.Path)

	return sb.String()
}

// SetOption sets the option for the service.
func (s *ServiceXML) SetOption(key string, value interface{}) {
	switch key {
	case "method":
		s.method = value.(string)
	case "scheme":
		s.scheme = value.(string)
	case "isTLS":
		s.isTLS = value.(bool)
	default:
		s.vars[key] = value.(string)
	}
}
