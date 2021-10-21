// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package actionagent

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/modules/actionagent/filewatch"
)

func TestAgent_asyncPushCollectorLog(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ctx, cancel := context.WithCancel(context.Background())
	agent := Agent{
		Ctx:         ctx,
		Cancel:      cancel,
		FileWatcher: &filewatch.Watcher{GracefulDoneC: make(chan struct{})},
	}
	go func() {
		time.Sleep(time.Second)
		cancel()
		<-agent.FileWatcher.GracefulDoneC
	}()
	agent.asyncPushCollectorLog()
}
