.PHONY: build install lint test stage
MAKEFLAGS += --no-print-directory
GIT_BRANCH := $(shell git branch --show-current)
GIT_REMOTE := git@github.com:calyrexx/dlog.git
BINARY     := dlog
BUILD_DIR  := .

CHECK_EMOJI := ✅
ERROR_EMOJI := ❌
INFO_EMOJI  := ℹ️
ARROW_UP    := ⬆️

build:
	@echo "$(INFO_EMOJI) Building $(BINARY)..."
	@go build -o $(BUILD_DIR)/$(BINARY) .
	@echo "$(CHECK_EMOJI) Built $(BUILD_DIR)/$(BINARY)"

install:
	@echo "$(INFO_EMOJI) Installing $(BINARY) to GOPATH/bin..."
	@go install .
	@echo "$(CHECK_EMOJI) Installed $(BINARY)"

checks: test lint

test:
	@echo "$(INFO_EMOJI) Running tests..."
	@(go test ./... > test.log 2>&1 || \
		(cat test.log && echo "$(ERROR_EMOJI) Tests failed! Check logs $(ARROW_UP)" && exit 1))
	@rm -f test.log
	@echo "$(CHECK_EMOJI) All tests passed!"

lint:
	@echo "$(INFO_EMOJI) Running linters..."
	@(golangci-lint run ./... > lint.log 2>&1 || \
		(cat lint.log && echo "$(ERROR_EMOJI) Linter found issues! Check logs $(ARROW_UP)" && exit 1))
	@rm -f lint.log
	@echo "$(CHECK_EMOJI) No lint errors found!"

stage:
	@echo "$(INFO_EMOJI) Running pre-commit checks..."
	@$(MAKE) checks
	@echo "$(CHECK_EMOJI) Pre-commit checks passed! Moving to staging..."
	@echo "$(INFO_EMOJI) Staging changes..."
	@git add .
	@git commit -m "$(m)"
	@git push $(GIT_REMOTE) $(GIT_BRANCH)
	@echo "$(CHECK_EMOJI) Changes pushed to $(GIT_BRANCH)!"
