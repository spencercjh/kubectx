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
	"github.com/spencercjh/sshctx/internal/printer"
	"github.com/spencercjh/sshctx/internal/sshconfig"
	"io"

	"github.com/pkg/errors"
)

type PreviousOp struct{}

func (op PreviousOp) Run(stdout, _ io.Writer) error {
	sc := new(sshconfig.SSHConfig).WithLoader(sshconfig.DefaultLoader)

	defer func(sshConfig *sshconfig.SSHConfig) {
		_ = sshConfig.Close()
	}(sc)

	if err := sc.Parse(); err != nil {
		return errors.Wrap(err, "sshconfig error")
	}
	if sc.PreviousHost == sshconfig.EmptyHost {
		return errors.New("No previous host in sshctx")
	}
	_ = printer.Success(stdout, "Previous host: %s", sc.PreviousHost.ToSSHParameter())
	return nil
}
