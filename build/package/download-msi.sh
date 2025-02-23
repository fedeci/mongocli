#!/usr/bin/env bash
# Copyright 2021 MongoDB Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -Eeou pipefail

VERSION="$(git tag --list "${TOOL_NAME:?}/v*" --sort=taggerdate | tail -1 | cut -d "v" -f 2)"
if [[ -z "${VERSION}" ]]; then
    VERSION="$(git describe --abbrev=0 | cut -d "v" -f 2)"
fi

PACKAGE_NAME=mongocli_"${VERSION}"_windows_x86_64.msi
if [[ "${TOOL_NAME:?}" == atlascli ]]; then
  PACKAGE_NAME=mongodb-atlas-cli_${VERSION}_windows_x86_64.msi
fi

pushd bin

curl https://fastdl.mongodb.org/mongocli/"${PACKAGE_NAME}" --output "${PACKAGE_NAME}"
