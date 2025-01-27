# Exit codes
$exit_success = 0
$exit_test_failure = 1
$exit_script_error = 2

# Check if Get-Clipboard and Set-Clipboard cmdlets are available
if (-not (Get-Command Set-Clipboard -ErrorAction SilentlyContinue)) {
    Write-Output "Set-Clipboard is required to run this script"
    exit $exit_script_error
}

if (-not (Get-Command Get-Clipboard -ErrorAction SilentlyContinue)) {
    Write-Output "Get-Clipboard is required to run this script"
    exit $exit_script_error
}

# 0 = success, 1 = failure
$test_status = 0

# Start server in background
Start-Process -NoNewWindow -FilePath "powershell" -ArgumentList "-Command ./limeade server"

./limeade --version

# Test copy command
$arg_copy_text = "I will show you something different from either"
$stdin_copy_text = "Your shadow at morning striding behind you or your shadow at evening rising to meet you"

Set-Clipboard -Value ""

./limeade copy $arg_copy_text
if ((Get-Clipboard) -eq $arg_copy_text) {
    Write-Output "✅ Copy arg test passed"
} else {
    Write-Output "❌ Copy arg test failed"
    $test_status = 1
}

Set-Clipboard -Value ""

$stdin_copy_text | ./limeade copy
if ((Get-Clipboard) -eq $stdin_copy_text) {
    Write-Output "✅ Copy stdin test passed"
} else {
    Write-Output "❌ Copy stdin test failed"
    $test_status = 1
}

# Test paste command

$paste_text = "I will show you fear in a handful of dust."

# Set the clipboard to the paste text
Set-Clipboard -Value $paste_text

if ((./limeade paste) -eq $paste_text) {
    Write-Output "✅ Paste test passed"
} else {
    Write-Output "❌ Paste test failed"
    $test_status = 1
}

if ($test_status -eq 0) {
    Write-Output "✅ All tests passed"
    exit $exit_success
} else {
    Write-Output "❌ Some tests failed"
    exit $exit_test_failure
}