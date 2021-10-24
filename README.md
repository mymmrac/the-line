# The Line

Simple extendable line counter

## Installation

```shell
go install github.com/mymmrac/the-line@latest

# Optionally rename binary 
mv $GOPATH/bin/the-line $GOPATH/bin/tl
```

> NOTE: Make sure to add `$GOPATH/bin` to `$PATH` to use

## Usage

```shell
# Search for all files/folders in current dir and count lines in Go files
tl -p go -r .

# Help info
tl --help
```

## Features

List of supported features (or TBD features)

### Filters:

- [X] Match
- [X] Prefix/suffix
- [X] Contains
- [X] Length
- [X] Regexp
- [X] Blank
- [X] Any
- [X] Union
- [X] Intersection
- [X] Not

### Modifiers:

- [X] Trim spaces/prefix/suffix
- [X] To lower

### Configs:

- [X] Embedded (default config included in binary)
- [X] YAML (user defined config)

### CLI:

- [X] Pretty interface using [Charm](https://charm.sh/)
- [X] Minimal output (regular stdout)
- [X] Output only total
- [ ] Output grouped by folders
- [ ] Show used/matched by profiles/skipped files

### Refactoring

- [ ] Main
- [ ] Display
- [ ] Make packages