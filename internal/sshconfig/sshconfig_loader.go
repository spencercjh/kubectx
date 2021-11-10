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

package sshconfig

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sshctx/internal/cmdutil"
	"sshctx/internal/printer"

	"github.com/pkg/errors"
)

var (
	DefaultLoader Loader = new(StandardLoader)
)

type StandardLoader struct{}

// LoadSSHConfig loads the SSH config from the given path.
// return: sshconfig, sshctxData, error
func (*StandardLoader) LoadSSHConfig() (io.ReadWriteCloser, error) {
	path, err := getSSHConfigPath()
	if err != nil {
		return nil, errors.Wrap(err, "Can't determine sshconfig path")
	}
	file, err := openFile(path, "sshconfig")
	if err != nil {
		return nil, errors.Wrap(err, "Can't open sshconfig")
	}
	return io.ReadWriteCloser(file), nil
}

func (*StandardLoader) LoadSSHCTXData() (io.ReadWriteCloser, error) {
	path, err := GetSSHCtxDataPath()
	if err != nil {
		return nil, errors.Wrap(err, "Can't determine sshconfig path")
	}
	var file *os.File
	// open sshctxData in env: SSHCTX or default one
	file, err = openFile(path, "sshctxData")
	if err != nil {
		_ = printer.Warning(os.Stderr, "Can't open given sshctxData by path: %s", path)
		dir, _ := getSSHCtxDataDir()
		defaultPath := filepath.Join(dir, "config.yaml")
		// try to open default sshctxData
		file, err = openFile(defaultPath, "sshctxData")
		if err != nil {
			_ = printer.Warning(os.Stderr, "Can't open given sshctxData by path: %s", defaultPath)
			// try to create sshctx dir
			if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
				// TODO: consider to ignore the error
				if err := os.Mkdir(dir, 0777); err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("Can't create sshctxData dir: %s", dir))
				}
			}
			// try to create sshctxData
			file, err = os.Create(defaultPath)
			_ = os.Chmod(defaultPath, 0777)
			// TODO: consider to ignore the error
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Can't create sshctxData file: %s", defaultPath))
			}
		}
	}
	return io.ReadWriteCloser(file), nil
}

func getSSHConfigPath() (string, error) {
	// for dev
	if v := os.Getenv("SSHCONFIG"); v != "" {
		list := filepath.SplitList(v)
		if len(list) > 1 {
			// TODO SSHCONFIG=file1:file2 currently not supported
			return "", errors.New("multiple files in SSHCONFIG are currently not supported")
		}
		return v, nil
	}

	// default path is ~/.ssh/config
	home := cmdutil.HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".ssh", "config"), nil
}

func getSSHCtxDataDir() (string, error) {
	home := cmdutil.HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".sshctx"), nil
}

func GetSSHCtxDataPath() (string, error) {
	// for dev
	if v := os.Getenv("SSHCTX"); v != "" {
		list := filepath.SplitList(v)
		if len(list) > 1 {
			// TODO SSHCONFIG=file1:file2 currently not supported
			return "", errors.New("multiple files in SSHCTX are currently not supported")
		}
		return v, nil
	}

	dir, _ := getSSHCtxDataDir()
	return filepath.Join(dir, "config.yaml"), nil
}

func openFile(path string, name string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Wrap(err, fmt.Sprintf("%s doesn't exist", name))
		}
		return nil, errors.Wrap(err, fmt.Sprintf("fail to open %s", name))
	}
	return f, nil
}
