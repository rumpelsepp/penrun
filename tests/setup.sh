# SPDX-FileCopyrightText: AISEC Pentesting Team
#
# SPDX-License-Identifier: Apache-2.0

setup() {
    # https://bats-core.readthedocs.io/en/stable/tutorial.html#let-s-do-some-setup
    DIR="$( cd "$( dirname "$BATS_TEST_FILENAME" )" >/dev/null 2>&1 && pwd )"
    PATH="$DIR/..:$PATH"
    
    cd "$BATS_TEST_TMPDIR" || exit 1
}

trim_string() {
    sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' <<< "$1"
}
