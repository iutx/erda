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

package k8s

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/scheduler/events"
	eventboxapi "github.com/erda-project/erda/modules/scheduler/events"
	"github.com/erda-project/erda/modules/scheduler/events/eventtypes"
	"github.com/erda-project/erda/modules/scheduler/executor/executortypes"
	"github.com/erda-project/erda/pkg/jsonstore/storetypes"
	"github.com/erda-project/erda/pkg/strutil"
)

func (k *Kubernetes) registerEvent(localStore *sync.Map, stopCh chan struct{}, notifier eventboxapi.Notifier) error {

	name := string(k.name)

	logrus.Infof("in k8s registerEvent, executor: %s", name)

	// watch event handler for a specific etcd directory
	syncRuntimeToLstore := func(key string, value interface{}, t storetypes.ChangeType) error {

		runtimeName := etcdKeyToMapKey(key)
		if len(runtimeName) == 0 {
			return nil
		}

		// Deal with the delete event first
		if t == storetypes.Del {
			_, ok := localStore.Load(runtimeName)
			if ok {
				var e events.RuntimeEvent
				e.RuntimeName = runtimeName
				e.IsDeleted = true
				localStore.Delete(runtimeName)

			}
			return nil
		}

		sg, ok := value.(*apistructs.ServiceGroup)
		if !ok {
			logrus.Errorf("failed to parse value to servicegroup, key: %v, value: %v", key, value)
			return nil
		}

		// Filter events that do not belong to this executor
		if sg.Executor != name {
			return nil
		}

		switch t {
		// Add or update event
		case storetypes.Add, storetypes.Update:
			if _, err := k.Status(context.Background(), *sg); err != nil {
				logrus.Errorf("failed to get k8s status in event, name: %s", sg.ID)
				return nil
			}
			e := GenerateEvent(sg)
			localStore.Store(runtimeName, e)

		default:
			logrus.Errorf("failed to get watch type, type: %s, name: %s", t, runtimeName)
			return nil
		}

		logrus.Infof("succeed to stored key: %s, executor: %s", key, name)
		return nil
	}

	// Correspond the name of the registered executor and its event channel
	getEventFn := func(executorName executortypes.Name) (chan *eventtypes.StatusEvent, chan struct{}, *sync.Map, error) {
		logrus.Infof("in RegisterEvChan executor(%s)", name)
		if string(executorName) == name {
			return k.evCh, stopCh, localStore, nil
		}
		return nil, nil, nil, errors.Errorf("this is for %s executor, not %s", executorName, name)
	}

	return executortypes.RegisterEvChan(executortypes.Name(name), getEventFn, syncRuntimeToLstore)
}

func GenerateEvent(sg *apistructs.ServiceGroup) events.RuntimeEvent {
	var e events.RuntimeEvent
	e.EventType = events.EVENTS_TOTAL
	e.RuntimeName = strutil.Concat(sg.Type, "/", sg.ID)
	e.ServiceStatuses = make([]events.ServiceStatus, len(sg.Services))
	for i, srv := range sg.Services {
		e.ServiceStatuses[i].ServiceName = srv.Name
		e.ServiceStatuses[i].Replica = srv.Scale
		e.ServiceStatuses[i].ServiceStatus = string(srv.Status)
		// Status is empty string
		if len(e.ServiceStatuses[i].ServiceStatus) == 0 {
			e.ServiceStatuses[i].ServiceStatus = convertServiceStatus(apistructs.StatusProgressing)
		}
		if e.ServiceStatuses[i].Replica == 0 {
			e.ServiceStatuses[i].ServiceStatus = convertServiceStatus(apistructs.StatusHealthy)
		}
	}
	return e
}

func convertServiceStatus(serviceStatus apistructs.StatusCode) string {
	switch serviceStatus {
	case apistructs.StatusReady:
		return string(apistructs.StatusHealthy)

	case apistructs.StatusProgressing:
		return string(apistructs.StatusUnHealthy)

	default:
		return string(apistructs.StatusUnknown)
	}
}

// todo: refactor this function
// e.g. /dice/service/services/staging-99 -> services/staging-99
func etcdKeyToMapKey(eKey string) string {
	fields := strings.Split(eKey, "/")
	if l := len(fields); l > 2 {
		return fields[l-2] + "/" + fields[l-1]
	}
	return ""
}
