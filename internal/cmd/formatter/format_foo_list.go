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
	"encoding/json"
	"fmt"
	"io"
)

type Foo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

type FooList []Foo

func (f FooList) FormatJSON(opts *Opts) (io.Reader, error) {
	return formatJSON(f, opts)
}

func (f FooList) FormatYAML(opts *Opts) (io.Reader, error) {
	return formatYAML(f, opts)
}

func (f FooList) FormatTable(_ *Opts) (io.Reader, error) {
	return formatTable(f)
}

func (f FooList) formatJSON(opts *Opts) ([]byte, error) {
	return json.MarshalIndent(f, "", "  ")
}

func (f FooList) formatHeader() []string {
	return []string{
		"ID",
		"NAME",
		"AGE",
	}
}

func (f FooList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(f))

	fooList := f

	const txtLen = 10

	for i := range fooList {
		data = append(data, map[string]string{
			"ID":   fmt.Sprintf("%d", fooList[i].ID),
			"NAME": fooList[i].Name,
			"AGE":  fooList[i].Age,
		})
	}

	return data
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}

	return s[:length] + "..."
}
