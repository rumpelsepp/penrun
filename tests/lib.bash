# SPDX-FileCopyrightText: Stefan Tatschner
#
# SPDX-License-Identifier: MIT

trim_string() {
	sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' <<<"$1"
}
