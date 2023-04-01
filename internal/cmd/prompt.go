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
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/edsonmichaque/template-cli/internal/config"
)

// execConfigPrompt
func execConfigPrompt(c *config.Config) (*config.Config, string, error) {
	res, err := execPrompt(
		execAccountPrompt(c.Account),
		execAccessTokenPrompt(c.AccessToken),
		execBaseURLPrompt("https://example.con"),
		execFileFmtPrompt(cfgFmtJSON),
		execConfirmPrompt("Do you want to save?", true),
	)
	if err != nil {
		return nil, "", err
	}

	var cfg config.Config

	if value := res.GetString(optAccount); value != "" {
		cfg.Account = value
	}

	if value := res.GetString(optAccessToken); value != "" {
		cfg.AccessToken = value
	}

	if value := res.GetString(optBaseURL); value != "" {
		cfg.BaseURL = value
	}

	var format string
	if value := res.GetString(optFormat); value != "" {
		format = value
	}

	fmt.Printf("%#v\n", cfg)

	return &cfg, format, nil
}

// execAccessTokenPrompt
func execAccessTokenPrompt(value string) runPromptFunc {
	return runPromptFunc(func() (*promptRunnerResult, error) {
		var token string

		if err := survey.AskOne(
			&survey.Input{
				Message: "Access token",
				Default: value,
			},
			&token,
		); err != nil {
			return nil, err
		}

		return &promptRunnerResult{
			Name:  optAccessToken,
			Value: token,
		}, nil
	})
}

// execAccountPrompt
func execAccountPrompt(value string) runPromptFunc {
	return runPromptFunc(func() (*promptRunnerResult, error) {
		var account string

		if err := survey.AskOne(
			&survey.Input{
				Message: "Account",
				Default: value,
			},
			&account,
		); err != nil {
			return nil, err
		}

		return &promptRunnerResult{
			Name:  optAccount,
			Value: account,
		}, nil
	})
}

// execEnvPrompt
func execEnvPrompt(value string) runPromptFunc {
	return runPromptFunc(func() (*promptRunnerResult, error) {
		var env string

		if err := survey.AskOne(
			&survey.Select{
				Message: "Environment",
				Options: []string{
					envProd,
					envSandbox,
					envDev,
				},
				Default: value,
			},
			&env,
		); err != nil {
			return nil, err
		}

		return &promptRunnerResult{
			Name:  optAccessToken,
			Value: env,
		}, nil
	})
}

// execBaseURLPrompt
func execBaseURLPrompt(value string) runPromptFunc {
	return runPromptFunc(func() (*promptRunnerResult, error) {
		var url string

		if err := survey.AskOne(
			&survey.Input{
				Message: "Base URL",
				Default: value,
			},
			&url,
		); err != nil {
			return nil, err
		}

		return &promptRunnerResult{
			Name:  optBaseURL,
			Value: url,
		}, nil
	})
}

// execFileFmtPrompt
func execFileFmtPrompt(value string) runPromptFunc {
	return runPromptFunc(func() (*promptRunnerResult, error) {
		var format string

		if err := survey.AskOne(
			&survey.Select{
				Message: "File format",
				Options: []string{
					cfgFmtJSON,
					cfgFmtYAML,
					cfgFmtTOML,
				},
				Default: value,
			},
			&format,
		); err != nil {
			return nil, err
		}

		return &promptRunnerResult{
			Name:  optFormat,
			Value: format,
		}, nil
	})
}

// execConfirmPrompt
func execConfirmPrompt(msg string, value bool) runPromptFunc {
	return runPromptFunc(func() (*promptRunnerResult, error) {
		var confirmation bool

		if err := survey.AskOne(
			&survey.Confirm{
				Message: msg,
				Default: false,
			},
			&confirmation,
		); err != nil {
			return nil, err
		}

		return &promptRunnerResult{
			Name:  "confirmation",
			Value: confirmation,
		}, nil
	})
}

// promptRunnerResult
type promptRunnerResult struct {
	Name  string
	Value interface{}
}

// promptRunner
type promptRunner interface {
	runPrompt() (*promptRunnerResult, error)
}

// runPromptFunc
type runPromptFunc func() (*promptRunnerResult, error)

// runPrompt
func (rpf runPromptFunc) runPrompt() (*promptRunnerResult, error) {
	return rpf()
}

// promptResponse
type promptResponse map[string]interface{}

// Get
func (rpr promptResponse) Get(key string) interface{} {
	if s, ok := rpr[key]; ok {
		return s
	}

	return nil
}

// GetString
func (rpr promptResponse) GetString(key string) string {
	value := rpr.Get(key)
	if value == nil {
		return ""
	}

	return value.(string)
}

// GetInt
func (rpr promptResponse) GetInt(key string) int {
	value := rpr.Get(key)
	if value == nil {
		return 0
	}

	return value.(int)
}

// GetBool
func (rpr promptResponse) GetBool(key string) bool {
	value := rpr.Get(key)
	if value == nil {
		return false
	}

	return value.(bool)
}

// GetInt64
func (rpr promptResponse) GetInt64(key string) int64 {
	value := rpr.Get(key)
	if value == nil {
		return 0
	}

	return value.(int64)
}

// GetInt32
func (rpr promptResponse) GetInt32(key string) int32 {
	value := rpr.Get(key)
	if value == nil {
		return 0
	}

	return value.(int32)
}

// HasKey
func (rpr promptResponse) HasKey(key string) bool {
	_, ok := rpr[key]

	return ok
}

// execPrompt
func execPrompt(runners ...promptRunner) (promptResponse, error) {
	result := make(map[string]interface{})

	for _, runner := range runners {
		res, err := runner.runPrompt()
		if err != nil {
			return nil, err
		}

		result[res.Name] = res.Value
	}

	return result, nil
}
