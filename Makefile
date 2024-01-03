# SPDX-FileCopyrightText: Stefan Tatschner
#
# SPDX-License-Identifier: CC0-1.0

GO ?= go

.PHONY: penrun
penrun:
	$(GO) build $(GOFLAGS) -o $@ .

.PHONY: lint
lint:
	find . \( -iname "*.sh" -or -iname "*.bats" \) | xargs shellcheck
	find . \( -iname "*.sh" -or -iname "*.bats" \) | xargs shfmt -d

.PHONY: fmt
fmt:
	find . \( -iname "*.sh" -or -iname "*.bats" \) | xargs shfmt -w

.PHONY: test
test:
	bats -x -r tests
