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

var specs = []Spec{}

type Vars map[string]string

type Service interface {
	// Init initializes the service with the given Service Options.
	Init(route string, vars Vars, opts ...Option) error
	// Parse parses the route.
	Send(title, body string) error
	// SetOption sets the option for the service.
	SetOption(key string, value interface{})
}

type Spec struct {
	Template []string
	Init     func() Service
}

func Specs() []Spec {
	return specs
}

type Option func(Service)

func SetTLS(isTLS bool) Option {
	return func(s Service) {
		s.SetOption("isTLS", isTLS)
	}
}
