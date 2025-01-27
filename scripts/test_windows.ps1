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

function Cleanup-LimeadeProcesses {
    Stop-Process -Name limeade -ErrorAction SilentlyContinue
    Remove-Item -Path /tmp/limeade.sock -ErrorAction SilentlyContinue
}

function Start-Server {
    Cleanup-LimeadeProcesses
    Start-Process -NoNewWindow -FilePath "./limeade" -ArgumentList "server"
}

function Test-CopyStdin {
    $stdinCopyText = "Your shadow at morning striding behind you or your shadow at evening rising to meet you"

    # Clear the clipboard
    Set-Clipboard -Value ""

    $stdinCopyText | ./limeade copy
    $output = Get-Clipboard | Out-String
    if ($output.Trim() -eq $stdinCopyText) {
        return 0
    } else {
        return 1
    }
}

function Test-CopyArg {
    $argCopyText = "I will show you something different from either"

    # Clear the clipboard
    Set-Clipboard -Value ""

    ./limeade copy $argCopyText
    $output = Get-Clipboard | Out-String
    if ($output.Trim() -eq $argCopyText) {
        return 0
    } else {
        return 1
    }
}

function Test-Paste {
    $pasteText = "I will show you fear in a handful of dust."

    # Set the clipboard to the paste text
    Set-Clipboard -Value $pasteText

    $output = ./limeade paste | Out-String
    if ($output.Trim() -eq $pasteText) {
        return 0
    } else {
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

    $result = Test-CopyStdin
    if ($result -eq 0) {
        Write-Output "✅ Copy arg test passed"
    } else {
        Write-Output "❌ Copy arg test failed"
        $numFails += 1
    }

    $result = Test-CopyArg
    if ($result -eq 0) {
        Write-Output "✅ Copy arg test passed"
    } else {
        Write-Output "❌ Copy arg test failed"
        $numFails += 1
    }

    $result = Test-Paste
    if ($result -eq 0) {
        Write-Output "✅ Paste test passed"
    } else {
        Write-Output "❌ Paste test failed"
        $numFails += 1
    }

    Summarise-Tests -numFails $numFails
}

try {
    Run-Tests
} finally {
    Cleanup-LimeadeProcesses
}