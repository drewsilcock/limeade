# Lemonade

remote...lemote...lemode......Lemonade.........Limeade! 🍋‍🟩

Lemonade is a remote utility tool.
(copy, paste and open browser) over TCP.

[![Build Status](https://travis-ci.org/lemonade-command/lemonade.svg?branch=master)](https://travis-ci.org/lemonade-command/lemonade)

## Installation

### Installer script

curl -sSfL https://raw.githubusercontent.com/drewsilcock/lemonade/main/install.sh | sh -s

### Go install

```sh
go install github.com/drewsilcock/lemonade-command@latest
```

### Build from source

```sh
git clone git@github.com:drewsilcock/lemonade
cd lemonade
make install
```

## Example of use

![Example](http://f.st-hatena.com/images/fotolife/P/Pocke/20150823/20150823173041.gif)

For example, you use a Linux as a virtual machine on Windows host.
You connect to Linux by SSH client(e.g. PuTTY).
When you want to copy text of a file on Linux to Windows, what do you do?
One solution is doing `cat file.txt` and drag displayed text.
But this answer is NOT elegant! Because your hand leaves from the keyboard to use the mouse.

Another solution is using the Lemonade.
You input `cat file.txt | lemonade copy`. Then, lemonade copies text of the file to clipboard of the Windows!

In addition to the above, lemonade supports pasting and opening URL.


## Usage

```sh
Usage: lemonade [options]... SUB_COMMAND [arg]
Sub Commands:
  open [URL]                  Open URL by browser
  copy [text]                 Copy text.
  paste                       Paste text.
  server                      Start lemonade server.

Options:
  --port=2489                 TCP port number
  --line-ending               Convert Line Ending(CR/CRLF)
  --allow="0.0.0.0/0,::/0"    Allow IP Range                [Server only]
  --host="localhost"          Destination hostname          [Client only]
  --no-fallback-messages      Do not show fallback messages [Client only]
  --trans-loopback=true       Translate loopback address    [open subcommand only]
  --trans-localfile=true      Translate local file path     [open subcommand only]
  --help                      Show this message
```

### On server (in the above, Windows)

```sh
$ lemonade server
```

### Client (in the above, Linux)

```sh
# You want to copy a text
$ cat file.txt | lemonade copy

# You want to paste a text from the clipboard of Windows
$ lemonade paste

# You want to open an URL to a browser on Windows.
$ lemonade open 'https://google.com'
```

## Configuration

You can override command line options by configuration file.
There is configuration file at `~/.config/lemonade.toml`.

### Server

```toml
port = 1234
allow = '192.168.0.0/24'
line-ending = 'crlf'
```

- `port` is a listening port of TCP.
- `allow` is a comma separated list of a allowed IP address(with CIDR block).


### Client

```toml
port = 1234
host = '192.168.x.x'
trans-loopback = true
trans-localfile = true
line-ending = 'crlf'
```

- `port` is a port of server.
- `host` is a hostname of server.
- `trans-loopback` is a flag of translation loopback address.
- `trans-localfile` is a flag of translation localfile.

Detail of `trans-loopback` and `trans-localfile` are described Advanced Usage.


##Advanced Usage

### trans-loopback

Default: true

This option works with `open` command only.

If this option is true, lemonade translates loopback address to address of client.

For example, you input `lemonade open 'http://127.0.0.1:8000'`.
If this option is false, server receives loopback address.
But this isn't expected.
Because, at server, loopback address is server itself.

If this option is true, server receives IP address of client.
So, server can open URL!

### trans-localfile

Default: true

This option works with `open` command only.

If this option is true, lemonade translates path of local file to address of client.

For example, you input `lemonade open ./file.txt`.
If this option is false, server receives `./file.txt`.
But this isn't expected.
Because, at server, `./file.txt` doesn't exist.

If this option is true, server receives IP address of client. And client serve the local file.
So, server can open the local file!

### line-ending

Default: "" (NONE)

This options works with `copy` and `paste` command only.

If this option is `lf` or `crlf`, lemonade converts the line ending of text to the specified.

### Alias

You can use lemonade as a `xdg-open`, `pbcopy` and `pbpaste`.

For example.

```sh
$ ln -s /path/to/lemonade /usr/bin/xdg-open
$ xdg-open 'http://example.com' # Same as lemonade open 'http://example.com'
```

### Usage over SSH

You can use Lemonade to provide copy and paste functionality over SSH using SSH port forwarding.

(Neovim will automatically detect the presence of the `lemonade` command and use it for copy and paste operations.)

```shell
lemonade install user@host
lemonade server --allow 127.0.0.1,::1 &

ssh -R 2489:127.0.0.1:2489 user@host
# You can now use lemonade copy and lemonade paste to copy and paste text
# between the local and remote machines.
```

You can also add the following to your `~/.ssh/config` file so that you don't need to specify the `-R` flag every time you SSH into the remote machine:

```shell
Host host
    User user
    HostName host
    ...
    RemoteForward 2489 localhost:2489
```

See:

- [SSH/OpenSSH/PortForwarding - Community Help Wiki](https://help.ubuntu.com/community/SSH/OpenSSH/PortForwarding)
- [WOW! and security? · Issue #14 · lemonade-command/lemonade](https://github.com/lemonade-command/lemonade/issues/14#)

## Limeade vs. Lemonade

This project is a fork of [lemonade-command/lemonade](https://github.com/lemonade-command/lemaonde) with a few key differences:

- Modernise codebase and keep dependencies up-to-date.
- Remove `open` functionality as this is unneeded and a security concern. You can just paste into your browser window.
- Remove old and unmaintained dependencies, replacing with maintained alternatives where required and removing altogether where not required.
- Add `install` command to easily install tool on remote machine over SSH.
- Set up CI/CD with GitHub actions and add installer script.
- Stop using TPC and RPC and use Unix sockets with Protobuf instead.
- Focus on Posix-compatible systems instead of Windows, removing option for automatic line-ending translation.
- The :lemon: emojis have been replaced with :lime: emojis.

## Links

- https://speakerdeck.com/pocke/remote-utility-tool-lemonade
- [リモートのPCのブラウザやクリップボードを操作するツール Lemonade を作った - pockestrap](http://pocke.hatenablog.com/entry/2015/07/04/235118)
- [リモートユーティリティーツール、Lemonade v0.2.0 をリリースした - pockestrap](http://pocke.hatenablog.com/entry/2015/08/23/221543)
- [lemonade v1.0.0をリリースした - pockestrap](http://pocke.hatenablog.com/entry/2016/04/19/233423)

## Todo

- Remove log15 package and use either standard log or maintained package.
- Make compliant w/ golangci-lint.
- Make tests not fail anymore.
- Set up GH Actions workflows for linting, building and creating releases.
- Add simple script for installing correct version from GH releases without needing to build from source.
- Move from ioutil to io/os.
- Update tested Go versions from 1.9-1.11 to 1.21-2.23.
- Use `os.UserHomeDir()` instead of `mitchellh/go-homedir`.
- Migrate from `monochromgane/conflag` to `spf13/viper`.
- Use posix compliant flags.
- Consider moving pocke/go-iprange to just be a subpackage of this repo.
- Use Unix socket instead of TCP and drop Windows support – `/var/run/limeade.sock`.
- Drop RPC and use simple, custom binary encoding.

## Notes

- The clipboard library we're using (atotto/clipboard) hasn't been updated is quite a while, but the other major library (golang-design/clipboard) requires CGO which is lame.