#### WIP - project it's in early stage of development. Use it with caution.

## Overview

Go DSU - provides tools to update Go dependencies with more control than default Go modules.

### Implemented
- [x] Simple updater
- [ ] Check internet connection before starting a process
- [ ] Return table with available updates
- [ ] Run local tests before and after update of each dependency with rollback if tests fail
- [ ] Run tests of dependencies
- [ ] Check if license of dependencies changed
- [ ] Check if in dependency changed contract of implemented interface / method

## Installation

**Ensure Go modules are enabled: GO111MODULE=on and go/bin is in your PATH variable.**

    go get github.com/dpcat237/go-dsu

## Usage

```
$ go run main.go             
Go DSU - provides tools to update Go dependencies with more control than default Go modules.

Usage:
  go-dsu [command]

Available Commands:
  help        Help about any command
  update      Update dependencies

Flags:
  -h, --help   help for go-dsu

Use "go-dsu [command] --help" for more information about a command.
```