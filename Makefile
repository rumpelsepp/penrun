# SPDX-FileCopyrightText: Stefan Tatschner
#
# SPDX-License-Identifier: CC0-1.0

.PHONY: lint
lint:
	find . \( -iname "penrun" -or -iname "*.sh" -or -iname "*.bats" \) | xargs shellcheck
	find . \( -iname "penrun" -or -iname "*.sh" -or -iname "*.bats" \) | xargs shfmt -d

.PHONY: fmt
fmt:
	find . \( -iname "penrun" -or -iname "*.sh" -or -iname "*.bats" \) | xargs shfmt -w

.PHONY: test
test:
	bats -x -r tests
