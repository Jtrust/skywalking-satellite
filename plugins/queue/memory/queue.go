// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package memory

import (
	"github.com/enriquebris/goconcurrentqueue"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
	"github.com/apache/skywalking-satellite/plugins/queue/api"
	"github.com/apache/skywalking-satellite/protocol/gen-codes/satellite/protocol"
)

const (
	Name = "memory-queue"
)

type Queue struct {
	config.CommonFields
	// config
	EventBufferSize int `mapstructure:"event_buffer_size"` // The maximum buffer event size.

	// components
	buffer *goconcurrentqueue.FixedFIFO
}

func (q *Queue) Name() string {
	return Name
}

func (q *Queue) Description() string {
	return "This is a memory queue to buffer the input event."
}

func (q *Queue) DefaultConfig() string {
	return `
# The maximum buffer event size.
event_buffer_size: 5000
`
}

func (q *Queue) Initialize() error {
	q.buffer = goconcurrentqueue.NewFixedFIFO(q.EventBufferSize)
	return nil
}

func (q *Queue) Enqueue(e *protocol.Event) error {
	if err := q.buffer.Enqueue(e); err != nil {
		log.Logger.Errorf("error in enqueue: %v", err)
		return api.ErrFull
	}
	return nil
}

func (q *Queue) Dequeue() (*api.SequenceEvent, error) {
	element, err := q.buffer.Dequeue()
	if err != nil {
		log.Logger.Debugf("error in dequeue: %v", err)
		return nil, api.ErrEmpty
	}
	return &api.SequenceEvent{
		Event:  element.(*protocol.Event),
		Offset: "no_offset_in_memory_queue",
	}, nil
}

func (q *Queue) Close() error {
	return nil
}

func (q *Queue) Ack(_ event.Offset) {
}
