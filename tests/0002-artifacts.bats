#!/usr/bin/env bats

# SPDX-FileCopyrightText: Stefan Tatschner
#
# SPDX-License-Identifier: MIT

load setup.bash
load lib.bash

@test "check META.json" {
	penrun ls -lah

	local meta
	meta="$(zstdcat penrun-artifacts/ls/LATEST/META.json.zst)"

	(("$(jq ".exit_code" <<<"$meta")" == "0"))
	[[ "$(jq -r '.cli | join(" ")' <<<"$meta")" == "ls -lah" ]]
	[[ "$(jq -r '.cli_string' <<<"$meta")" == "ls -lah" ]]
	[[ "$(jq -r '.start_time' <<<"$meta")" != "" ]]
	[[ "$(jq -r '.end_time' <<<"$meta")" != "" ]]
}

@test "check environment variables in META.json" {
	penrun ls -lah

	local meta
	meta="$(zstdcat penrun-artifacts/ls/LATEST/META.json.zst)"

	jq -r '.environ' <<<"$meta" | grep "PENRUN_CLI_STRING=ls -lah"
	jq -r '.environ' <<<"$meta" | grep -E "PENRUN_ARTIFACTS_DIR=.*/ls/run-.*"
}

@test "check OUTPUT.zstd file" {
	penrun echo "hans"

	[[ "$(trim_string "$(zstdcat penrun-artifacts/echo/LATEST/OUTPUT.zst)")" == "hans" ]]
}
