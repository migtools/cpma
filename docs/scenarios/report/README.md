## Generating pre-migration report with CPMA

---

### 1. Prerequisites

Prior to working with CPMA you need to deploy an OCP3 cluster(OCP 3.7, 3.9, 3.10, 3.11 should be supported).
In order to generate a report, CPMA interacts with OCP using it's API. This means KUBECONFIG is required as a well as user used to talk with OCP cluster API needs to have priviligies to list nodes, projects, pods etc. Recommended cluster role is `system:admin`. It can be configured with following command: `oc adm policy add-cluster-role-to-user system:admin <username>`.

---

### 2. General CPMA configuration

CPMA can be configured in a couple different ways.

1. Interactive prompts. This is simillar to `openshift-install`. The tool can be run with no config, all missing values will be prompted. You can see an example below. This is the most recommended way, because prompts will guide you through needed values and it can generate a config based on prompted values that will be used later.

![prompt](https://user-images.githubusercontent.com/20123872/60581251-c0f57100-9d86-11e9-9ab3-7681b840731a.gif)


2. Flags. All config values can be passed using flags. For example: `./cpma --source cluster.example.com --work-dir ./dir` Refer to CPMA's [README.md](https://github.com/fusor/cpma#usage) for full list of flags.

3. Predefined config file. You can manually create a yaml config based on this [example](https://github.com/fusor/cpma/blob/master/examples/cpma-config.example.yaml).

4. Environmental variables. It is also possible to pass all configuration values as environmental variables. List of variables can be found in [README.md](https://github.com/fusor/cpma#e2e-tests)

---

### 3. Using CPMA to generate report

In case you decided to configure CPMA using interactive prompt, cluster report will be generated right after prompting all values.

If you have a predifined yaml config, you can pass it using `--config` flag, or place `cpma.yaml` in the directory from which you run cpma.

---

## 4. Reading report.json

Generated report will be placed inside specified working directory in format of a json file. We are still working on passing this report to UI.

You can find report example report in this scenario directory.