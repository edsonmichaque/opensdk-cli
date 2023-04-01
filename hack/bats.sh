#!/usr/bin/env sh
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


if [ ! -d "e2e/bats" ]; then
    git clone https://github.com/bats-core/bats e2e/bats
fi

if [ ! -d "e2e/test_helper/bats-assert" ]; then
    git clone https://github.com/bats-core/bats-assert e2e/test_helper/bats-assert
fi

if [ ! -d "e2e/test_helper/bats-support" ]; then
    git clone https://github.com/bats-core/bats-support e2e/test_helper/bats-support
fi
