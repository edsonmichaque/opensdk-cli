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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	profile    string
)

const (
	cmdName           = "template"
	defaultProfile    = "main"
	envCfgFile        = "TEMPLATE_CONFIG_FILE"
	envCfgHome        = "XDG_CONFIG_HOME"
	envDev            = "DEV"
	envPrefix         = "TEMPLATE"
	envProd           = "PROD"
	envProfile        = "TEMPLATE_PROFILE"
	envSandbox        = "SANDBOX"
	optAccessToken    = "access-token"
	optAccount        = "account"
	optBaseURL        = "base-url"
	optCollaboratorID = "collaborator-id"
	optConfigFile     = "config-file"
	optConfirm        = "confirm"
	optDomain         = "domain"
	optFormat         = "format"
	optFromFile       = "from-file"
	optOutput         = "output"
	optPage           = "page"
	optPerPage        = "per-page"
	optProfile        = "profile"
	optNoInteractive  = "no-interactive"
	optQuery          = "query"
	optRecordID       = "record-id"
	optSandbox        = "sandbox"
	outputJSON        = "json"
	outputTable       = "table"
	outputText        = "text"
	outputYAML        = "yaml"
	pathConfigFile    = "/etc/template"
)

// init
func init() {
	cobra.OnInitialize(initCfg)
	viperBindFlags()
}

// Run
func Run() error {
	return run()
}

// run
func run() error {
	opts, err := InitOpts()
	if err != nil {
		return err
	}

	return runWithOpts(opts)
}

// runWithOpts
func runWithOpts(opts *Opts) error {
	return cmdRoot(opts).Execute()
}

// cmdRoot
func cmdRoot(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use: cmdName,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.PersistentFlags())
		},
		SilenceUsage: true,
	}

	return initCmd(
		cmd,
		withCmd(cmdFoo(opts)),
		withCmd(cmdBar(opts)),
		withCmd(cmdCfg(opts)),
		withCmd(cmdVersion(opts)),
		withFlagsGlobal(),
	)
}

// initCfg
func initCfg() {
	var (
		cfgFile string
		cfgName string
		cfgDir  string
	)

	var err error
	if configFile != "" {
		cfgFile = configFile
	}

	if path := os.Getenv(envCfgFile); path != "" && configFile == "" {
		cfgFile = path
	}

	cfgName = defaultProfile

	if dir := os.Getenv(envCfgHome); dir != "" {
		dir, err = os.UserConfigDir()
		cobra.CheckErr(err)

		cfgDir = dir
	} else {
		if os.Getenv(envCfgHome) != "" {
			dir := os.Getenv(envCfgHome)
			if dir == "" {
				dir, err = os.UserConfigDir()
				cobra.CheckErr(err)
			}

			cfgDir = filepath.Join(dir, cmdName)

			if env := os.Getenv(envProfile); env != "" {
				cfgName = env
			}
		}
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if cfgDir != "" && cfgName != "" {
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName(cfgName)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Found error: ", err.Error())
		}
	}
}

// Cmd
type Cmd struct {
	*cobra.Command
}

// cmdOption
type cmdOption func(*cobra.Command)

// initCmd
func initCmd(cmd *cobra.Command, opts ...cmdOption) *Cmd {
	for _, opt := range opts {
		opt(cmd)
	}

	return &Cmd{
		Command: cmd,
	}
}

// viperBindFlags
func viperBindFlags() {
	for _, env := range os.Environ() {
		envParts := strings.Split(env, "=")
		if len(envParts) != 2 {
			continue
		}

		if !strings.HasPrefix(envParts[0], fmt.Sprintf("%s_", envPrefix)) {
			continue
		}

		flag, err := convertEnvToFlag(env)
		if err != nil {
			continue
		}

		_ = viper.BindEnv(flag, envParts[0])
	}
}

// convertFlagToEnv
func convertFlagToEnv(flag string) string {
	env := strings.ToUpper(strings.ReplaceAll(flag, "-", "_"))

	return fmt.Sprintf("%s_%s", envPrefix, env)
}

// convertEnvToFlag
func convertEnvToFlag(env string) (string, error) {
	envParts := strings.Split(
		strings.TrimPrefix(env, fmt.Sprintf("%s_", envPrefix)), "=",
	)
	if len(envParts) != 2 {
		return "", errors.New("Invalid env var")
	}

	flag := strings.ReplaceAll(strings.ToLower(envParts[0]), "_", "-")

	return flag, nil
}

// cmdPrint
func cmdPrint(cmd *cobra.Command, r io.Reader) error {
	if _, err := io.Copy(cmd.OutOrStdout(), r); err != nil {
		return err
	}

	return nil
}

// flagContains
func flagContains(flag string, values []string) error {
	flagValue := viper.GetString(flag)

	for _, value := range values {
		if flagValue == value {
			return nil
		}
	}

	return fmt.Errorf(`flag "%s" has invalid value "%s"`, flag, flagValue)
}

// cmdPreRun
func cmdPreRun(fn ...func() error) error {
	for _, preRun := range fn {
		if err := preRun(); err != nil {
			return err
		}
	}

	return nil
}
