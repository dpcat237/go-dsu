## Overview

Go DSU - provides tools to update Go dependencies with more control than default Go tools. It isn't stable yet. Use it with caution, and any feedback would be very appreciated.  
Go DSU works on base of Go modules and Git. Git uses public SSH key for access to private repositories.  

### Implement for first stable version
- [x] Simple update of direct / indirect modules
- [x] Return table with available updates
- [x] Allow select modules to update
- [x] Optionally run local tests before and after update of each module with rollback if tests fail
- [x] Check if license of direct and indirect dependencies changed
- [x] Check vulnerabilities from OSS Index
- [x] Adapt download with Git for not a known dependency
- [x] Current command (analyze) - list all current licenses / vulnerabilities of dependencies
- [ ] Implement more tests and improve documentation

## Installation

    go get github.com/dpcat237/go-dsu

## Usage

```
$ go-dsu
Go DSU - provides tools to update Go dependencies with more control than default Go modules.

Usage:
  go-dsu [command]

Available Commands:
  clean       Clean modules
  help        Help about any command
  preview     Preview updates
  update      Update modules

Flags:
  -h, --help   help for go-dsu

Use "go-dsu [command] --help" for more information about a command.
```

## Examples

### Preview available updates with changes
`$ go-dsu preview`

![](doc/images/preview.png)