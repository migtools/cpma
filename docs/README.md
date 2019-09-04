## CPMA capabilities

This documents list CPMA capabilities.

---

### Report

Report consist of 2 parts:

1. Pre-migration analytics that goes under "cluster" in json. It contains general information about:
  * Cluster Resources Quotas - name, spec, selectors(labels, annotations)
  * Nodes - name, сpu usage, memory capacity, consumed memory, number of running pods, pod capacity, information that
  indicates if node is master.
  * Namespaces - name, latest change timestamp, container count, total cpu/memory usage.
    * Pods - name
    * Route - name, spec
    * Daemonsets - name, latest change timestamp
    * Deployments - name, latest change timestampg
    * Resource Quotas - name, spec, selectors(labels, annotations)
  * Storage classes - name, provisioner
  * RBAC - information about users, groups, roles, cluster roles, cluster role bindings, security context constraints
  * Persistent volumes - names, storage class, driver, capacity, phase

2. Information about configurations that indicates what can/can't be migrated. Following configurations are included:
  * API
  * CRI-O
  * Docker
  * Etcd
  * Image
  * OAuth
  * Project
  * Scheduler
  * SDN

---

### Manifests

List of supported configuration to manifest translations:
  * API Certificate
    * If API/HTTP TLS Certificate defined in Master configuration file under ServingInfo section is not signed by openshif itself (therefore changed by user for a proper CA signed certificate) then it's ported to a TLS secret and saved under the '100_CPMA-cluster-config-APISecret.yaml' file. To apply the secret and update the API server, follow this [procedure](https://docs.openshift.com/container-platform/4.1/authentication/certificates/api-server.html#add-named-api-server_api-server-certificates).
  * CRI-O
    * If defined in OCP3's cluster, the CRI-O configuration defined in 'crio.conf' is ported to a machineconfiguration.openshift.io resource and saved under '100_CPMA-crio-config.yaml'.
  * Cluster Resources Quotas
    * Every quota defined at the cluster level is exported into equivalent OCP4 CR file. The file name is in the form: '100_CPMA-cluster-quota-resource-<ClusterQuota name>.yaml'
  * Resources Quotas
    * Within a given namespace/project, each quota resource is exported into equivalent OCP4 CR file in the form of '100_CPMA-<Namespace>-resource-quota-<ResourceQuota Name>.yaml'
  * Image Configuration
    * An image.config.openshift.io resource, saved under file name '100_CPMA-cluster-config-image.yaml' is created from the following sources:
      * Portable Image registries information from OCP 3 master file /etc/registries/registries.conf.
      * Portable Image Policy Configuration information from OCP 3 master configuration file etc/origin/master/master-config.yaml.
  * OAuth Providers
    * All OAuth providers defined in OCP 3 are ported to OCP4 as an OAuth resource CR file 100_CPMA-cluster-config-oauth.yaml.
  * Projects Configuration
    * Existing project configuration information that are portable are created in the projects.config.openshift.io resource file 100_CPMA-cluster-config-project.yaml.
  * Scheduler
    * The initial OCP 3 master's scheduler configuration is ported as an OCP 4 operator. The generated CR file is dubbed '100_CPMA-cluster-config-scheduler.yaml'.
  * SDN
    * Network configuration from the OCP's master is ported to OCP 4 as a network operator. The CR file is saved under '100_CPMA-cluster-config-sdn.yaml'.

### See Also

[Master configuration supported items](Supported.md)
