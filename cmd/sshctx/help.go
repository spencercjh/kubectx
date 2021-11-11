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
	"io"
	"strings"

	"github.com/pkg/errors"
)

// HelpOp describes printing help.
type HelpOp struct{}

func (op HelpOp) Run(stdout, _ io.Writer) error {
	return printUsage(stdout)
}

func printUsage(out io.Writer) error {
	help := `USAGE:
  %PROG%                       : list the hosts
  %PROG% <HOST>                : connect to <HOST>
  %PROG% -                     : connect to the previous successfully connected host
  %PROG% -p, --previous        : show the previous successfully connected host
  %PROG% -h,--help             : show this message
  %PROG% -v,-V,--version       : show version`
	help = strings.ReplaceAll(help, "%PROG%", selfName())

	_, err := fmt.Fprintf(out, "%s\n", help)
	return errors.Wrap(err, "write error")
}

// selfName
func selfName() string {
	return "sshctx"
}
