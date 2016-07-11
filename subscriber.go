/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package scheduler

import (
	"encoding/json"

	"github.com/nats-io/nats"
	"gopkg.in/redis.v3"
)

type Scheduler struct {
	natsURI string
	R       *redis.Client
	N       *nats.Conn
}

type Error struct {
	Code    json.Number `json:"code,Number"`
	Message string      `json:"message"`
}

type response struct {
	Uuid      string           `json:"uuid"`
	Service   string           `json:"service"`
	Action    string           `json:"action"`
	Component *json.RawMessage `json:"component"`
	Error     Error            `json:"error"`
}

func (s *Scheduler) ProcessRequest(from string, to string) {
	s.N.Subscribe(from, func(m *nats.Msg) {
		s.manageRequest(m.Data, to)
	})
}

func (s *Scheduler) ProcessSuccessResponse(from string, nex string, done string) {
	s.N.Subscribe(from, func(m *nats.Msg) {
		s.manageSuccessResponse(m.Data, nex, done)
	})
}

func (s *Scheduler) ProcessFailedResponse(from string, to string) {
	s.N.Subscribe(from, func(m *nats.Msg) {
		s.manageFailedResponse(m.Data, to)
	})
}
