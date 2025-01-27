# 🍋‍🟩 Limeade

[![Lint](https://github.com/drewsilcock/limeade/actions/workflows/lint.yaml/badge.svg)](https://github.com/drewsilcock/limeade/actions/workflows/lint.yaml)
[![Test](https://github.com/drewsilcock/limeade/actions/workflows/test.yaml/badge.svg)](https://github.com/drewsilcock/limeade/actions/workflows/test.yaml)

remote...lemote...lemode......Lemonade.........Limeade! 🍋‍🟩

Limeade is a tool for remote clipboard access, i.e. copy and paste over SSH.

This is a fork of [lemonade-command/lemonade](https://github.com/lemonade-command/lemonade).

## Installation

### Installer script

curl -sSfL https://raw.githubusercontent.com/drewsilcock/limeade/main/install.sh | sh -s

### Go install

```sh
go install github.com/drewsilcock/limeade@latest
```

### Build from source

```sh
git clone git@github.com:drewsilcock/limeade
cd limeade
go build
```

## Usage

For options, run `limeade --help`.

First, run the server on the machine you want to access the clipboard of.

```sh
limeade server
````

Then, SSH into a remote machine:

```shell
ssh -R /tmp/limeade.sock:/tmp/limeade.sock user@host
```

Finally, run the client on the remote machine:

```shell
limeade copy 'Hello, world!'

# You can also pipe text into the copy command
echo 'Hello, world!' | limeade copy
```

You can also add the following to your `~/.ssh/config` file so that you don't need to specify the `-R` flag every time you SSH into the remote machine:

```shell
Host host
    User user
    HostName host
    ...
    RemoteForward /tmp/limeade.sock localhost:/tmp/limeade.sock
```

### Alias

You can use limeade as an alias of `pbcopy` and `pbpaste`.

For example.

```sh
$ ln -s /path/to/limeade /usr/bin/pbcopy
$ echo "Hello" | pbcopy  # Same as 'echo "Hello" | limeade copy'
```

## Working with neovim

Neovim will automatically detect lemonade in a remote environment, but not limeade. For this reason, you can just alias limeade to lemonade:

```shell
ln -s /path/to/limeade /usr/bin/lemonade
```

## Limeade vs. Lemonade

This project is a fork of [lemonade-command/lemonade](https://github.com/lemonade-command/lemaonde) with a few key differences:

- Modernise codebase and keep dependencies up-to-date.
- Remove `open` functionality as this is unneeded and a security concern. You can just paste into your browser window.
- Remove old and unmaintained dependencies, replacing with maintained alternatives where required and removing altogether where not required.
- Add `install` command to easily install tool on remote machine over SSH.
- Add installer script for easy installation.
- Set up CI/CD with GitHub actions.
- Stop using IP and RPC and use Unix domain sockets with simple custom serialisation format instead.
- Remove feature for automatic line-ending translation.
- The 🍋 emojis have been replaced with 🍋‍🟩 emojis.

## Todo

- ~~Remove log15 package and use either standard log or maintained package.~~
- ~~Make compliant w/ golangci-lint.~~
- ~~Make tests not fail anymore.~~
- ~~Set up GH Actions workflows for linting, building and creating releases.~~
- ~~Move from ioutil to io/os or remove code altogether.~~
- ~~Update tested Go versions from 1.9-1.11 to 1.21-2.23.~~
- ~~Use `os.UserHomeDir()` instead of `mitchellh/go-homedir` or just remove altogether.~~
- ~~Migrate from `monochromgane/conflag` to `spf13/viper`.~~
- ~~Use posix compliant flags.~~
- ~~Use Unix socket instead of TCP and drop Windows support – `/var/run/limeade.sock`.~~
- ~~Drop RPC and use simple, custom binary encoding.~~
- Write lots of tests – unit and integration.
- Add simple script for installing correct version from GH releases without needing to build from source.

## Notes

The clipboard library we're using (atotto/clipboard) hasn't been updated is quite a while, but the other major library (golang-design/clipboard) requires CGO which is lame.

On Linux/Unix, this library has a runtime dependency on the executable used for clipboard interfacing. It seems to support [xsel](https://github.com/kfish/xsel), [xclip](https://github.com/astrand/xclip), [wl-clipboard](https://github.com/bugaevc/wl-clipboard) among others.

It'd be nice to use a supported library, but the golang-design/clipboard library still hasn't been updated in 2 years and doesn't support libwayland.

## Troubleshooting

### `Warning: remote port forwarding failed for listen path /tmp/limeade.sock`

This can happen if the server fails to listen on the desired socket path, e.g. if it is not allowed in the sshd config or the path is already in use.

First, check the [sshd config](https://man7.org/linux/man-pages/man5/sshd_config.5.html) (usually `/etc/ssh/sshd_config` and sometimes `/etc/ssh/sshd_config.d/*`) – if `AllowStreamLocalForwarding no` is set, you won't be able to forward any sockets. The default is `yes` so if it's not specified, it should work.

Second, check the sshd logs. It's not easy to do this for a production service so the best thing to do is start another sshd service on another port which you can safely run in debug mode without risking locking yourself out of your remote machine:

```shell
# On the remote machine – check your firewall to make sure you can access
# TCP on this port from your local machine.
/usr/sbin/sshd -d -p 2222
```

Then, in another terminal tab, on your local machine, connect to the remote machine on the new port:

```shell
ssh -p 2222 -R /tmp/limeade.sock:/tmp/limeade.sock user@host
```

You should see a spew of logs coming from the sshd process on the remote machine. If you see something like this:

```text
...
debug1: active: key options: agent-forwarding port-forwarding pty user-rc x11-forwarding
debug1: Entering interactive session for SSH2.
debug1: server_init_dispatch
debug1: server_input_global_request: rtype streamlocal-forward@openssh.com want_reply 1
debug1: server_input_global_request: streamlocal-forward listen path /tmp/limeade.sock
unix_listener: cannot bind to path /tmp/limeade.sock: Address already in use
debug1: server_input_channel_open: ctype session rchan 0 win 1048576 max 16384
debug1: input_session_request
debug1: channel 0: new [server-session]
debug1: session_new: session 0
debug1: session_open: channel 0
debug1: session_open: session 0: link with channel 0
...
```

Then this means that something is already binding to the socket file `/tmp/limeade.sock`. You can check what is using this file with `lsof`:

```shell
lsof /tmp/limeade.sock
```

Is that doesn't show anything, you might've just accidentally left the file there from a previous run. You can remove it with:

```shell
rm /tmp/limeade.sock
```

Then try again. You can also enable the sshd_config setting `StreamLocalBindUnlink yes` which will automatically remove the socket file immediately before binding to it, which will ensure no-one else is binding to the socket but will delete the socket file.