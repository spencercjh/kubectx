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
	"facette.io/natsort"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spencercjh/sshctx/internal/env"
	"github.com/spencercjh/sshctx/internal/printer"
	"github.com/spencercjh/sshctx/internal/sshconfig"
	"io"
	"os"
	"sort"
	"strings"
)

type PromptuiOp struct{}

func (op PromptuiOp) Run(stdout, stderr io.Writer) error {
	sc := new(sshconfig.SSHConfig).WithLoader(sshconfig.DefaultLoader)

	defer func(sshConfig *sshconfig.SSHConfig) {
		_ = sshConfig.Close()
	}(sc)

	if err := sc.Parse(); err != nil {
		return errors.Wrap(err, "sshconfig error")
	}

	sort.SliceStable(sc.Hosts, func(i, j int) bool {
		return natsort.Compare(sc.Hosts[i].Host, sc.Hosts[j].Host)
	})

	items := []string{}
	for _, h := range sc.Hosts {
		str := h.ToSSHParameter()
		_, ok := os.LookupEnv(env.StrictMode)
		if ok && !env.SSHParameterRegexp.MatchString(str) {
			_ = printer.Warning(stdout, "%s is an illegal ssh parameter", str)
			continue
		}
		str = "💻: " + h.DisplayName + "#" + str
		if h == sc.PreviousHost {
			str = printer.ActiveItemColor.Sprint(str)
		}
		items = append(items, str)
	}

	prompt := promptui.Select{
		Label: "Select a host to connect",
		Items: items,
	}
	_, result, err := prompt.Run()

	if err != nil {
		return errors.Wrap(err, "promptui error")
	}
	choice := strings.TrimSpace(result)
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
