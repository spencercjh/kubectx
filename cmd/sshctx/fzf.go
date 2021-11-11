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
	"fmt"
	"github.com/spencercjh/sshctx/internal/env"
	"github.com/spencercjh/sshctx/internal/sshconfig"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type InteractiveSwitchOp struct {
	SelfCmd string
}

func (op InteractiveSwitchOp) Run(stdout, stderr io.Writer) error {
	// parse sshconfig just to see if it can be loaded
	sc := new(sshconfig.SSHConfig).WithLoader(sshconfig.DefaultLoader)
	defer func(sshConfig *sshconfig.SSHConfig) {
		_ = sshConfig.Close()
	}(sc)

	if err := sc.Parse(); err != nil {
		return errors.Wrap(err, "sshconfig error")
	}

	cmd := exec.Command("fzf", "--ansi", "--no-preview")
	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = stderr
	cmd.Stdout = &out

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", op.SelfCmd),
		fmt.Sprintf("%s=1", env.ForceColor))
	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.Is(err, exitError); !ok {
			return err
		}
	}
	choice := strings.TrimSpace(out.String())
	if choice == "" {
		return errors.New("you did not choose any of the options")
	}
	displayName, sshPara, err := connectTarget(choice, stderr)
	if err != nil {
		return errors.Wrap(err, "failed to switch host")
	}
	if err := savePreviousHost(stdout, displayName, sshPara); err != nil {
		return errors.Wrap(err, "failed to save previous host")
	}
	return nil
}
