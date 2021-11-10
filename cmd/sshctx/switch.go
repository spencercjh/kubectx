// Copyright 2021 Google LLC
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

package main

import (
	"bytes"
	"github.com/spencercjh/sshctx/internal/env"
	"github.com/spencercjh/sshctx/internal/printer"
	"github.com/spencercjh/sshctx/internal/sshconfig"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
)

// SwitchOp indicates intention to switch contexts.
type SwitchOp struct {
	Target string // '-' for back and forth, or NAME
}

func (op SwitchOp) Run(_, stderr io.Writer) error {
	var target string
	var err error
	if op.Target == "-" {
		target, err = connectPrevious(stderr)
	} else {
		target, err = connectTarget(op.Target, stderr)
	}
	if err != nil {
		return errors.Wrap(err, "failed to connect host")
	}
	// save previous host
	e := savePreviousHost(target)
	if e != nil {
		return errors.Wrap(e, "failed to save previous host")
	}
	return nil
}

func savePreviousHost(target string) error {
	matches := env.SSHParameterRegexp.FindStringSubmatch(target)
	port, _ := strconv.Atoi(matches[2])
	var host = map[string]sshconfig.Host{"previous": {matches[1], matches[0], port}}
	data, err := yaml.Marshal(&host)

	if err != nil {
		return errors.Wrap(err, "failed to marshal host")
	}

	sshCtxDataPath, _ := sshconfig.GetSSHCtxDataPath()
	if err := ioutil.WriteFile(sshCtxDataPath, data, 0); err != nil {
		return errors.Wrap(err, "failed to write host to sshctxData file")
	}
	return nil
}

// connectTarget switches to specified context name.
func connectTarget(target string, stderr io.Writer) (string, error) {
	_ = printer.Success(stderr, "Switched to target \"%s\".", printer.SuccessColor.Sprint(target))
	cmd := exec.Command("ssh", target)
	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = stderr
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.Is(err, exitError); !ok {
			return target, err
		}
	}
	return target, nil
}

// connectPrevious switches to previously switch context.
func connectPrevious(stderr io.Writer) (string, error) {
	sc := new(sshconfig.SSHConfig).WithLoader(sshconfig.DefaultLoader)

	defer func(sshConfig *sshconfig.SSHConfig) {
		_ = sshConfig.Close()
	}(sc)

	if err := sc.Parse(); err != nil {
		return "", errors.Wrap(err, "sshconfig error")
	}

	if sc.PreviousHost == sshconfig.EmptyHost {
		return "", errors.New("No previous host")
	}

	return connectTarget(sc.PreviousHost.ToSSHParameter(), stderr)
}
