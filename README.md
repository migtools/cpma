## cpma [![Build Status](https://travis-ci.com/fusor/cpma.svg?branch=master)](https://travis-ci.com/fusor/cpma) [![Maintainability](https://api.codeclimate.com/v1/badges/aac7d46fd7899042ce52/maintainability)](https://codeclimate.com/github/fusor/cpma/maintainability)
Control Plane Migration Assistant (CPMA) is a Command Line interface to help as much as possible users migrating an Openshift 3.7+ control plane configuration to an Openshift 4.x.
The utility provides Custom Ressource (CR) manifests and reports informing users which aspects of configuration can and cannot be migrated.

## Introduction

### Background

Openshift Container Platform (OCP) version 3 uses Ansible's openshift-ansible modules to allow for extensive configuration.

Meanwhile OCP 4.x uses openshift-installer which integrates with supported clouds for deployment, openshift-installer is a Day-1 operation and relies on [Operators](https://www.openshift.com/learn/topics/operators) to install and configure the cluster.

Many of the installation options available during install time in Openshift 3 are configurable by the user as Day-2 operations in Openshift 4.
In many cases this is done by writing CRs for operators which affect the configuration changes.

Since there is no direct upgrade path from OCP 3 to OCP 4, the goal of CPMA is to ease the transition between the 2 OCP versions involved.

### Goals

Bring the target Cluster to be as close as possible from the source cluster with the use of CR Manifests and report which configuration aspects can and cannot be migrated.
The tool effectively provides confidence information about each processed options to explain what is supported, fully, partially or not.
The generated CR Manifests can be applied to an Openshift 4 cluster as Day-2 operations.

### Non-Goals

Applying the CRs to the Openshift 4 cluster directly is currently not an intended feature for this utility.
The user must review the output CRs and filter the desired ones to be applied to a targeted OCP 4 cluster.
Migrating workloads are not to be done with CPMA, please use the appropriate migration tool for this purpose.

### Operation Overview

The purpose of the CPMA tool is to assist an administrator to migrate the control plane of an OpenShift cluster from version 3.7+ to its next major OpenShift version 4.x.
For that purpose CPMA sources information from:
- Cluster Configuration files:
  - Master Node:
    - Master configuration file - Usually /etc/origin/master/master-config.yaml
    - CRIO configuration file - Usually /etc/crio/crio.conf
    - ETC configuration file - Usually /etc/etcd/etcd.conf
    - Image Registries file - Usually /etc/containers/registries.conf
    - Dependent configuration files:
      - Password files
        - HTPasswd, etc.
      - Configmaps
      - Secrets
- APIs
  - Kubernetes
  - Openshift

Configuration files are processed to generate equivalent Custom Resource manifest files which can then to be consumed by OCP 4.x Operators.
During that processs, every parameter that is analysed is ported when compatible to its equivalent.
A feature fully supported or not means there is a direct or not equivalent in OCP 4.
A partially supported parameter indicates the feature is note entirelly equivalent.
The reason for the latter is because some features are deprecated or used differently in OCP 4.
OCP 3 and 4 approach configuration management completely differently across.
Therefore it’s expected the tool cannot port all features across.

For more information about CPMA coverage please see [docs](./docs).

CPMA uses an ETL pattern to process the configuration files and query APIs which produce the output in two different forms:
- Custom Resource Manifest files in YAML format
- A report file (by default report.json) is produced by the reporting process

The user can then review the new configuration from the generated manifests and must aslo use the reports as a guide how to leverage the CRs to apply configuration to a newly installed Openshift 4 cluster. The reviewed configuration can then be used to update a targeted OpenShift cluster.

## Prerequisites
Prior to the execution of the assistant tool, the following must be met:

* The OCP source cluster must have be updated to the latest asynchronous release.
* The environment health check must have been executed to confirm there are no diagnostic errors or warnings.
* The OCP source cluster must meet all current prerequisites for the given version of OCP.
* The OCP source cluster must be at least of any of the versions 3.7, 3.9, 3.10 or 3.11.

## Warning

The CPMA tool is to assist with complex settings and therefore it’s mandatory to have read the corresponding documentation.
For more information please refef to:
- OCP 3.x documentation
- OCP 4.x documentation

## Installation

### Build from source

Requires go >= v1.11

This project uses Go Modules for managing dependencies. This means it can be
cloned and compiled outside of `GOPATH`.

```console
$ git clone https://github.com/fusor/cpma.git
$ cd cpma
$ make
$ ./bin/cpma
```
### CPMA Image

CPMA is also available in as an image, quay.io/ocpmigrate/cpma

Example Usage:
```
docker run -it --rm -v ${PWD}:/mnt:z -v $HOME/.kube:/.kube:z -v $HOME/.ssh:/.ssh:z -u ${UID} \
quay.io/ocpmigrate/cpma:latest
```

Where `${PWD}` is mounted in the working directory of the image (`/mnt`). This means that paths provided to --config and --work-dir will need to specified be relative to your present working directory.

To make it a little more intuitive it can also be run via an alias, for example:
```console
$ alias cpma="docker run -it --rm -v ${PWD}:/mnt:z -v $HOME/.kube:/.kube:z -v $HOME/.ssh:/.ssh:z \
-u ${UID} quay.io/ocpmigrate/cpma:latest"
```

## Getting Started

CPAM needs information about the source cluster to anaylyse.

CPMA configuration information can be provided in following ways:
- Environment (ENV) variables
- Command Line (CLI) parameters
- User prompt
- Configuration file

### ENV
CPMA specific environment variable must be prefixed with `CPMA_`:
- CPMA_CONFIGSOURCE
- CPMA_CLUSTERNAME
- CPMA_CRIOCONFIGFILE
- CPMA_DEBUG
- CPMA_ETCDCONFIGFILE
- CPMA_HOSTNAME
- CPMA_INSECUREHOSTKEY
- CPMA_NODECONFIGFILE
- CPMA_MANIFESTS
- CPMA_MASTERCONFIGFILE
- CPMA_REGISTRIESCONFIGFILE
- CPMA_REPORTING
- CPMA_SSHPRIVATEKEY
- CPMA_SSHLOGIN
- CPMA_SSHPORT
- CPMA_VERBOSE
- CPMA_WORKDIR

### CLI
```console
$ ./bin/cpma -h
Usage:
  cpma [flags]

Flags:
  -i, --allow-insecure-host        allow insecure ssh host key
  -c, --cluster-name string        OCP3 cluster kubeconfig name
      --config string              config file (Default searches ./cpma.yaml, $HOME/cpma.yml)
      --config-source string       source for OCP3 config files, accepted values: remote or local
      --crio-config string         path to crio config file
  -d, --debug                      show debug ouput
      --etcd-config string         path to etcd config file
  -h, --help                       help for cpma
  -n, --hostname string            OCP3 cluster hostname
  -m, --manifests                  Generate manifests (default true)
      --master-config string       path to master config file
      --node-config string         path to node config file
      --registries-config string   path to registries config file
  -r, --reporting                  Generate reporting  (default true)
  -k, --ssh-keyfile string         OCP3 ssh keyfile path
  -l, --ssh-login string           OCP3 ssh login
  -p, --ssh-port int16             OCP3 ssh port
  -v, --verbose                    verbose output
  -w, --work-dir string            set application data working directory (Default ".")
```

Example:
```console
$ ./bin/cpma --config /path/to/config/.yml --verbose --debug
```

### User Prompt
The user will be prompted for required parameters

### Configuration file
The default CPMA configuration file is either `./cpma.json` or `~/cpma.json` unless an explicit path/name is provided (see --config option)

```json
clustername: openshift-testuser-example-com
configsource: remote
crioconfigfile: /etc/crio/crio.conf
debug: false
etcdconfigfile: /etc/etcd/etcd.conf
fetchfromremote: true
home: /home/testuser
hostname: master0.example.com
insecurehostkey: false
manifests: true
masterconfigfile: /etc/origin/master/master-config.yaml
nodeconfigfile: /etc/origin/node/node-config.yaml
registriesconfigfile: /etc/containers/registries.conf
reporting: true
saveconfig: true
sshlogin: testuser
sshport: 0
sshprivatekey: /home/users/test/.ssh/testuser
verbose: true
workdir: data

```

### Data file structure

CPMA data file structure looks like the following example.

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
|   └── etc
|       └── origin
|           └── node
|               └── node-config.yaml
└── report.json
```

The `WorkDir` option ('data' in above example) determines the directory containing all data handled by CPMA.

The `manifests` subfolder will contain all generated Custom Resource manifests.

Each cluster endpoint is identified by its FQDN (master-0.example.com in above example) and its subfolders contain the configuration files retrieved and to
process.

And finally `report.json` contains the result of the reporting analysis.

### Local or remote modes

CPMA can be used in either remote or local modes.
- In remote mode, CPMA retrieves the configuration files itself using SSH, which are then stored locally.
- Conversely, in local mode, configuration files must have been copied locally prior to launching CPMA.

Whether the files are retrieved by CPMA (remote mode) or manually (local mode), the tool always relies on the local file system to process a cluster node files using `<workDir>/<Hostname>/`.

## Contribute

### Unit tests and integration tests

In order to add new unit test bundle create `*_test.go` file in package you
want to test(ex: `foo.go`, `foo_test.go`).  To execute tests run `make test`.

https://golang.org/pkg/testing/

### Functional tests

Tests are located under `tests` directory, based on project layout [standards](https://github.com/golang-standards/project-layout).

### E2E tests

Prerequisites for base cluster report e2e tests:

| ENV variable            |   Expected value                                                                                                  |
|-------------------------|-------------------------------------------------------------------------------------------------------------------|
|    `CPMA_PASSWD`        |    Password for the cluster                                                                           |
|    `CPMA_LOGIN`         |    Login for the cluster                                                                                          |
|    `CPMA_HOSTNAME`      |    Hostname of the stable cluster                                                                                 |
|    `CPMA_CLUSTERNAME`   |    Hostname from the current context of the stable cluster (could be extracted with `oc whoami --show-context`)   |
|    `CPMA_SSHLOGIN`      |    SSH login for the cluster                                                                                      |
|    `CPMA_SSHPORT`       |    SSH port                                                                                                       |
|    `CPMA_SSHPRIVATEKEY` |    Path for the private key                                                                                       |

The workflow involves:
1) Ensuring the cluster session is open.
2) Generating cluster report from cluster.
3) Parsing generated and local reports.
4) Comparison of reports.

Those are needed as a substitution of configuration or CLI input, so the test could run autonomously.
