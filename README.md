# Limeade

[![Lint](https://github.com/drewsilcock/limeade/actions/workflows/lint.yaml/badge.svg)](https://github.com/drewsilcock/limeade/actions/workflows/lint.yaml)
[![Test](https://github.com/drewsilcock/limeade/actions/workflows/test.yaml/badge.svg)](https://github.com/drewsilcock/limeade/actions/workflows/test.yaml)

remote...lemote...lemode......Lemonade.........Limeade! 🍋‍🟩

Limeade is a tool for remote clipboard access, i.e. copy and paste over SSH.

It is a fork of [lemonade-command/lemonade](https://github.com/lemonade-command/lemonade).

## Installation

### Installer script

curl -sSfL https://raw.githubusercontent.com/drewsilcock/limeade/main/install.sh | sh -s

### Go install

```sh
go install github.com/drewsilcock/limeade/cmd@latest
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
$ ln -s /path/to/lemonade /usr/bin/pbcopy
$ echo "Hello" | pbcopy  # Same as 'echo "Hello" | lemonade copy'
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
- Set up CI/CD with GitHub actions and add installer script.
- Stop using TPC and RPC and use Unix sockets with simple custom serialisation format instead.
- Focus on Posix-compatible systems instead of Windows, removing option for automatic line-ending translation.
- The 🍋 emojis have been replaced with  🍋‍🟩 emojis.

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
- Write lots of tests.
- Add simple script for installing correct version from GH releases without needing to build from source.

## Notes

- The clipboard library we're using (atotto/clipboard) hasn't been updated is quite a while, but the other major library (golang-design/clipboard) requires CGO which is lame.