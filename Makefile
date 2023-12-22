# SPDX-FileCopyrightText: Stefan Tatschner
#
# SPDX-License-Identifier: CC0-1.0

.PHONY: lint
lint:
	find . -iname "*.sh" -or -iname "*.bats" -exec shellcheck '{}' \;

.PHONY: fmt
fmt:
	find . -iname "*.sh" -or -iname "*.bats" -exec shfmt -w '{}' \;
	
.PHONY: test
test:
	bats -x -r tests
