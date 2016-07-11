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

func TestBasicProcessSuccessResponse(t *testing.T) {
	ch1 := make(chan *nats.Msg, 1)
	ch2 := make(chan *nats.Msg, 1)
	setup()
	s.Setup(os.Getenv("NATS_URI"))
	s.ProcessSuccessResponse("component.process.done", "component.process", "components.process.done")
	Convey("Given the scheduler has a stored request for components.process", t, func() {
		body := `{"service":"foo", "components":[{"_uuid":"ua","id":"a","status":"completed"},{"_uuid":"ub","id":"b"}]}`

		Convey("When it receives a component.process.done message", func() {
			Convey("And the component is the last one to be processed", func() {
				sub, _ := s.N.ChanSubscribe("components.process.done", ch1)
				s.R.Set("uuid", body, 0)
				s.N.Publish("component.process.done", []byte(`{"_batch_id":"uuid","_uuid":"ub","id":"b","status":"completed","extra_field":"extra"}`))
				Convey("Then it should store the result and send a components.process.done", func() {
					msg, err := waitMsg(ch1)
					So(err, ShouldBeNil)
					res := request{}
					json.Unmarshal(msg.Data, &res)
					So(len(res.Components), ShouldEqual, 2)
					So(res.Service, ShouldEqual, "foo")
					So(res.Status, ShouldEqual, "completed")
					So(res.ErrorCode, ShouldEqual, "")
					So(res.ErrorMessage, ShouldEqual, "")
					c1 := s.mapJson(res.Components[0])
					So(c1["status"], ShouldEqual, "completed")
					c2 := s.mapJson(res.Components[1])
					So(c2["status"], ShouldEqual, "completed")
					So(c2["extra_field"], ShouldEqual, "extra")
				})
				sub.Unsubscribe()
			})

			Convey("And the component is not the last one to be processed", func() {
				sub, _ := s.N.ChanSubscribe("component.process", ch2)
				body := `{"service":"foo", "components":[{"_uuid":"ua","id":"a"},{"_uuid":"ub","id":"b"}]}`
				s.R.Set("uuid", body, 0)
				s.N.Publish("component.process.done", []byte(`{"_batch_id":"uuid","_uuid":"ub","id":"b","status":"completed","extra_field":"extra"}`))
				Convey("Then it should store the result and send a component.process for the next component", func() {
					msg, err := waitMsg(ch2)
					So(err, ShouldBeNil)
					So(string(msg.Data), ShouldEqual, `{"_uuid":"ua","id":"a"}`)
				})
				sub.Unsubscribe()
			})
		})
	})
}

func TestBasicProcessErrorResponse(t *testing.T) {
	ch3 := make(chan *nats.Msg, 1)
	setup()
	s.Setup(os.Getenv("NATS_URI"))
	s.ProcessFailedResponse("component.process.error", "components.process.error")
	Convey("Given the scheduler has a stored request for components.process", t, func() {
		body := `{"service":"foo", "components":[{"_uuid":"ua","id":"a","status":"completed"},{"_uuid":"ub","id":"b"}]}`
		Convey("When it receives a component.process.error message", func() {
			sub, _ := s.N.ChanSubscribe("components.process.error", ch3)
			s.R.Set("uuid", body, 0)
			s.N.Publish("component.process.error", []byte(`{"_batch_id":"uuid","_uuid":"ub","id":"b","status":"completed","extra_field":"extra"}`))
			Convey("Then it should store the result and send a components.process.error", func() {
				msg, err := waitMsg(ch3)
				So(err, ShouldBeNil)
				res := request{}
				json.Unmarshal(msg.Data, &res)
				So(len(res.Components), ShouldEqual, 2)
				So(res.Service, ShouldEqual, "foo")
				So(res.Status, ShouldEqual, "completed")
				So(res.ErrorCode, ShouldEqual, "")
				So(res.ErrorMessage, ShouldEqual, "")
				c1 := s.mapJson(res.Components[0])
				So(c1["status"], ShouldEqual, "completed")
				c2 := s.mapJson(res.Components[1])
				So(c2["status"], ShouldEqual, "completed")
				So(c2["extra_field"], ShouldEqual, "extra")
			})
			sub.Unsubscribe()
		})
	})
}
