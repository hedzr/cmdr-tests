# Cmdr CLI

## Getting Started

```bash
go install github.com/hedzr/cmdr-cli/cli/cmdr
cmdr --help
cmdr new
```

### Create a new app

```bash
$ cmdr new myapp
$ cd myapp
$ cmdr new cmd jump.to.here
$ cmdr new flg jump.to.here count int 0
$ cmdr new cmd <<<EOT
jump.to
jump.to.{here,there}
server.{start,stop,restart,status,reload}
EOT
```
