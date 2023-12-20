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
    local tmpdir
    tmpdir="$(mktemp -d)"
    (
        cd "$tmpdir"
        penrun -c /dev/null ls -lah
        if [[ ! -d "ls" ]]; then
            echo "output directory is missing"
            return 1
        fi

        mapfile -t meta < "ls/LATEST/META.json"
        for line in "${meta[@]}"; do
            local cmd
            local exit_code
            if [[ "$line" =~ EXIT:(.+) ]]; then
                exit_code="${BASH_REMATCH[1]}"
                if ((exit_code != 0)); then
                    return 1
                fi
            fi
            if [[ "$line" =~ COMMAND:(.+)\s$ ]]; then
                cmd="${BASH_REMATCH[1]}"
                if [[ "$cmd" != "ls -lah" ]]; then
                    echo "$cmd"
                    return 1
                fi
            fi
        done

        rm -rf "$tmpdir"
    )
}
