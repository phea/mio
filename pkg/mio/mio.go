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

import "github.com/phea/mio/pkg/service"

type route struct {
	schema string
	svc    service.Service
}

type register struct {
	svcs []route
}

// add takes a service definition and adds the service routes
// to the register.
func (r *register) add(s service.Spec) {
	for _, t := range s.Template {
		r.svcs = append(r.svcs, route{schema: t, svc: s.Init()})
	}
}

// reg is the global register for all services.
var reg = register{svcs: []route{}}

func init() {
	// Register all services here.
	reg.add(service.RegisterJSON())
}
