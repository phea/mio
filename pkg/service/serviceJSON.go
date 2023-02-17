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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var jsonTemplates = []string{
	"json://{host}",
	"json://{host}:{port}",
	"json://{user}@{host}",
	"json://{user}@{host}:{port}",
	"json://{user}:{pass}@{host}",
	"json://{user}:{pass}@{host}:{port}",
}

func init() {
	specs = append(specs, jsonSpec())
}

// check if ServiceJSON implements Service interface
var _ Service = (*ServiceJSON)(nil)

type ServiceJSON struct {
	method   string
	scheme   string
	isTLS    bool
	rawRoute string
	vars     Vars // this should probably be a Field struct
}

func jsonSpec() Spec {
	return Spec{
		Template: jsonTemplates,
		Init: func() Service {
			return defaultJSONService()
		},
	}
}

func defaultJSONService() *ServiceJSON {
	return &ServiceJSON{
		method: "POST",
		scheme: "json",
		isTLS:  true,
	}
}

// Init initializes the service
func (s *ServiceJSON) Init(route string, vars Vars, opts ...Option) error {
	for _, opt := range opts {
		opt(s)
	}
	s.rawRoute = route
	s.vars = vars

	return nil
}

// Send sends a JSON message.
func (s *ServiceJSON) Send(title, body string) error {
	data, err := json.Marshal(map[string]string{
		"title": title,
		"body":  body,
	})

	if err != nil {
		return err
	}

	// create a io.Reader from data
	reader := bytes.NewReader(data)

	client := http.DefaultClient
	req, err := http.NewRequestWithContext(context.Background(),
		s.method, s.Endpoint(), reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	_, err = client.Do(req)
	return err
}

// Endpoint returns the endpoint URL for the http request.
func (s *ServiceJSON) Endpoint() string {
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
	if url.Port() != "" {
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%s", url.Port()))
	}
	sb.WriteString(url.Path)

	return sb.String()
}

// SetOption sets options for the service.
func (s *ServiceJSON) SetOption(key string, value interface{}) {
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
