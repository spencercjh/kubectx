# `sshctx`: Power tool for `ssh`

This repository provides `sshctx` tool.

**`sshctx`** helps you switch between hosts back and forth.

# sshctx(1)

sshctx is a utility to switch between ssh(1) hosts.

```
USAGE:
  kubectx                   : list the contexts
  kubectx <NAME>            : switch to context <NAME>
  kubectx -                 : switch to the previous context
  kubectx -c, --current     : show the current context name
  kubectx <NEW_NAME>=<NAME> : rename context <NAME> to <NEW_NAME>
  kubectx <NEW_NAME>=.      : rename current-context to <NEW_NAME>
  kubectx -d <NAME>         : delete context <NAME> ('.' for current-context)
                              (this command won't delete the user/cluster entry
                              that is used by the context)
  kubectx -u, --unset       : unset the current context
```

### Usage

```sh
$ sshctx 192.168.1.2
Connect to context `192.168.1.2`.

$ sshctx -
Connect to context `192.168.1.2`.
```

`sshctx` supports <kbd>Tab</kbd> completion on bash/zsh/fish shells to help with long host names. You don't have to
remember full host names anymore.

-----
