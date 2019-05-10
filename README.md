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

# IO

The data file structure looks like the following tree structure example.
The cluster endpoints subfolders contain the configuration files retrieved and to process.
The manifests directory contains the generated CRDs.

```
data
├── manifests
├── master-0.example.com
|   └── etc
|       └── origin
|           ├── master
|               ├── htpasswd
|               └── master-config.yaml
└── node-1.example.com
    └── etc
        └── origin
            └── node
                └── node-config.yaml
```

The configuration files are retrieved from local disk (outputDir/<Hostname>/),
If a file is not available it's retrieved from <Hostname> and stored on local disk.

To trigger a total or partial network file fetch, remove any prior data from <Hostname> sub directory.


# Unit tests

In order to add new unit test bundle create `*_test.go` file in package you want to test(ex: `foo.go`, `foo_test.go`).
To execute tests run `make test`.

https://golang.org/pkg/testing/
