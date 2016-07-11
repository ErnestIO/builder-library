/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package scheduler

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicProcessRequest(t *testing.T) {
	processCh := make(chan *nats.Msg, 2)
	setup()
	s.Setup(os.Getenv("NATS_URI"))
	Convey("Given the scheduler is set up to process request", t, func() {
		s.ProcessRequest("components.process", "component.process")
		sub, _ := s.N.ChanSubscribe("component.process", processCh)

		Convey("When it recieves a components.process message", func() {
			body := []byte(`{"service":"aaa", "components":[{"id":"kk"},{"id":"oo"}]}`)
			s.N.Publish("components.process", body)

			Convey("Then it should send a message for its component", func() {
				var res struct {
					ID      string `json:"id"`
					UUID    string `json:"_uuid"`
					BatchID string `json:"_batch_id"`
				}
				msg, err := waitMsg(processCh)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &res)
				So(res.ID, ShouldEqual, "kk")
				So(res.UUID, ShouldNotBeNil)

				msg2, err := waitMsg(processCh)
				So(err, ShouldBeNil)
				json.Unmarshal(msg2.Data, &res)
				So(res.ID, ShouldEqual, "oo")
				So(res.UUID, ShouldNotBeNil)

				stored, err := s.R.Get(res.BatchID).Result()
				req := request{}
				json.Unmarshal([]byte(stored), &req)
				So(len(req.Components), ShouldEqual, 2)
				So(req.Service, ShouldEqual, "aaa")
			})
		})
		sub.Unsubscribe()

		s.ProcessRequest("components.process", "component.process")
		sub, _ = s.N.ChanSubscribe("component.process", processCh)
		Convey("When its requested to be processed sequentially", func() {
			body := []byte(`{"service":"aaa", "components":[{"id":"kk"},{"id":"oo"}],"sequential_processing":true}`)
			s.N.Publish("components.process", body)

			Convey("Then it should only process the first component", func() {
				var res struct {
					ID      string `json:"id"`
					UUID    string `json:"_uuid"`
					BatchID string `json:"_batch_id"`
				}
				msg, err := waitMsg(processCh)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &res)
				So(res.ID, ShouldEqual, "kk")
				So(res.UUID, ShouldNotBeNil)

				_, err = waitMsg(processCh)
				So(err, ShouldBeNil)

				stored, err := s.R.Get(res.BatchID).Result()
				req := request{}
				json.Unmarshal([]byte(stored), &req)
				So(len(req.Components), ShouldEqual, 2)
				So(req.Service, ShouldEqual, "aaa")
			})
		})
		sub.Unsubscribe()
	})
}
