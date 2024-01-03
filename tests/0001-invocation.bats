#!/usr/bin/env bats

# SPDX-FileCopyrightText: Stefan Tatschner
#
# SPDX-License-Identifier: MIT

load setup.bash

@test "invoke penrun without parameters" {
	mkdir foo

	penrun ls -lah

	[[ -d "penrun-artifacts" ]]
	[[ -r "penrun-artifacts/ls/LATEST/META.json.zst" ]]
	[[ -r "penrun-artifacts/ls/LATEST/OUTPUT.zst" ]]
}

@test "invoke penrun with set artifacts basedir" {
	mkdir artifacts

	penrun -d artifacts true

	[[ -r "artifacts/META.json.zst" ]]
	[[ -r "artifacts/OUTPUT.zst" ]]
}

@test "invoke penrun with set artifacts dir from config" {
	echo 'artifacts-dir = "artifacts"' >penrun.toml

	mkdir artifacts

	penrun true

	[[ -r "artifacts/META.json.zst" ]]
	[[ -r "artifacts/OUTPUT.zst" ]]
}

@test "invoke penrun with set artifacts basedir from config" {
	echo 'artifacts-base = "artifacts"' >penrun.toml

	mkdir artifacts

	penrun true

	[[ -r "artifacts/true/LATEST/META.json.zst" ]]
	[[ -r "artifacts/true/LATEST/OUTPUT.zst" ]]
}
