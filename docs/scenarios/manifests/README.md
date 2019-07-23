## Generating manifests using CPMA and applying them

CPMA analyzes OCP3 configuration files and when possible transforms them to CRs, generated CRs are used to configure OCP4 cluster as `day two` operations.

---


### 1. Prerequisites

Prior to working with CPMA you need to deploy an OCP 3.7+ cluster(OCP 3.7, 3.9, 3.10, 3.11 will be supported).
There are to 2 ways for cpma to obtain desired configurations:

1. SSH, CPMA is capable to connect to cluster using SSH and fetch configurations to a local directory. For using SSH it's required to have a proper private key and SSH user may require sudo permissions for reading configurations.

2. Configurations stored in local directory, it's possible to manualy download configuration and feed them to CPMA.

---

### 2. General CPMA configuration

CPMA can be configured using either:

1. Interactive prompts. This is simillar to `openshift-install`. The tool can be run with no configuration, all missing values will be prompted. You can see an example below. This is the most recommended way, because prompts will guide you through needed values and it can generate a configuration based on prompted values that will be used later.

![prompt](https://user-images.githubusercontent.com/20123872/60581251-c0f57100-9d86-11e9-9ab3-7681b840731a.gif)


2. Flags. All configuration values can be passed using flags. For example: `./cpma --source cluster.example.com --work-dir ./dir` Refer to CPMA's [README.md](https://github.com/fusor/cpma#usage) for full list of flags.

3. Predefined configuration file. You can manually create a yaml configuration based on this [example](https://github.com/fusor/cpma/blob/master/examples/cpma-config.example.yaml).

4. Environmental variables. It is also possible to pass all configuration values as environmental variables. List of variables can be found in [README.md](https://github.com/fusor/cpma#e2e-tests)

---

### 3. Using CPMA to generate manifests

Once the configuration has been provided, either by prompt or by providing configuration file, manifests will be generated and placed in `output-directory/manifests`

If you have a predifined yaml configuration, you can pass it using --config flag, or place cpma.yaml in the directory from which you run cpma.

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

WIP


