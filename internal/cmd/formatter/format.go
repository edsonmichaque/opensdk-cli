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

package formatter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/jmespath/go-jmespath"
	"gopkg.in/yaml.v3"
)

type Output string

const (
	OutputText  = Output("text")
	OutputTable = Output("table")
	OutputJSON  = Output("json")
	OutputYAML  = Output("yaml")
)

type Opts struct {
	Output Output
	Query  string
}

type YAMLFormatter interface {
	FormatYAML(opts *Opts) (io.Reader, error)
}

type JSONFormatter interface {
	FormatJSON(opts *Opts) (io.Reader, error)
}

type TableFormatter interface {
	FormatTable(opts *Opts) (io.Reader, error)
}

type TextFormatter interface {
	FormatText(opts *Opts) (io.Reader, error)
}

func Format(data interface{}, opts *Opts) (io.Reader, error) {
	if opts.Output == OutputJSON {
		if formatter, ok := data.(JSONFormatter); ok {
			return formatter.FormatJSON(opts)
		}

		return nil, errors.New("json formatter is not implemented")
	}

	if opts.Output == OutputYAML {
		if formatter, ok := data.(YAMLFormatter); ok {
			return formatter.FormatYAML(opts)
		}

		return nil, errors.New("yaml formatter is not implemented")
	}

	if opts.Output == OutputTable {
		if formatter, ok := data.(TableFormatter); ok {
			return formatter.FormatTable(opts)
		}

		return nil, errors.New("table formatter is not implemented")
	}

	if opts.Output == OutputText {
		if formatter, ok := data.(TextFormatter); ok {
			return formatter.FormatText(opts)
		}

		return nil, errors.New("table formatter is not implemented")
	}

	return nil, errors.New("invalid formatter")
}

func formatJSON(j jsonFormatter, opts *Opts) (io.Reader, error) {
	data, err := j.formatJSON(opts)
	if err != nil {
		return nil, err
	}

	var result interface{}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	if opts.Query != "" {
		result, err = jmespath.Search(opts.Query, result)
		if err != nil {
			return nil, err
		}
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(out), nil
}

func formatYAML(j jsonFormatter, opts *Opts) (io.Reader, error) {
	data, err := j.formatJSON(opts)
	if err != nil {
		return nil, err
	}

	var result interface{}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	if opts.Query != "" {
		result, err = jmespath.Search(opts.Query, result)
		if err != nil {
			return nil, err
		}
	}

	out, err := yaml.Marshal(result)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(out), nil
}

type jsonFormatter interface {
	formatJSON(opts *Opts) ([]byte, error)
}

type tableFormatter interface {
	formatHeader() []string
	formatRows() []map[string]string
}

func formatTable(t tableFormatter) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)

	if _, err := fmt.Fprintln(tw, strings.Join(t.formatHeader(), "\t")); err != nil {
		return nil, err
	}

	for _, v := range t.formatRows() {
		row := make([]string, 0)

		for _, col := range t.formatHeader() {
			if v, ok := v[col]; ok {
				row = append(row, v)
			}
		}

		if _, err := fmt.Fprintln(tw, strings.Join(row, "\t")); err != nil {
			return nil, err
		}
	}

	if err := tw.Flush(); err != nil {
		return nil, err
	}

	return buf, nil
}
