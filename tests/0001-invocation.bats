#!/usr/bin/env bats

# SPDX-FileCopyrightText: AISEC Pentesting Team
#
# SPDX-License-Identifier: Apache-2.0

load setup.sh

@test "invoke penrun without parameters" {
	mkdir foo

	penrun ls -lah

	[[ -d "penrun-artifacts" ]]
}
