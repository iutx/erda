// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// Package event manipulates the k8s api of event object
package event

import (
	"github.com/erda-project/erda/pkg/clientgo/kubernetes"
)

// Event is the object to encapsulate docker
type Event struct {
	client *kubernetes.Clientset
}

// Option configures an Event
type Option func(*Event)

// New news an Event
func New(options ...Option) *Event {
	ns := &Event{}

	for _, op := range options {
		op(ns)
	}

	return ns
}

// WithKubernetesClient provides an Option
func WithKubernetesClient(client *kubernetes.Clientset) Option {
	return func(event *Event) {
		event.client = client
	}
}
