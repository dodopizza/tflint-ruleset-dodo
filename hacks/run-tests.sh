#!/bin/bash

if [[ -z ${TEST_DIR} ]]; then
    TEST_DIR="./..."
fi

go test \
    -race \
    -cover \
    -covermode=atomic \
    -coverprofile cover.out \
    ${TEST_RUN_ARGS} \
    ${TEST_DIR} \
    -timeout 60s |
    sed "/PASS/s//$(printf "\033[32mPASS\033[0m")/" |
    sed "/FAIL/s//$(printf "\033[31mFAIL\033[0m")/"

exit ${PIPESTATUS[0]}
