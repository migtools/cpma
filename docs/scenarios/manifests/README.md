## Generating manifests using CPMA and applying them

CPMA analyzes OCP3 configuration files and when possible transforms them to CRs, generated CRs are used to configure OCP4 cluster as `day two` operations.

---


### 1. Prerequisites

Prior to working with CPMA you need to deploy an OCP 3.7+ cluster(OCP 3.7, 3.9, 3.10, 3.11 are supported).
There are to 2 ways for cpma to obtain desired configurations:

1. Remotely retrieve configuration files using SSH, which are then stored locally (directory determined by `<workdir>` parameter) and under remote hostname sub directory. Proper SSH credentials (user, port and ssh key) are required with user having access to configuration files (sudo access might be necessary).

2. Locally stored configuration files. If option 1 is not possible, files can be retrieved prior to executing CPMA. The configuration files are expected to be located in `<workdir>` parameter and under the FQDN target hostname.

---

### 2. General CPMA configuration

CPMA can be configured using either:

1. Interactive prompts. This is simillar to `openshift-install`. The tool can be run with no configuration, all required values will be prompted. You can see an example below. This is the most recommended way, because prompts will guide you through needed values and it can generate a configuration based on prompted values that will be used later.

![prompt](https://user-images.githubusercontent.com/20123872/60581251-c0f57100-9d86-11e9-9ab3-7681b840731a.gif)


2. CLI parameters. All configuration values can be passed using CLI parameters. For example: `./cpma --source cluster.example.com --work-dir ./dir` Refer to CPMA's [README.md](https://github.com/fusor/cpma#usage) for full list of parameters.

3. Predefined configuration file. You can manually create a yaml configuration based on this [example](https://github.com/fusor/cpma/blob/master/examples/cpma-config.example.yaml). Configuration file path can be passed using `--config` parameter, or place `cpma.yaml` in your home directory.

4. Environmental variables. It is also possible to pass all configuration values as environmental variables. List of variables can be found in [README.md](https://github.com/fusor/cpma#e2e-tests)

---

### 3. Using CPMA to generate manifests

Once the configuration has been provided, by either prompt, CLI or ENV parameters, or configuration file, manifests will be generated and placed in `output-directory/manifests`

---

### 4. Applying generated CRs

For applying generated refer to [OCP 4 documentation](https://docs.openshift.com/container-platform/4.1/welcome/index.html). You can see an example of configuring OAuth below:

```bash
$ oc apply -f outputDirectory/100_CPMA-cluster-config-secret-htpasswd-secret.yaml
$ oc apply -f outputDirectory/100_CPMA-cluster-config-oauth.yaml
```

---

### 5. Complete scenario example

Using configuration files stored remotely:

![remote](https://user-images.githubusercontent.com/20123872/61694754-cff29200-ad3a-11e9-9254-de5f738e9c7d.gif)

Using configuration files stored locally:

![local](https://user-images.githubusercontent.com/20123872/61719959-9471bb00-ad6e-11e9-9d2a-59bfc223b1e5.gif)



