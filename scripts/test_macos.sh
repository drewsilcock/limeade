#!/bin/sh

set -eou pipefail

local exit_success=0
local exit_test_failure=1
local exit_script_error=2

if ! command -v pbcopy &> /dev/null; then
  echo "pbcopy is required to run this script"
  exit $exit_script_error
fi

if ! command -v pbpaste &> /dev/null; then
  echo "pbpaste is required to run this script"
  exit $exit_script_error
fi

# 0 = success, 1 = failure
local test_status=0

# Start server in background
./limeade server &

./limeade --version

# Test copy command
local arg_copy_text="I will show you something different from either"
local stdin_copy_text="Your shadow at morning striding behind you or your shadow at evening rising to meet you"

echo "" | pbcopy

./limeade copy $arg_copy_text
if [[ $(./pbpaste) = $arg_copy_text ]]; then
  echo "✅ Copy arg test passed"
else
  echo "❌ Copy stdin test failed"
  test_status=1
fi

echo "" | pbcopy

echo $stdin_copy_text | ./limeade copy
if [[ $(./pbpaste) = $stdin_copy_text ]]; then
  echo "✅ Copy stdin test passed"
else
  echo "❌ Copy arg test failed"
  test_status=1
fi

local paste_text="I will show you fear in a handful of dust."

echo $paste_text | pbcopy

if [[ $(./limeade paste) = $paste_text ]]; then
  echo "✅ Paste test passed"
else
  echo "❌ Paste test failed"
  test_status=1
fi

if [[ $test_status == 0 ]]; then
  echo "✅ All tests passed"
  exit $exit_success
else
  echo "❌ Some tests failed"
  exit $exit_test_failure
fi