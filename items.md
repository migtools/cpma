| Component | OCP3 | OCP4 | Manifests | Reported | OCP4 support |
| :--- | :--- | :--- | :---: | :---: | :--- |
| Authentication and Authorization Configuration | authConfig | Incompatible | No | No | |
| Authentication and Authorization Configuration | AuthenticationCacheSize  | Incompatible | No | No | |
| Authentication and Authorization Configuration | AuthorizationCacheTTL | Incompatible | No | No | |
| etcd Configuration | Address | Future | No | No | >= OCP4.4 |
| etcd Configuration | etcdClientInfo | Future | No | No | >= OCP4.4 |
| etcd Configuration | etcdConfig | Future | No | No | >= OCP4.4 |
| etcd Configuration | etcdStorageConfig | Future | No | No | >= OCP4.4 |
| etcd Configuration | KubernetesStoragePrefix  | Future | No | No | >= OCP4.4 |
| etcd Configuration | KubernetesStorageVersion | Future | No | No | >= OCP4.4 |
| etcd Configuration | OpenShiftStoragePrefix | Future | No | No | >= OCP4.4 |
| etcd Configuration | OpenShiftStorageVersion  | Future | No | No | >= OCP4.4 |
| etcd Configuration | PeerAddress | Future | No | No | >= OCP4.4 |
| etcd Configuration | PeerServingInfo | Future | No | No | >= OCP4.4 |
| etcd Configuration | ServingInfo | Future | No | No | >= OCP4.4 |
| etcd Configuration | StorageDir | Future | No | No | >= OCP4.4 |
| Image Policy Configuration | DisableScheduledImport | No | No | Yes  | OCP4: Always enabled |
| Image Policy Configuration | MaxImagesBulkImportedPerRepository  | No | No | Yes  | OCP4: no limit |
| Image Policy Configuration | MaxScheduledImageImportsPerMinute | No | No | Yes  | OCP4: 60 per minute |
| Image Policy Configuration | ScheduledImageImportMinimumIntervalSeconds | No | No | Yes  | OCP4: 15 minutes |
| Image Policy Configuration | AllowedRegistriesForImport | Yes | Yes | Yes  | List (DomainName \| Insecure) |
| Image Policy Configuration | AdditionalTrustedCA | No | No | Yes  | |
| Image Policy Configuration | InternalRegistryHostname | No | No | Yes  | OCP4 Configured via registry operator |
| Image Policy Configuration | ExternalRegistryHostname | Yes | Yes | Yes  | |
| Network Configuration | ClusterNetworkCIDR | Yes | Yes | Yes  | |
| Network Configuration | externalIPNetworkCIDRs | No | No | Yes  | |
| Network Configuration | ingressIPNetworkCIDR  | Yes | No | No | >= OCP4.4 |
| Network Configuration | HostSubnetLength  | No | No | Yes  | |
| Network Configuration | NetworkPluginName | Yes | Yes | Yes  | |
| Network Configuration | serviceNetworkCIDR | Yes | Yes | Yes  | |
| OAuth Authentication Configuration | AlwaysShowProviderSelection  | No | No | Yes  | |
| OAuth Authentication Configuration | AssetPublicURL | No | No | Yes  | |
| OAuth Authentication Configuration | Template:IdentityProviders | Yes | Yes | Yes  | OAuth CRD:spec:identityProviders |
| OAuth Authentication Configuration | Template:ProviderSelection | Yes | No | No | OAuth CRD:spec:template:providerSelection:name |
| OAuth Authentication Configuration | Template:Login | Yes | No | No | OAuth CRD:spec:template:login:name |
| OAuth Authentication Configuration | Template:Error | Yes | No | No | OAuth CRD:spec:template:error:name |
| OAuth Authentication Configuration | MasterCA | No | No | Yes  | |
| OAuth Authentication Configuration | MasterPublicURL| No | No | Yes  | |
| OAuth Authentication Configuration | MasterURL  | No | No | Yes  | |
| OAuth Authentication Configuration | grantConfig | No | No | Yes  | The method must now be specified by OAuth Client (grantMethod) |
| OAuth Authentication Configuration | SessionConfig:sessionMaxAgeSeconds  | No | No | Yes  | |
| OAuth Authentication Configuration | SessionConfig:sessionName | No | No | Yes  | |
| OAuth Authentication Configuration | SessionConfig:sessionSecretsFile | No | No | Yes  | |
| OAuth Authentication Configuration | TokenConfig:accessTokenMaxAgeSeconds | Yes | Yes | Yes  | OAuth CRD:spec:tokenConfig:accessTokenMaxAgeSeconds |
| OAuth Authentication Configuration | TokenConfig:accessTokenMaxAgeSeconds | Incompatible | No | Yes  | Hard coded: 5min |
| Project Configuration | DefaultNodeSelector | No | No | Yes  | |
| Project Configuration | SecurityAllocator:mcsAllocatorRange | No | No | Yes  | |
| Project Configuration | SecurityAllocator:mcsLabelsPerProject | No | No | Yes  | |
| Project Configuration | SecurityAllocator:uidAllocatorRange | No | No | Yes  | |
| Project Configuration | ProjectRequestMessage | Yes | Yes | Yes  | |
| Project Configuration | ProjectRequestTemplate | Yes | Yes | Yes  | |
| Scheduler Configuration | SchedulerConfigFile | Yes | No | No | See ConfigMap default scheduler; Default policy applies if not defined |
| Service Account Configuration | LimitSecretReferences | Incompatible | No | No | |
| Service Account Configuration | ManagedNames| Incompatible | No | No | |
| Service Account Configuration | MasterCA | Incompatible | No | No | |
| Service Account Configuration | PrivateKeyFile | Incompatible | No | No | |
| Service Account Configuration | PublicKeyFiles | Incompatible | No | No | |
| Service Account Configuration | ServiceAccountConfig  | Incompatible | No | No | |
| Specifying TLS ciphers for etcd | | Incompatible | No | No | |
