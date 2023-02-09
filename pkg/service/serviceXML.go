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

func RegisterXML() Spec {
	return Spec{
		Template: []string{"xml://"},
		Init: func() Service {
			return &ServiceXML{}
		},
	}
}

// check if ServiceXML implements Service interface
var _ Service = (*ServiceXML)(nil)

type ServiceXML struct{}

// Send sends a XML message.
func (s *ServiceXML) Send(title, body string) error {
	fmt.Printf("XML: %s - %s", title, body)
	return nil
}

// Route returns the service route.
func (s *ServiceXML) Route() string {
	return "xml://"
}
