/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package scheduler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats"
	"gopkg.in/redis.v3"
)

func (s *Scheduler) Setup(n string) {
	if n == "" {
		n = nats.DefaultURL
	}

	s.natsURI = n
	s.N = s.natsClient()
	s.R = s.redisClient()
}

func (s *Scheduler) natsClient() *nats.Conn {
	n, err := nats.Connect(s.natsURI)
	if err != nil {
		log.Println("Could not connect to NATS server")
		panic(err)
	}

	return n
}

func (s *Scheduler) redisClient() *redis.Client {
	var cfg struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int64  `json:"DB"`
	}

	resp, err := s.N.Request("config.get.redis", nil, time.Second)
	if err != nil {
		log.Println("could not load config")
		log.Panic(err)
	}

	err = json.Unmarshal(resp.Data, &cfg)
	if err != nil {
		log.Panic(err)
	}
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
