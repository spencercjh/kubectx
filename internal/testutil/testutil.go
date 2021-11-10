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

package testutil

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// WithEnvVar sets an env var temporarily. Call its return value
// in defer to restore original value in env (if exists).
func WithEnvVar(key, value string) func() {
	orig, ok := os.LookupEnv(key)
	_ = os.Setenv(key, value)
	return func() {
		if ok {
			_ = os.Setenv(key, orig)
		} else {
			_ = os.Unsetenv(key)
		}
	}
}

var noDefaultSSHConfig = false

var cwd, _ = os.Getwd()

func SetupSSHConfig(t *testing.T) {
	sshconfigPath := filepath.Join(os.Getenv("HOME"), ".ssh", "config")
	if _, err := os.Stat(sshconfigPath); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			t.Log("Default sshconfig doesn't exist, create it!")
			noDefaultSSHConfig = true
			existedExampleConfigPath := filepath.Join(cwd, "..", "..", "test", "ssh-config-example")
			exampleConfig, err := os.Open(existedExampleConfigPath)
			if err != nil {
				t.Error("Can't open existedExampleConfig")
			}
			defer func(exampleConfig *os.File) {
				_ = exampleConfig.Close()
			}(exampleConfig)
			sshconfigFile, err := os.Create(sshconfigPath)
			if err != nil {
				t.Error("Can't open sshconfig")
			}
			defer func(sshconfigFile *os.File) {
				_ = sshconfigFile.Close()
			}(sshconfigFile)
			written, err := io.Copy(sshconfigFile, exampleConfig)
			if err != nil {
				t.Error("Can't copy sshconfig")
			}
			t.Logf("Copy %d bytes to sshconfig", written)
		}
	}
}

func TearDownSSHConfig() {
	if noDefaultSSHConfig {
		sshconfigPath := filepath.Join(os.Getenv("HOME"), ".ssh", "config")
		_ = os.Remove(sshconfigPath)
	}
}
