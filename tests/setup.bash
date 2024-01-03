# SPDX-FileCopyrightText: Stefan Tatschner
#
# SPDX-License-Identifier: MIT

setup() {
	# https://bats-core.readthedocs.io/en/stable/tutorial.html#let-s-do-some-setup
	DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")" >/dev/null 2>&1 && pwd)"
	PATH="$DIR/..:$PATH"

	cd "$BATS_TEST_TMPDIR" || exit 1
}

