# cpma [![Build Status](https://travis-ci.com/fusor/cpma.svg?branch=master)](https://travis-ci.com/fusor/cpma) [![Maintainability](https://api.codeclimate.com/v1/badges/aac7d46fd7899042ce52/maintainability)](https://codeclimate.com/github/fusor/cpma/maintainability)
Control Plane Migration Assistance:  Intended to help migration cluster configuration of a OCP 3.x cluster to OCP 4.x


# Build

Tested on go 1.12.4

This project is Go Modules for managing dependencies. This means it can be cloned and compiled outside of GOPATH.

```console
$ git checkout https://github.com/fusor/cpma.git
$ cd cpma
$ make
$ ./bin/cpma
```

# Usage

Flags:
```
--config string           config file (default is $HOME/.cpma.yaml)
--debug                   show debug ouput
-h, --help                help for cpma
-o, --output-dir string   set the directory to store extracted configuration.
```

You can find example config in `examples/` 

Example:

```console
$ ./bin/cpma --config /path/to/config/.yml --debug
```

# Unit tests

In order to add new unit test bundle create `*_test.go` file in package you want to test(ex: `foo.go`, `foo_test.go`).
To execute tests run `make test`.

https://golang.org/pkg/testing/
