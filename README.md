# `sshctx`: Power tool for `ssh`

This repository provides `sshctx` tool.

**`sshctx`** helps you switch between hosts back and forth.

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
List hosts in `~/.ssh/config`

$ sshctx 192.168.1.2
Connect to context `192.168.1.2`.

$ sshctx -
Connect to context `192.168.1.2`.
```

-----
