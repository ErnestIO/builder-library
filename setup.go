/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package scheduler

import (
	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
	"gopkg.in/redis.v3"
)

func (s *Scheduler) Setup(n string) {
	s.configSetup(n)
	s.N = s.natsClient()
	s.R = s.redisClient()
}

func (s *Scheduler) configSetup(n string) {
	s.C = ecc.NewConfig(n)
}

func (s *Scheduler) natsClient() *nats.Conn {
	if s.C == nil {
		s.configSetup(s.natsURI)
	}
	return s.C.Nats()
}

func (s *Scheduler) redisClient() *redis.Client {
	if s.C == nil {
		s.configSetup(s.natsURI)
	}
	return s.C.Redis()
}
