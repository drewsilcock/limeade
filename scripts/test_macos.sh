#!/bin/bash

set -eou pipefail

function copy {
  echo "$@" | pbcopy
}

function paste {
  pbpaste
}

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
. "$SCRIPT_DIR/test_common.sh"

require_command pbcopy "Are you definitely running from macOS?"
require_command pbpaste "Are you definitely running from macOS?"

echo "Running integration tests for macOS"
run copy paste