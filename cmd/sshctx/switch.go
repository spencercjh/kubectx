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
	"fmt"
	"github.com/pkg/errors"
	"github.com/spencercjh/sshctx/internal/env"
	"github.com/spencercjh/sshctx/internal/printer"
	"github.com/spencercjh/sshctx/internal/sshconfig"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// SwitchOp indicates intention to switch contexts.
type SwitchOp struct {
	Target string // '-' for back and forth, or NAME
}

func (op SwitchOp) Run(stdout, stderr io.Writer) error {
	var displayName string
	var sshPara string
	var err error
	if op.Target == "-" {
		displayName, sshPara, err = connectPrevious(stderr)
	} else {
		displayName, sshPara, err = connectTarget(op.Target, stderr)
	}
	if err != nil {
		return errors.Wrap(err, "failed to connect host")
	}
	if err := savePreviousHost(stdout, displayName, sshPara); err != nil {
		return errors.Wrap(err, "failed to save previous host")
	}
	return nil
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func savePreviousHost(stdin io.Writer, displayName, sshPara string) error {
	matches := env.SSHParameterRegexp.FindStringSubmatch(sshPara)
	matches = deleteEmpty(matches)
	var host map[string]sshconfig.Host
	switch {
	case len(matches) == 5:
		port, _ := strconv.Atoi(matches[4])
		host = map[string]sshconfig.Host{"previous": {Host: matches[2], DisplayName: displayName, Username: matches[1], Port: port}}
	case len(matches) == 3:
		host = map[string]sshconfig.Host{"previous": {Host: matches[2], DisplayName: displayName, Username: matches[1]}}
	default:
		return fmt.Errorf("illegal SSH parameter: %s", sshPara)
	}

	data, err := yaml.Marshal(&host)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to marshal host: %v", host))
	}

	sshCtxDataPath, _ := sshconfig.GetSSHCtxDataPath()
	if err := ioutil.WriteFile(sshCtxDataPath, data, 0666); err != nil {
		return errors.Wrap(err, "failed to write host to sshctxData file")
	}
	_ = printer.Success(stdin, "Saved previous host successfully: %v", host)
	return nil
}

// extract
func extract(target string) (string, string, error) {
	sshParaBeginIndex := strings.IndexAny(target, "#")
	if sshParaBeginIndex == -1 {
		return "", "", errors.New("invalid target")
	}
	displayName := target[:sshParaBeginIndex]
	sshPara := target[sshParaBeginIndex+1:]
	return displayName, sshPara, nil
}

// connectTarget
func connectTarget(target string, stderr io.Writer) (string, string, error) {
	displayName, sshPara, err := extract(target)
	if err != nil {
		return "", "", err
	}
	return connectTargetWithDisplayName(displayName, sshPara, stderr)
}

// connectTargetWithDisplayName
func connectTargetWithDisplayName(displayName string, sshPara string, stderr io.Writer) (string, string, error) {
	_ = printer.Success(stderr, "Switched to target %s.", printer.SuccessColor.Sprint(displayName))

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		cmd := exec.Command("ssh", "-t", "-t", sshPara)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			_ = printer.Error(stderr, "Failed to connect to target %s because: %v.", printer.ErrorColor.Sprint(displayName), err)
		}
		waitGroup.Done()
	}()
	waitGroup.Wait()
	return displayName, sshPara, nil
}

// connectPrevious switches to previously switch context.
func connectPrevious(stderr io.Writer) (string, string, error) {
	sc := new(sshconfig.SSHConfig).WithLoader(sshconfig.DefaultLoader)

	defer func(sshConfig *sshconfig.SSHConfig) {
		_ = sshConfig.Close()
	}(sc)

	if err := sc.Parse(); err != nil {
		return "", "", errors.Wrap(err, "sshconfig error")
	}

	if sc.PreviousHost == sshconfig.EmptyHost {
		return "", "", errors.New("No previous host")
	}

	return connectTargetWithDisplayName(sc.PreviousHost.DisplayName, sc.PreviousHost.ToSSHParameter(), stderr)
}
