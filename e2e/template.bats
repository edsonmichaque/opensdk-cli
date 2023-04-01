#!/usr/bin/env bats

setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'

    DIR="$( cd "$( dirname "$BATS_TEST_FILENAME" )" >/dev/null 2>&1 && pwd )"
    PATH="$DIR/../bin:$PATH"
}

function template() {
    ./bin/template $@
}

@test "template bar" {
    run template bar
    assert_failure
    [ "$status" -eq 1 ]
    [ "${lines[0]}" = "Error: account id is required" ]
}
