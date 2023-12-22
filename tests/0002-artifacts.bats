#!/usr/bin/env bats

# SPDX-FileCopyrightText: AISEC Pentesting Team
#
# SPDX-License-Identifier: Apache-2.0

load setup.sh

@test "check penrun environment variables" {
	penrun true

	local command
	command="$(trim_string "$(grep PENRUN_COMMAND <true/LATEST/ENV)")"
	[[ "$command" == "PENRUN_COMMAND=true" ]]

	local artifacts
	artifacts="$(trim_string "$(grep PENRUN_ARTIFACTS <true/LATEST/ENV)")"
	[[ "$(basename "$(dirname "${artifacts#*=}")")" == "true" ]]
}

@test "check META.json file" {
	penrun ls -lah

	local meta
	meta="$(<ls/LATEST/META.json)"

	(("$(jq ".exit_code" <<<"$meta")" == "0"))
	[[ "$(jq -r '.command | join(" ")' <<<"$meta")" == "ls -lah" ]]
	[[ "$(jq -r '.start_time' <<<"$meta")" != "" ]]
	[[ "$(jq -r '.end_time' <<<"$meta")" != "" ]]
}

@test "check OUTPUT.zstd file" {
	penrun echo "hans"

	[[ "$(trim_string "$(zstdcat echo/LATEST/OUTPUT.zst)")" == "hans" ]]
}
