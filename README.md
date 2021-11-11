# `sshctx`: Power tool for `ssh`

This repository provides `sshctx` tool.

**`sshctx`** helps you switch between hosts back and forth.

[Install &rarr;](#Installation)

# sshctx(1)

sshctx is a utility to switch between ssh(1) hosts.

```
USAGE:
  sshctx                       : list the hosts
  sshctx <HOST>                : connect to <HOST>
  sshctx -                     : connect to the previous successfully connected host
  sshctx -p, --previous        : show the previous successfully connected host
  sshctx -h,--help             : show this message
  sshctx -v,-V,--version       : show version
```

### Usage

```sh
$ sshctx
List hosts in `~/.ssh/config`.

$ sshctx test
Connect to host `test` in your `~/.ssh/config`.

$ sshctx -
Connect to context `test`.

$ sshctx -p
Show the latest connected host
```

-----

## Installation

### Macos

Distribution with Homebrew: [Formula](https://github.com/spencercjh/homebrew-sshctx)

```shell
brew tap spencercjh/sshctx

brew install sshctx

brew upgrade sshctx
```

### Manual Installation (macOS and Linux)

Example installation steps:

```shell
# get source code, no matter where
git clone https://github.com/spencercjh/sshctx
# go to src dir
cd sshctx/cmd/sshctx
# build
go build
# make it executable
chmod +x ./sshctx
# link to path
sudo ln -s $PWD/sshctx /usr/local/bin/sshctx
```

-----

## Interactive mode

If you want `sshctx` command to present you an interactive menu with fuzzy searching, you just need
to [install `fzf`](https://github.com/junegunn/fzf) in your `$PATH`.

If you have `fzf` installed, but want to opt out of using this feature, set the environment
variable `SSHCTX_IGNORE_FZF=1`.

If `fzf` isn't installed, but you still want to use `sshctx` in interactive
mode, [promptui](https://github.com/manifoldco/promptui) provides a simple interactive mode without fuzzy searching.

If you want to keep `fzf` interactive mode but need the default behavior of the command, you can do it using Unix
composability:

```
sshctx | cat
```

-----
