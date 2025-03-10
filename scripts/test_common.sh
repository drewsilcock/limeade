EXIT_SUCCESS=0
EXIT_TEST_FAILURE=1
EXIT_SCRIPT_ERROR=2

TOTAL_NUM_TESTS=4

function require_command {
  local cmd="$1"
  local install_help="$2"

  if ! command -v $cmd &> /dev/null; then
    echo "$cmd is required to run this script"
    echo "$install_help"
    exit $EXIT_SCRIPT_ERROR
  fi
}

function cleanup {
  local exit_code=$?
  kill -- -$$ 2> /dev/null || true
  exit $exit_code
}

function start_server {
  trap cleanup SIGINT SIGTERM EXIT
  go run ./... server &
}

function test_copy_stdin {
  local copy_cmd=$1
  local paste_cmd=$2

  local stdin_copy_text="I will show you something different from either"

  $copy_cmd ""

  echo $stdin_copy_text | go run ./... copy
  if [[ $($paste_cmd) = $stdin_copy_text ]]; then
    echo "✅ Copy stdin test passed"
    return 0
  else
    echo "❌ Copy arg test failed"
    return 1
  fi
}

function test_copy_file_stdin {
  local copy_cmd=$1
  local paste_cmd=$2

  local stdin_copy_text="Your shadow at morning striding behind you"

  $copy_cmd ""

  local fname="/tmp/limeade-test-copy-file-stdin.txt"
  echo $stdin_copy_text > $fname
  cat $fname | go run ./... copy
  if [[ $($paste_cmd) = $stdin_copy_text ]]; then
    echo "✅ Copy file stdin test passed"
    return 0
  else
    echo "❌ Copy file stdin test failed"
    return 1
  fi
}

function test_copy_arg {
  local copy_cmd=$1
  local paste_cmd=$2

  local arg_copy_text="Or your shadow at evening rising to meet you"

  $copy_cmd ""

  go run ./... copy "$arg_copy_text"
  if [[ $($paste_cmd) = $arg_copy_text ]]; then
    echo "✅ Copy arg test passed"
    return 0
  else
    echo "❌ Copy stdin test failed"
    return 1
  fi
}

function test_paste {
  local copy_cmd=$1
  local paste_cmd=$2

  local paste_text="I will show you fear in a handful of dust."

  $copy_cmd "$paste_text"

  if [[ $(go run ./... paste) = $paste_text ]]; then
    echo "✅ Paste test passed"
    return 0
  else
    echo "❌ Paste test failed"
    return 1
  fi
}

function summarise_tests {
  local num_fails=$1

  if [[ $num_fails == 0 ]]; then
    echo "✅ All $TOTAL_NUM_TESTS/$TOTAL_NUM_TESTS tests passed"
    exit $EXIT_SUCCESS
  else
    echo "❌ $num_fails/$TOTAL_NUM_TESTS tests failed"
    exit $EXIT_TEST_FAILURE
  fi
}

function run {
  copy_cmd=$1
  paste_cmd=$2

  local num_fails=0

  go run ./... --version

  echo "Starting limeade server"
  start_server

  test_copy_stdin $copy_cmd $paste_cmd
  num_fails=$(($num_fails+$?))

  test_copy_file_stdin $copy_cmd $paste_cmd
  num_fails=$(($num_fails+$?))

  test_copy_arg $copy_cmd $paste_cmd
  num_fails=$(($num_fails+$?))

  test_paste $copy_cmd $paste_cmd
  num_fails=$(($num_fails+$?))

  summarise_tests $num_fails
}