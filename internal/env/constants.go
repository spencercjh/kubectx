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

package env

import "regexp"

var SSHParameterRegexp = regexp.MustCompile(`(?m)^(\w+)@((?:1[0-9][0-9]\.|2[0-4][0-9]\.|25[0-5]\.|[1-9][0-9]\.|[0-9]\.){3}(?:1[0-9][0-9]|2[0-4][0-9]|25[0-5]|[1-9][0-9]|[0-9])+|\w+[^\s]+\.[^\s]+)+:(\d+)$`)

const (
	// FZFIgnore describes the environment variable to set to disable
	// interactive context selection when fzf is installed.
	FZFIgnore = "KUBECTX_IGNORE_FZF"

	// NoColor describes the environment variable to disable color usage
	// when printing current context in a list.
	NoColor = `NO_COLOR`

	// ForceColor describes the "internal" environment variable to force
	// color usage to show current context in a list.
	ForceColor = `_KUBECTX_FORCE_COLOR`

	// Debug describes the internal environment variable for more verbose logging.
	Debug = `DEBUG`

	// StrictMode describes the internal environment to force host name with SSHParameterRegexp
	StrictMode = `STRICT_MODE`
)
