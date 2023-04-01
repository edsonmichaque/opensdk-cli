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

package main

import (
	"os"

	"github.com/edsonmichaque/template-cli/internal/cmd"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

func main() {
	os.Exit(runCmd())
}

func runCmd() int {
	if err := cmd.Run(); err != nil {
		if cmdErr, ok := err.(cmd.CmdError); ok {
			return cmdErr.Code
		}

		return exitFailure
	}

	return exitSuccess
}
