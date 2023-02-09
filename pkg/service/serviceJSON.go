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

import "fmt"

// RegisterJSON returns the JSON service spec.
func RegisterJSON() Spec {
	return Spec{
		Template: []string{"json://"},
		Init: func() Service {
			return &ServiceJSON{}
		},
	}
}

// check if ServiceJSON implements Service interface
var _ Service = (*ServiceJSON)(nil)

type ServiceJSON struct{}

// Send sends a JSON message.
func (s *ServiceJSON) Send(title, body string) error {
	fmt.Printf("JSON: %s - %s", title, body)
	return nil
}

// Route returns the service route.
func (s *ServiceJSON) Route() string {
	return "json://"
}
