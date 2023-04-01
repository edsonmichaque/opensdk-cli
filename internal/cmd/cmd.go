// Copyright 2023 Edson Michaque
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"io"
	"os"
)

// CmdError
type CmdError struct {
	Code int
	Err  error
}

// Error
func (e CmdError) Error() string {
	return e.Err.Error()
}

// newError
func newError(code int, err string) CmdError {
	return wrapError(code, errors.New(err))
}

// wrapError
func wrapError(code int, err error) CmdError {
	return CmdError{
		Code: code,
		Err:  err,
	}
}

// InitOpts
func InitOpts() (*Opts, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Opts{
		Stdin:   os.Stdin,
		Stderr:  os.Stderr,
		Stdout:  os.Stdout,
		WorkDir: wd,
	}, nil
}

// Opts
type Opts struct {
	Stdout  io.Writer
	Stdin   io.Reader
	Stderr  io.Writer
	WorkDir string
}

// Validate
func (c Opts) Validate() error {
	return nil
}
