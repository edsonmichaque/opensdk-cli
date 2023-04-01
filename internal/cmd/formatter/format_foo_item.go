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
	"fmt"
	"io"
)

type FooItem Foo

func (f FooItem) FormatText(opts *Opts) (io.Reader, error) {
	keys := []string{
		"id",
		"name",
		"age",
	}

	const txtLen = 8

	values := map[string]interface{}{
		"id":   f.ID,
		"name": f.Name,
		"age":  f.Age,
	}

	titles := map[string]string{
		"id":   "ID",
		"name": "Name",
		"age":  "Age",
	}

	buf := new(bytes.Buffer)
	for _, v := range keys {
		buf.WriteString(fmt.Sprintf("%-20s%v\n", titles[v]+":", values[v]))
	}

	return buf, nil
}

func (d FooItem) FormatJSON(opts *Opts) (io.Reader, error) {
	return formatJSON(d, opts)
}

func (d FooItem) FormatYAML(opts *Opts) (io.Reader, error) {
	return formatYAML(d, opts)
}

func (d FooItem) formatJSON(opts *Opts) ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}
