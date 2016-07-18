/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package scheduler

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/satori/go.uuid"
)

type request struct {
	Service              string             `json:"service"`
	Status               string             `json:"status"`
	ErrorCode            string             `json:"error_code"`
	ErrorMessage         string             `json:"error_message"`
	Components           []*json.RawMessage `json:"components"`
	SequentialProcessing bool               `json:"sequential_processing"`
	Action               string             `json:"action"`
}

func (s *Scheduler) manageRequest(inputBody []byte, from, to string) {
	req := request{}
	if err := json.Unmarshal(inputBody, &req); err != nil {
		log.Println("Error unmarshalling request: " + string(inputBody))
		return
	}

	if req.Status == "completed" {
		s.N.Publish(from+".done", inputBody)
		return
	}

	batchID := uuid.NewV4()

	for i, c := range req.Components {
		req.Components[i] = s.identifyComponent(*c, batchID, req)
	}

	somethingToProcess := false
	for i, c := range req.Components {
		if s.readyToBeProcessed(c) == true {
			somethingToProcess = true
			// TODO update status as in progress
			c = s.markAs(*c, "processing")
			req.Components[i] = c
			s.publishNext(to, c)
			if req.SequentialProcessing == true {
				break
			}
		}
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Println("Error marshalling result: " + string(body))
		return
	}
	s.R.Set(batchID.String(), string(body), 0)

	if somethingToProcess == false {
		s.N.Publish(from+".done", inputBody)
		return
	}
}
func (s *Scheduler) markAs(c json.RawMessage, status string) *json.RawMessage {
	dec := json.NewDecoder(bytes.NewReader(c))
	var a map[string]interface{}
	dec.Decode(&a)
	a["status"] = status

	b, err := json.Marshal(a)
	if err != nil {
		log.Println("An error occurred processing the input json: " + string(c))
	}

	m := json.RawMessage(b)
	return &m
}

func (s *Scheduler) identifyComponent(c json.RawMessage, batchID uuid.UUID, req request) *json.RawMessage {
	dec := json.NewDecoder(bytes.NewReader(c))
	var a map[string]interface{}
	dec.Decode(&a)

	a["_batch_id"] = batchID
	a["_uuid"] = uuid.NewV4()
	a["service"] = req.Service
	b, err := json.Marshal(a)
	if err != nil {
		log.Println("An error occurred processing the input json: " + string(c))
	}

	m := json.RawMessage(b)

	return &m
}

func (s *Scheduler) readyToBeProcessed(c *json.RawMessage) bool {
	dec := json.NewDecoder(bytes.NewReader(*c))
	var a map[string]interface{}
	dec.Decode(&a)

	if a["status"] == nil {
		return true
	}

	if a["status"].(string) != "completed" {
		return true
	}

	return false
}

func (s *Scheduler) publishNext(to string, component *json.RawMessage) {
	s.N.Publish(to, *component)
}
