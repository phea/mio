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
	"log"
	"testing"

	"github.com/phea/mio/internal/matcher"
)

// TestGnomeTemplate tests the gnome template.
func TestGnomeTemplate(t *testing.T) {
	var gnomeMatchers []*matcher.Matcher
	for _, tmpl := range gnomeTemplates {
		gnomeMatchers = append(gnomeMatchers, matcher.New(tmpl))
	}

	tests := []struct {
		route string
		match bool
	}{
		{"gnome://", true},
		{"gnome:///", true},
		{"gnome://foo", true},
		{"gnome://foo/", true},
		{"gnome://foo/bar", true},
		{"gnome://?icon=foo", true},
	}

	for _, test := range tests {
		var matched bool
		for _, m := range gnomeMatchers {
			if m.IsMatch(test.route) {
				matched = true
				break
			}
		}

		if matched != test.match {
			t.Errorf("route %s should match: %t", test.route, test.match)
		}
	}
}

// TestGnomeInit tests the gnome service initialization.
func TestGnomeInit(t *testing.T) {
	svc := defaultGnomeService()

	m := matcher.New("gnome://")
	mVars, err := m.Vars("gnome://?icon=foo")
	if err != nil {
		log.Fatalf("expected no error, got %v", err)
	}

	err = svc.Init("gnome://?icon=foo", mVars)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if svc.icon != "foo" {
		t.Errorf("expected icon to be foo, got %s", svc.icon)
	}
}

// TestGnomeSend tests the gnome service send method.
func TestGnomeSend(t *testing.T) {
	svc := defaultGnomeService()
	svc.sender = func(payload gnomePayload) error {
		return nil
	}

	err := svc.Send("test", "message")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
