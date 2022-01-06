DPhoto - CLI
====================================

Installed on the end-user computer, backup photos and videos using command line interface.

Getting Started
------------------------------------

Usage is available with `help` command:

```
$ dphoto help
Backup photos and videos to your personal AWS Cloud.

Usage:
  dphoto [command]

Available Commands:
  album        Organise your collection into albums
  backup       Backup photos and videos to personal cloud
  configure    Configuration wizard to grant dphoto access AWS resources.
  help         Help about any command
  housekeeping Run housekeeping script to perform delayed operations
  scan         Discover directory structure to suggest new albums to create
  version      Print the version

Flags:
      --config string   use configuration file provided instead of searching in ./ , $HOME/.dphoto, and /etc/dphoto
      --debug           enable debug logging
  -h, --help            help for dphoto

Use "dphoto [command] --help" for more information about a command.
```

Contribute
------------------------------------

Test and build:

    go test ./...
    go build

    # or
    make