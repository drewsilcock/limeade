EXIT_SUCCESS=0
EXIT_TEST_FAILURE=1
EXIT_SCRIPT_ERROR=2

TOTAL_NUM_TESTS=3

function require_command {
  cmd="$1"
  install_help="$2"

  if ! command -v $cmd &> /dev/null; then
    echo "$cmd is required to run this script"
    echo "$install_help"
    exit $EXIT_SCRIPT_ERROR
  fi
}

function cleanup {
  local exit_code=$?
  kill -- -$$
  exit $exit_code
}

function start_server {
  trap cleanup SIGINT SIGTERM EXIT
  ./limeade server &
}

function test_copy_stdin {
  copy_cmd=$1
  paste_cmd=$2

  stdin_copy_text="I will show you something different from either"

  $copy_cmd ""

  echo $stdin_copy_text | ./limeade copy
  if [[ $($paste_cmd) = $stdin_copy_text ]]; then
    echo "✅ Copy stdin test passed"
    return 0
  else
    echo "❌ Copy arg test failed"
    return 1
  fi
}

function test_copy_arg {
  copy_cmd=$1
  paste_cmd=$2

  arg_copy_text="Your shadow at morning striding behind you or your shadow at evening rising to meet you"

  $copy_cmd ""

  ./limeade copy "$arg_copy_text"
  if [[ $($paste_cmd) = $arg_copy_text ]]; then
    echo "✅ Copy arg test passed"
    return 0
  else
    echo "❌ Copy stdin test failed"
    return 1
  fi
}

function test_paste {
  copy_cmd=$1
  paste_cmd=$2

  paste_text="I will show you fear in a handful of dust."

  $copy_cmd "$paste_text"

  if [[ $(./limeade paste) = $paste_text ]]; then
    echo "✅ Paste test passed"
    return 0
  else
    echo "❌ Paste test failed"
    return 1
  fi
}

function summarise_tests {
  num_fails=$1

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

  num_fails=0

  ./limeade --version

  echo "Starting limeade server"
  start_server

  test_copy_stdin $copy_cmd $paste_cmd
  num_fails=$(($num_fails+$?))

  test_copy_arg $copy_cmd $paste_cmd
  num_fails=$(($num_fails+$?))

  test_paste $copy_cmd $paste_cmd
  num_fails=$(($num_fails+$?))

  summarise_tests $num_fails
}