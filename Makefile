default: check

.PHONY: check
check: tidy lint test

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: test
test:
	TEST_RUN_ARGS="$(TEST_RUN_ARGS)" TEST_DIR="$(TEST_DIR)" ./hacks/run-tests.sh

.PHONY: lint
lint:
	golangci-lint run
