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

package actionagent

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (agent *Agent) Execute(r io.Reader) {

	// log level
	debug, _ := strconv.ParseBool(os.Getenv("ACTIONAGENT_DEBUG"))
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Infof("before get openapi info")
	agent.getOpenAPIInfo()
	if len(agent.Errs) > 0 {
		return
	}
	logrus.Infof("before parse arg")
	agent.parseArg(r)
	if len(agent.Errs) > 0 {
		return
	}

	logrus.Infof("before pull bootstrap info")
	agent.pullBootstrapInfo()
	if len(agent.Errs) > 0 {
		return
	}

	logrus.Infof("before validate")
	// 1. validate
	agent.validate()
	if len(agent.Errs) > 0 {
		return
	}

	logrus.Infof("before prepare")
	// 2. prepare
	agent.prepare()
	if len(agent.Errs) > 0 {
		return
	}

	logrus.Infof("before restore")
	// 3. restore / store
	agent.restore()
	if len(agent.Errs) > 0 {
		return
	}
	defer func() {
		agent.store()
	}()

	go agent.ListenSignal()
	go agent.watchFiles()

	logrus.Infof("before logic")
	// 4. logic
	agent.logic()
	if len(agent.Errs) > 0 {
		return
	}
}

func (agent *Agent) parseArg(r io.Reader) {
	// base64 decode
	encodedArg, err := ioutil.ReadAll(r)
	if err != nil {
		agent.AppendError(err)
		return
	}
	decodedArg, err := base64.StdEncoding.DecodeString(string(encodedArg))
	if err != nil {
		agent.AppendError(err)
		return
	}
	agent.Arg = &AgentArg{}
	if err := json.Unmarshal(decodedArg, agent.Arg); err != nil {
		agent.AppendError(err)
		return
	}
}
