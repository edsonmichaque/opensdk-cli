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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/edsonmichaque/template-cli/internal/cmd/formatter"
	"github.com/edsonmichaque/template-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgFmtJSON = "json"
	cfgFmtYAML = "yaml"
	cfgFmtTOML = "toml"
	cfgFmtYML  = "yml"
)

var (
	configProps = map[string]struct{}{
		optAccount:     {},
		optBaseURL:     {},
		optAccessToken: {},
		optSandbox:     {},
	}

	cfgValidateFuncs = map[string]func(string) (interface{}, error){
		optSandbox: func(value string) (interface{}, error) {
			return strconv.ParseBool(value)
		},
		optAccount: func(value string) (interface{}, error) {
			return strconv.ParseInt(value, 10, 64)
		},
		optBaseURL: func(value string) (interface{}, error) {
			return value, nil
		},
		optAccessToken: func(value string) (interface{}, error) {
			return value, nil
		},
	}
)

// cmdCfg
func cmdCfg(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configurations",
	}

	return initCmd(
		cmd,
		withOpts(opts),
		withCmd(
			cmdCfgInit(opts),
			cmdCfgGet(opts),
			cmdCfgSet(opts),
		),
	)
}

// cmdCfgInit
func cmdCfgInit(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadWithValidation(false)
			if err != nil {
				return wrapError(exitFailure, err)
			}

			cfg, ext, err := execConfigPrompt(cfg)
			if err != nil {
				return wrapError(exitFailure, err)
			}

			cfgDir, err := os.UserConfigDir()
			if err != nil {
				return wrapError(exitFailure, err)
			}

			target := filepath.Join(
				cfgDir,
				cmdName,
				fmt.Sprintf(
					"%s.%s",
					viper.GetString(optProfile),
					strings.ToLower(ext),
				),
			)

			if err := writeCfg(cfg, target); err != nil {
				return wrapError(exitFailure, err)
			}

			return nil
		},
	}

	return initCmd(
		cmd,
		withOpts(opts),
	)
}

// writeCfg
func writeCfg(cfg *config.Config, dst string) error {
	v := viper.New()

	v.Set(optAccount, cfg.Account)
	v.Set(optAccessToken, cfg.AccessToken)

	if cfg.BaseURL != "" {
		v.Set(optBaseURL, cfg.BaseURL)
	}

	if cfg.Sandbox {
		v.Set(optSandbox, cfg.Sandbox)
	}

	if err := v.WriteConfigAs(dst); err != nil {
		return wrapError(exitFailure, err)
	}

	return nil
}

// cmdCfgGet
func cmdCfgGet(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Manage configurations",
		Args: func(cmd *cobra.Command, args []string) error {
			prop := args[0]

			return cmdPreRun(
				func() error {
					return cobra.ExactArgs(1)(cmd, args)
				},
				func() error {
					if _, ok := configProps[prop]; !ok {
						return errors.New("not found")
					}

					return nil
				},
			)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &config.Config{
				Account:     "abc123456",
				AccessToken: "def456789",
				BaseURL:     "https://example.com",
			}

			resp, err := formatter.Format(
				formatter.ToConfigList(cfg),
				&formatter.Opts{
					Output: formatter.Output(
						viper.GetString(optOutput),
					),
				},
			)
			if err != nil {
				return wrapError(1, err)
			}

			if err := cmdPrint(cmd, resp); err != nil {
				return wrapError(1, err)
			}

			return nil
		},
	}

	return initCmd(
		cmd,
		withFlagOutput(outputTable),
		withOpts(opts),
	)
}

// cmdCfgSet
func cmdCfgSet(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Manage configurations",
		Args: func(cmd *cobra.Command, args []string) error {
			return cmdPreRun(
				func() error {
					return cobra.ExactArgs(2)(cmd, args)
				},
				func() error {
					if _, ok := configProps[args[0]]; !ok {
						return errors.New("not found")
					}

					return nil
				},
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			validateCfg := cfgValidateFuncs[args[0]]
			if validateCfg == nil {
				return newError(exitFailure, "no validator found")
			}

			value, err := validateCfg(args[1])
			if err != nil {
				return wrapError(exitFailure, err)
			}

			viper.Set(args[0], value)

			if err := viper.WriteConfig(); err != nil {
				return wrapError(exitFailure, err)
			}

			return nil
		},
	}

	return initCmd(
		cmd,
		withOpts(opts),
	)
}
