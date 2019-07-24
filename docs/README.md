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
  * Crio
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
  * Crio
  * Cluster Resources Quotas and Resources Quotas
  * Image
  * OAuth
  * Project
  * Scheduler
  * SDN

### See Also

[Master configuration supported items](Supported.md)
