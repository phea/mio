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

type Service interface {
	// Send sends a message.
	Send(title, body string) error
	// Route returns the service route.
	Route() string
}

type Spec struct {
	Template []string
	Init     func() Service
}
