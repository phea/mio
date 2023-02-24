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
	"time"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
)

var gnomeTemplates = []string{
	"gnome://",
}

func init() {
	specs = append(specs, gnomeSpec())
}

// check if ServiceGnome implements Service interface
var _ Service = (*ServiceGnome)(nil)

type gnomePayload struct {
	appName string
	icon    string
	title   string
	body    string
}

type senderFunc func(payload gnomePayload) error

type ServiceGnome struct {
	scheme string
	icon   string
	vars   Vars
	sender senderFunc
}

func gnomeSpec() Spec {
	return Spec{
		Template: gnomeTemplates,
		Init: func() Service {
			return defaultGnomeService()
		},
	}
}

// defaultSender is the default sender function for ServiceGnome.
// it uses the github.com/esiqveland/notify package to send the notification.
func defaultSender(payload gnomePayload) error {
	conn, err := dbus.SessionBusPrivate()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = conn.Auth(nil); err != nil {
		return err
	}

	if err = conn.Hello(); err != nil {
		return err
	}
	notification := notify.Notification{
		AppName:       payload.appName,
		AppIcon:       payload.icon,
		Summary:       payload.title,
		Body:          payload.body,
		ExpireTimeout: time.Second * 5,
	}

	if _, err := notify.SendNotification(conn, notification); err != nil {
		return err
	}
	return nil
}

func defaultGnomeService() *ServiceGnome {
	return &ServiceGnome{
		scheme: "gnome",
		icon:   "dialog-information",
		sender: defaultSender,
	}
}

// Init initializes the service with the given Service Options.
func (s *ServiceGnome) Init(route string, vars Vars, opts ...Option) error {
	s.vars = vars
	for _, opt := range opts {
		opt(s)
	}

	if s.vars["icon"] != "" {
		s.icon = s.vars["icon"]
	}

	return nil
}

// Send sends the notification to the service.
func (s *ServiceGnome) Send(title, body string) error {
	// create payload
	payload := gnomePayload{
		appName: s.vars["name"],
		icon:    s.icon,
		title:   title,
		body:    body,
	}

	// send notification
	err := s.sender(payload)
	if err != nil {
		log.Printf("error sending notification: %v", err.Error())
	}

	return nil
}

// SetOption sets the option for the service.
func (s *ServiceGnome) SetOption(key string, value interface{}) {
	s.vars[key] = value.(string)
}
