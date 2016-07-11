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

func (s *Scheduler) manageRequest(body []byte, to string) {
	req := request{}
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println("Error unmarshalling request: " + string(body))
		return
	}
	batchID := uuid.NewV4()

	for i, c := range req.Components {
		req.Components[i] = s.identifyComponent(*c, batchID)
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Println("Error marshalling result: " + string(body))
		return
	}
	s.R.Set(batchID.String(), string(body), 0)

	for _, c := range req.Components {
		s.publishNext(to, c)
		if req.SequentialProcessing == true {
			break
		}
	}
}

func (s *Scheduler) identifyComponent(c json.RawMessage, batchID uuid.UUID) *json.RawMessage {
	dec := json.NewDecoder(bytes.NewReader(c))
	var a map[string]interface{}
	dec.Decode(&a)

	a["_batch_id"] = batchID
	a["_uuid"] = uuid.NewV4()
	b, err := json.Marshal(a)
	if err != nil {
		log.Println("An error occurred processing the input json: " + string(c))
	}

	m := json.RawMessage(b)

	return &m
}

func (s *Scheduler) publishNext(to string, component *json.RawMessage) {
	s.N.Publish(to, *component)
}
