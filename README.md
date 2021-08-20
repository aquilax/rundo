# rundo

Local scripts runner

Run commands defined in .rundo file in the current directory

```shell
$cat .rundo
list=ls -lah
```

## Usage

```shell
$ rundo --help
NAME:
   rundo - Local scripts runner

USAGE:
   rundo [global options] command [command options] [arguments...]

COMMANDS:
   list, ls  Shows list of the defined commands in the current context
   run, r    Runs one of the commands
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```
