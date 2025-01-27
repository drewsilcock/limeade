#!/bin/sh

set -eou pipefail

function copy {
  echo "$@" | xclip -selection clipboard
}

function paste {
  xclip -o -selection clipboard
}

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
. "$SCRIPT_DIR/test_common.sh"

require_command xclip "Install with 'sudo apt install xclip', 'sudo dnf install xclip', 'sudo pacman -S xclip', etc."

echo "Running integration tests for Linux"
run copy paste