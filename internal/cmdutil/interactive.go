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

package cmdutil

import (
	"github.com/mattn/go-isatty"
	"github.com/spencercjh/sshctx/internal/env"
	"os"
	"os/exec"
)

// isTerminal determines if given fd is a TTY.
func isTerminal(fd *os.File) bool {
	return isatty.IsTerminal(fd.Fd())
}

// fzfInstalled determines if fzf(1) is in PATH.
func fzfInstalled() bool {
	v, _ := exec.LookPath("fzf")

	return v != ""
}

// UseFzf determines if we can do choosing with fzf.
func UseFzf(stdout *os.File) bool {
	v := os.Getenv(env.FZFIgnore)
	return v == "" && isTerminal(stdout) && fzfInstalled()
}

// UsePromptui determines if we can do choosing with promptui.
func UsePromptui(stdout *os.File) bool {
	v := os.Getenv(env.FZFIgnore)
	return isTerminal(stdout) && (!fzfInstalled() || v != "")
}
