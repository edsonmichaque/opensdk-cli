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

//go:build e2e
// +build e2e

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.Println("start")

	exit := m.Run()

	log.Println("stop")

	os.Exit(exit)
}

func TestA(t *testing.T) {
	t.Log("a")
}

func TestB(t *testing.T) {
	t.Log("b")
}

func TestBar(t *testing.T) {
	type testcase struct {
		args   []string
		stdin  io.Reader
		stdout io.Writer
		stderr io.Writer
	}

	tt := map[string]testcase{
		"foo": {
			args:   []string{"bar"},
			stdin:  os.Stdin,
			stdout: os.Stdout,
			stderr: os.Stderr,
		},
	}

	for tn, tc := range tt {
		t.Run(tn, func(t *testing.T) {
			process := &Process{
				Command: "../../bin/template",
				Args:    tc.args,
			}

			err := process.Run()
			require.NoError(t, err)
		})
	}
}

type Process struct {
	Command string
	Stdout  *bytes.Buffer
	Stdin   io.Reader
	Stderr  io.Writer
	Args    []string
	Env     []string
}

func (p Process) ErrContains(s string) bool {
}

func (m *Process) Run() error {
	cmd := exec.Command(m.Command, m.Args...)

	cmd.Args = m.Args
	cmd.Stdin = m.Stdin
	cmd.Stdout = m.Stdout
	cmd.Stderr = m.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
