/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package scheduler

import (
	"errors"
	"os"
	"time"

	"github.com/nats-io/nats"
)

var s Scheduler

func waitMsg(ch chan *nats.Msg) (*nats.Msg, error) {
	select {
	case msg := <-ch:
		return msg, nil
	case <-time.After(time.Millisecond * 500):
	}
	return nil, errors.New("timeout")
}

func setup() {
	s := Scheduler{natsURI: os.Getenv("NATS_URI")}
	n := s.natsClient()
	n.Subscribe("config.get.redis", func(msg *nats.Msg) {
		n.Publish(msg.Reply, []byte(`{"DB":0,"addr":"localhost:6379","password":""}`))
	})
}
