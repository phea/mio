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

import "net/url"

var smtpTemplates = []string{
	"smtp://",
	"smtp://{user}:{pass}",
	"smtp://{user}:{pass}@{host}",
	"smtp://{user}:{pass}@{host}:{port}",
}

var supportedSMTPProviders = []string{
	"gmail.com",
	"outlook.com",
	"hotmail.com",
	"yahoo.com",
}

func init() {
	specs = append(specs, smtpSpec())
}

var _ Service = (*ServiceSMTP)(nil)

type ServiceSMTP struct {
	port     string
	scheme   string
	url      url.URL
	isTLS    bool
	rawRoute string
	vars     Vars
}

func smtpSpec() Spec {
	return Spec{
		Template: smtpTemplates,
		Init: func() Service {
			return defaultSMTPService()
		},
	}
}

func defaultSMTPService() *ServiceSMTP {
	return &ServiceSMTP{
		scheme: "smtp",
		port:   "587",
		isTLS:  true,
	}
}

// Init initializes the service
func (s *ServiceSMTP) Init(route string, vars Vars, opts ...Option) error {
	for _, opt := range opts {
		opt(s)
	}
	s.rawRoute = route
	if vars["port"] != "" {
		s.port = vars["port"]
	}

	return nil
}

// Send sends the email
func (s *ServiceSMTP) Send(title, body string) error {
	return nil
}

// SetOption sets options for the service.
func (s *ServiceSMTP) SetOption(key string, value interface{}) {
	switch key {
	case "isTLS":
		s.isTLS = value.(bool)
	}
}
