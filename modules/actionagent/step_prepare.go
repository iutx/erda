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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/filehelper"
)

const (
	EnvStdErrRegexpList = "ACTIONAGENT_STDERR_REGEXP_LIST"
)

// 对于 custom action，需要将 commands 转换为 script 来执行
func (agent *Agent) prepare() {
	cmd := exec.Command("ls", "-l", "/home", "/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

	cmd = exec.Command("whoami")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

	if err := os.MkdirAll("/tmp/pipeline/container", 0777); err != nil {
		logrus.Errorf("create dir error: %v", err)
		agent.AppendError(err)
	}

	cmd = exec.Command("ls", "-la", "/tmp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

	// 1. create contextDir/workDir/dir(metaFile)/uploadDir/tempTarUploadDir
	if err := os.MkdirAll(agent.EasyUse.ContainerContext, 0777); err != nil {
		logrus.Errorf("create dir %s error: %v", agent.EasyUse.ContainerContext, err)
		agent.AppendError(err)
	}

	cmd = exec.Command("ls", "-la", "/home")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

	if err := os.MkdirAll(agent.EasyUse.ContainerWd, 0777); err != nil {
		logrus.Errorf("create dir %s error: %v", agent.EasyUse.ContainerWd, err)
		agent.AppendError(err)
	}
	if err := os.MkdirAll(filepath.Dir(agent.EasyUse.ContainerMetaFile), 0777); err != nil {
		logrus.Errorf("create dir %s error: %v", agent.EasyUse.ContainerMetaFile, err)
		agent.AppendError(err)
	}
	if agent.EasyUse.ContainerUploadDir != "" {
		if err := os.MkdirAll(agent.EasyUse.ContainerUploadDir, 0777); err != nil {
			logrus.Errorf("create dir %s error: %v", agent.EasyUse.ContainerUploadDir, err)
			agent.AppendError(err)
		}
	}
	if err := os.Mkdir(agent.EasyUse.ContainerTempTarUploadDir, 0777); err != nil {
		logrus.Errorf("create dir %s error: %v", agent.EasyUse.ContainerTempTarUploadDir, err)
		agent.AppendError(err)
	}
	// {
	//     "name":"a.cert",
	//     "value":"/.pipeline/container/context/.cms/dice_files/a.cert",
	//     "labels":{
	//         "diceFileUUID":"d31b0b31e85c467c8a54e4a9786363b7"
	//     }
	// }
	for _, f := range agent.Arg.Context.CmsDiceFiles {
		if err := os.MkdirAll(filepath.Dir(f.Value), 0777); err != nil {
			logrus.Errorf("create dir %s error: %v", filepath.Dir(f.Value), err)
			agent.AppendError(err)
		}
	}

	// 2. create custom script
	if len(agent.Arg.Commands) > 0 {
		if err := agent.setupScript(); err != nil {
			agent.AppendError(err)
			return
		}
	}

	// 3. compatible when_sigterm -> when_sig_15
	const oldSigTERMScript = "/opt/action/when_sigterm"
	if err := filehelper.CheckExist(oldSigTERMScript, false); err == nil {
		if err := filehelper.Copy(oldSigTERMScript, getSigScriptPath(syscall.SIGTERM)); err != nil {
			agent.AppendError(err)
			return
		}
	}

	// 4. multiWriter of stdout/stderr
	if f, err := filehelper.CreateFile3(agent.EasyUse.RunMultiStdoutFilePath, bytes.NewBufferString(""), 0644); err != nil {
		logrus.Printf("failed to create multi stdout, err: %v\n", err)
	} else {
		agent.EasyUse.RunMultiStdout = f
	}
	if f, err := filehelper.CreateFile3(agent.EasyUse.RunMultiStderrFilePath, bytes.NewBufferString(""), 0644); err != nil {
		logrus.Printf("failed to create multi stderr, err: %v\n", err)
	} else {
		agent.EasyUse.RunMultiStderr = f
	}

	// 5. set stderr regexp list
	envStdErrRegexpStr := os.Getenv(EnvStdErrRegexpList)
	regexpStrList := []string{}
	if err := json.Unmarshal([]byte(envStdErrRegexpStr), &regexpStrList); err != nil {
		agent.AppendError(err)
		return
	}
	for i := range regexpStrList {
		reg, err := regexp.Compile(regexpStrList[i])
		if err != nil {
			agent.AppendError(err)
		} else {
			agent.StdErrRegexpList = append(agent.StdErrRegexpList, reg)
		}
	}
}

func (agent *Agent) setupScript() error {
	var buf bytes.Buffer
	for _, command := range agent.Arg.Commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, `$`, `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}
	script := fmt.Sprintf(
		buildScript,
		buf.String(),
	)

	if err := filehelper.CreateFile(agent.EasyUse.RunScript, script, 0777); err != nil {
		return err
	}

	return nil
}

// buildScript is a helper script which add a shebang
// to the generated script.
const buildScript = `#!/bin/sh
set -e
%s
`

// traceScript is a helper script which is added to the
// generated script to trace each command.
const traceScript = `
echo + %s
%s || ((echo "- FAIL! exit code: $?") && false)
echo
`
