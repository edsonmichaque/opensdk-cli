#!/usr/bin/env bash
# Copyright 2023 Edson Michaque
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0


sed -i "s|edsonmichaque/template-cli|${1}|g" .goreleaser.yaml

find . -type f -name '*.go' -print0 | xargs -0 sed -i "s|edsonmichaque/template-cli|${1}|g"
find . -type f -name '*.mod' -print0 | xargs -0 sed -i "s|edsonmichaque/template-cli|${1}|g"
find . -type f -name '*.go' -print0 | xargs -0 sed -i "s|template|${2}|g"
find . -type f -name '*.go' -print0 | xargs -0 sed -i "s/TEMPLATE/$(echo "$2" | tr a-z A-Z)/g"

mv cmd/template "cmd/${2}"
