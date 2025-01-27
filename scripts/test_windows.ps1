$EXIT_SUCCESS = 0
$EXIT_TEST_FAILURE = 1
$EXIT_SCRIPT_ERROR = 2

$TOTAL_NUM_TESTS = 3

function Require-Command {
    param (
        [string]$cmd,
        [string]$installHelp
    )

    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        Write-Output "$cmd is required to run this script"
        Write-Output "$installHelp"
        exit $EXIT_SCRIPT_ERROR
    }
}

function Start-Server {
    # Start server in background
    Start-Process -NoNewWindow -FilePath "powershell" -ArgumentList "-Command ./limeade server"
}

function Test-CopyStdin {
    $stdinCopyText = "Your shadow at morning striding behind you or your shadow at evening rising to meet you"

    # Clear the clipboard
    Set-Clipboard -Value ""

    $stdinCopyText | ./limeade copy
    if (Get-Clipboard -eq $stdinCopyText) {
        Write-Output "✅ Copy stdin test passed"
        return 0
    } else {
        Write-Output "❌ Copy stdin test failed"
        return 1
    }
}

function Test-CopyArg {
    $argCopyText = "I will show you something different from either"

    # Clear the clipboard
    Set-Clipboard -Value ""

    ./limeade copy $argCopyText
    if (Get-Clipboard -eq $argCopyText) {
        Write-Output "✅ Copy arg test passed"
        return 0
    } else {
        Write-Output "❌ Copy arg test failed"
        return 1
    }
}

function Test-Paste {
    $pasteText = "I will show you fear in a handful of dust."

    # Set the clipboard to the paste text
    Set-Clipboard -Value $pasteText

    if ((./limeade paste) -eq $pasteText) {
        Write-Output "✅ Paste test passed"
        return 0
    } else {
        Write-Output "❌ Paste test failed"
        return 1
    }
}

function Summarise-Tests {
    param (
        [int]$numFails
    )

    if ($numFails -eq 0) {
        Write-Output "✅ All $TOTAL_NUM_TESTS/$TOTAL_NUM_TESTS tests passed"
        exit $EXIT_SUCCESS
    } else {
        Write-Output "❌ $numFails/$TOTAL_NUM_TESTS tests failed"
        exit $EXIT_TEST_FAILURE
    }
}

function Run-Tests {
    $numFails = 0

    Write-Output "Running integration tests for Windows"
    ./limeade --version

    Write-Output "Starting limeade server"
    Start-Server

    Test-CopyStdin
    $numFails += $?

    Test-CopyArg
    $numFails += $?

    Test-Paste
    $numFails += $?

    Summarise-Tests -numFails $numFails
}

Run-Tests