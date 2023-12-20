#!/usr/bin/env bats

# SPDX-FileCopyrightText: AISEC Pentesting Team
#
# SPDX-License-Identifier: Apache-2.0

setup() {
    # https://bats-core.readthedocs.io/en/stable/tutorial.html#let-s-do-some-setup
    DIR="$( cd "$( dirname "$BATS_TEST_FILENAME" )" >/dev/null 2>&1 && pwd )"
    PATH="$DIR/..:$PATH"
}

@test "invoke penrun without parameters" {
    cd "$BATS_TEST_TMPDIR"

    mkdir foo

    penrun -c /dev/null ls -lah

    if [[ ! -d "ls" ]]; then
        echo "output directory is missing"
        return 1
    fi

    if (( "$(jq ".exit_code" < ls/LATEST/META.json)" != "0" )); then
        echo "exit_code in META != 0"
        return 1
    fi

    if [[ "$(jq -r '.command | join(" ")' < ls/LATEST/META.json)" != "ls -lah" ]]; then
        echo "command != ls -lah"
        return 1
    fi
}
