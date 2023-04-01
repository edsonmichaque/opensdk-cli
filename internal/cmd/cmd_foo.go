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
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/edsonmichaque/template-cli/internal/cmd/formatter"
	"github.com/edsonmichaque/template-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cmdFoo
func cmdFoo(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "foo",
		Short: "List accounts",
		Example: heredoc.Doc(`
			template foo
			template foo --output=json
			template foo --output=yaml
			template foo --output=json --query="[].id"
		`),
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cmdPreRun(
				func() error {
					return viper.BindPFlags(cmd.Flags())
				},
				func() error {
					return flagContains(
						optOutput,
						[]string{
							outputJSON,
							outputYAML,
							outputTable,
						},
					)
				},
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.Init()
			if err != nil {
				return wrapError(exitFailure, err)
			}

			fooList := formatter.FooList{
				formatter.Foo{
					ID:   1,
					Name: "First Name",
					Age:  "19",
				},
				formatter.Foo{
					ID:   2,
					Name: "First Name",
					Age:  "19",
				},
			}

			fooOutput, err := formatter.Format(
				fooList, &formatter.Opts{
					Output: formatter.Output(
						viper.GetString(optOutput),
					),
				},
			)
			if err != nil {
				return wrapError(exitFailure, err)
			}

			if err := cmdPrint(cmd, fooOutput); err != nil {
				return wrapError(exitFailure, err)
			}

			return nil
		},
	}

	return initCmd(
		cmd,
		withFlagOutput(outputTable),
		withFlagQuery(),
		withOpts(opts),
	)
}
