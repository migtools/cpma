package decode

import (
	"io/ioutil"
	"testing"

	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestMasterConfig(t *testing.T) {
	t.Parallel()

	dynamicProvisioningEnabled := true
	masterCA := "ca-bundle.crt"
	expectedMasterConfig := &legacyconfigv1.MasterConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MasterConfig",
			APIVersion: "v1",
		},
		ServingInfo: legacyconfigv1.HTTPServingInfo{
			ServingInfo: legacyconfigv1.ServingInfo{
				BindAddress: "0.0.0.0:443",
				BindNetwork: "tcp4",
				CertInfo: legacyconfigv1.CertInfo{
					CertFile: "master.server.crt",
					KeyFile:  "master.server.key",
				},
				ClientCA:          "ca.crt",
				NamedCertificates: []legacyconfigv1.NamedCertificate(nil),
				MinTLSVersion:     "",
				CipherSuites:      []string(nil),
			},
			MaxRequestsInFlight:   500,
			RequestTimeoutSeconds: 3600,
		},
		AuthConfig: legacyconfigv1.MasterAuthConfig{
			RequestHeader: &legacyconfigv1.RequestHeaderAuthenticationOptions{
				ClientCA:            "front-proxy-ca.crt",
				ClientCommonNames:   []string{"aggregator-front-proxy"},
				UsernameHeaders:     []string{"X-Remote-User"},
				GroupHeaders:        []string{"X-Remote-Group"},
				ExtraHeaderPrefixes: []string{"X-Remote-Extra-"},
			},
			WebhookTokenAuthenticators: []legacyconfigv1.WebhookTokenAuthenticator(nil),
			OAuthMetadataFile:          "",
		},
		AggregatorConfig: legacyconfigv1.AggregatorConfig{
			ProxyClientInfo: legacyconfigv1.CertInfo{
				CertFile: "aggregator-front-proxy.crt",
				KeyFile:  "aggregator-front-proxy.key",
			},
		},
		CORSAllowedOrigins: []string{
			"(?i)//127\\.0\\.0\\.1(:|\\z)",
			"(?i)//localhost(:|\\z)",
			"(?i)//192\\.168\\.0\\.160(:|\\z)",
			"(?i)//kubernetes\\.default(:|\\z)",
			"(?i)//kubernetes(:|\\z)",
			"(?i)//master\\.ci\\-agnosticd\\-p\\-101\\.mg\\.dog8code\\.com(:|\\z)",
			"(?i)//openshift\\.default(:|\\z)",
			"(?i)//openshift\\.default\\.svc(:|\\z)",
			"(?i)//172\\.30\\.0\\.1(:|\\z)",
			"(?i)//master1\\.ci\\-agnosticd\\-p\\-101\\.internal(:|\\z)",
			"(?i)//openshift\\.default\\.svc\\.cluster\\.local(:|\\z)",
			"(?i)//kubernetes\\.default\\.svc(:|\\z)",
			"(?i)//kubernetes\\.default\\.svc\\.cluster\\.local(:|\\z)",
			"(?i)//openshift(:|\\z)",
		},
		APILevels:       []string{"v1"},
		MasterPublicURL: "https://master.ci-agnosticd-p-101.mg.dog8code.com:443",
		Controllers:     "*",
		AdmissionConfig: legacyconfigv1.AdmissionConfig{
			PluginConfig: map[string]*legacyconfigv1.AdmissionPluginConfig{
				"BuildDefaults": &legacyconfigv1.AdmissionPluginConfig{
					Location: "",
					Configuration: runtime.RawExtension{
						Raw:    []uint8{0x7b, 0x22, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x76, 0x31, 0x22, 0x2c, 0x22, 0x65, 0x6e, 0x76, 0x22, 0x3a, 0x5b, 0x5d, 0x2c, 0x22, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x3a, 0x22, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x2c, 0x22, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x22, 0x3a, 0x7b, 0x22, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x73, 0x22, 0x3a, 0x7b, 0x7d, 0x2c, 0x22, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x22, 0x3a, 0x7b, 0x7d, 0x7d, 0x7d},
						Object: runtime.Object(nil),
					},
				},
				"BuildOverrides": &legacyconfigv1.AdmissionPluginConfig{
					Location: "",
					Configuration: runtime.RawExtension{
						Raw:    []uint8{0x7b, 0x22, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x76, 0x31, 0x22, 0x2c, 0x22, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x3a, 0x22, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x7d},
						Object: runtime.Object(nil),
					},
				},
				"MutatingAdmissionWebhook": &legacyconfigv1.AdmissionPluginConfig{
					Location: "",
					Configuration: runtime.RawExtension{
						Raw:    []uint8{0x7b, 0x22, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x61, 0x70, 0x69, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x22, 0x2c, 0x22, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x3a, 0x22, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x41, 0x64, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x2c, 0x22, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x3a, 0x22, 0x2f, 0x64, 0x65, 0x76, 0x2f, 0x6e, 0x75, 0x6c, 0x6c, 0x22, 0x7d},
						Object: runtime.Object(nil),
					},
				},
				"ValidatingAdmissionWebhook": &legacyconfigv1.AdmissionPluginConfig{
					Location: "",
					Configuration: runtime.RawExtension{
						Raw:    []uint8{0x7b, 0x22, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x61, 0x70, 0x69, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x22, 0x2c, 0x22, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x3a, 0x22, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x41, 0x64, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x2c, 0x22, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x3a, 0x22, 0x2f, 0x64, 0x65, 0x76, 0x2f, 0x6e, 0x75, 0x6c, 0x6c, 0x22, 0x7d},
						Object: runtime.Object(nil),
					},
				},
			},
			PluginOrderOverride: []string(nil),
		},
		ControllerConfig: legacyconfigv1.ControllerConfig{
			Controllers: []string(nil),
			Election: &legacyconfigv1.ControllerElectionConfig{
				LockName:      "openshift-master-controllers",
				LockNamespace: "",
				LockResource: legacyconfigv1.GroupResource{
					Group:    "",
					Resource: "",
				},
			},
			ServiceServingCert: legacyconfigv1.ServiceServingCert{
				Signer: &legacyconfigv1.CertInfo{
					CertFile: "service-signer.crt",
					KeyFile:  "service-signer.key",
				},
			},
		},
		EtcdStorageConfig: legacyconfigv1.EtcdStorageConfig{
			KubernetesStorageVersion: "v1",
			KubernetesStoragePrefix:  "kubernetes.io",
			OpenShiftStorageVersion:  "v1",
			OpenShiftStoragePrefix:   "openshift.io",
		},
		EtcdClientInfo: legacyconfigv1.EtcdConnectionInfo{
			URLs: []string{
				"https://master1.ci-agnosticd-p-101.internal:2379",
			},
			CA: "master.etcd-ca.crt",
			CertInfo: legacyconfigv1.CertInfo{
				CertFile: "master.etcd-client.crt",
				KeyFile:  "master.etcd-client.key",
			},
		},
		KubeletClientInfo: legacyconfigv1.KubeletConnectionInfo{
			Port: 0x280a,
			CA:   "ca-bundle.crt",
			CertInfo: legacyconfigv1.CertInfo{
				CertFile: "master.kubelet-client.crt",
				KeyFile:  "master.kubelet-client.key",
			},
		},
		KubernetesMasterConfig: legacyconfigv1.KubernetesMasterConfig{
			APILevels:                  []string(nil),
			DisabledAPIGroupVersions:   map[string][]string(nil),
			MasterIP:                   "192.168.0.160",
			MasterEndpointReconcileTTL: 0,
			ServicesSubnet:             "172.30.0.0/16",
			ServicesNodePortRange:      "",
			SchedulerConfigFile:        "/etc/origin/master/scheduler.json",
			PodEvictionTimeout:         "",
			ProxyClientInfo: legacyconfigv1.CertInfo{
				CertFile: "master.proxy-client.crt",
				KeyFile:  "master.proxy-client.key",
			},
			APIServerArguments: legacyconfigv1.ExtendedArguments{
				"storage-backend": []string{
					"etcd3",
				},
				"storage-media-type": []string{
					"application/vnd.kubernetes.protobuf",
				},
			},
			ControllerArguments: legacyconfigv1.ExtendedArguments{
				"cluster-signing-cert-file": []string{
					"/etc/origin/master/ca.crt",
				},
				"cluster-signing-key-file": []string{
					"/etc/origin/master/ca.key",
				},
				"pv-recycler-pod-template-filepath-hostpath": []string{
					"/etc/origin/master/recycler_pod.yaml",
				},
				"pv-recycler-pod-template-filepath-nfs": []string{
					"/etc/origin/master/recycler_pod.yaml",
				},
			},
			SchedulerArguments: legacyconfigv1.ExtendedArguments(nil),
		},
		EtcdConfig: nil,
		OAuthConfig: &legacyconfigv1.OAuthConfig{
			MasterCA:                    &masterCA,
			MasterURL:                   "https://master1.ci-agnosticd-p-101.internal:443",
			MasterPublicURL:             "https://master.ci-agnosticd-p-101.mg.dog8code.com:443",
			AssetPublicURL:              "https://master.ci-agnosticd-p-101.mg.dog8code.com/console/",
			AlwaysShowProviderSelection: false,
			IdentityProviders: []legacyconfigv1.IdentityProvider{
				legacyconfigv1.IdentityProvider{
					Name:            "htpasswd_auth",
					UseAsChallenger: true,
					UseAsLogin:      true,
					MappingMethod:   "claim",
					Provider: runtime.RawExtension{
						Raw:    []uint8{0x7b, 0x22, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x76, 0x31, 0x22, 0x2c, 0x22, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x3a, 0x22, 0x2f, 0x65, 0x74, 0x63, 0x2f, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x68, 0x74, 0x70, 0x61, 0x73, 0x73, 0x77, 0x64, 0x22, 0x2c, 0x22, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x3a, 0x22, 0x48, 0x54, 0x50, 0x61, 0x73, 0x73, 0x77, 0x64, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x22, 0x7d},
						Object: runtime.Object(nil),
					},
				},
			},
			GrantConfig: legacyconfigv1.GrantConfig{
				Method:               "auto",
				ServiceAccountMethod: "",
			},
			SessionConfig: &legacyconfigv1.SessionConfig{
				SessionSecretsFile:   "/etc/origin/master/session-secrets.yaml",
				SessionMaxAgeSeconds: 3600,
				SessionName:          "ssn",
			},
			TokenConfig: legacyconfigv1.TokenConfig{
				AuthorizeTokenMaxAgeSeconds:         500,
				AccessTokenMaxAgeSeconds:            86400,
				AccessTokenInactivityTimeoutSeconds: (*int32)(nil),
			},
			Templates: nil,
		},
		DNSConfig: &legacyconfigv1.DNSConfig{
			BindAddress:           "0.0.0.0:8053",
			BindNetwork:           "tcp4",
			AllowRecursiveQueries: false,
		},
		ServiceAccountConfig: legacyconfigv1.ServiceAccountConfig{
			ManagedNames: []string{
				"default",
				"builder",
				"deployer",
			},
			LimitSecretReferences: false,
			PrivateKeyFile:        "serviceaccounts.private.key",
			PublicKeyFiles: []string{
				"serviceaccounts.public.key",
			},
			MasterCA: "ca-bundle.crt",
		},
		MasterClients: legacyconfigv1.MasterClients{
			OpenShiftLoopbackKubeConfig: "openshift-master.kubeconfig",
			OpenShiftLoopbackClientConnectionOverrides: &legacyconfigv1.ClientConnectionOverrides{
				AcceptContentTypes: "application/vnd.kubernetes.protobuf,application/json",
				ContentType:        "application/vnd.kubernetes.protobuf",
				QPS:                300,
				Burst:              600,
			},
		},
		ImageConfig: legacyconfigv1.ImageConfig{
			Format: "registry.redhat.io/openshift3/ose-${component}:${version}",
			Latest: false,
		},
		ImagePolicyConfig: legacyconfigv1.ImagePolicyConfig{
			MaxImagesBulkImportedPerRepository:         0,
			DisableScheduledImport:                     false,
			ScheduledImageImportMinimumIntervalSeconds: 0,
			MaxScheduledImageImportsPerMinute:          0,
			AllowedRegistriesForImport: &legacyconfigv1.AllowedRegistries{
				legacyconfigv1.RegistryLocation{
					DomainName: "registry1.test.com",
					Insecure:   true,
				},
				legacyconfigv1.RegistryLocation{
					DomainName: "registry2.test.com",
					Insecure:   false,
				},
			},
			InternalRegistryHostname: "docker-registry.default.svc:5000",
			ExternalRegistryHostname: "external-registry.example.com",
			AdditionalTrustedCA:      "",
		},
		PolicyConfig: legacyconfigv1.PolicyConfig{
			UserAgentMatchingConfig: legacyconfigv1.UserAgentMatchingConfig{
				RequiredClients:         []legacyconfigv1.UserAgentMatchRule(nil),
				DeniedClients:           []legacyconfigv1.UserAgentDenyRule(nil),
				DefaultRejectionMessage: "",
			},
		},
		ProjectConfig: legacyconfigv1.ProjectConfig{
			DefaultNodeSelector:    "node-role.kubernetes.io/compute=true",
			ProjectRequestMessage:  "To provision Projects you must request access in https://labs.opentlc.com or https://rhpds.redhat.com",
			ProjectRequestTemplate: "default/project-request",
			SecurityAllocator: &legacyconfigv1.SecurityAllocator{
				UIDAllocatorRange:   "1000000000-1999999999/10000",
				MCSAllocatorRange:   "s0:/2",
				MCSLabelsPerProject: 5,
			},
		},
		RoutingConfig: legacyconfigv1.RoutingConfig{
			Subdomain: "apps.ci-agnosticd-p-101.mg.dog8code.com",
		},
		NetworkConfig: legacyconfigv1.MasterNetworkConfig{
			NetworkPluginName:            "redhat/openshift-ovs-subnet",
			DeprecatedClusterNetworkCIDR: "",
			ClusterNetworks: []legacyconfigv1.ClusterNetworkEntry{
				legacyconfigv1.ClusterNetworkEntry{
					CIDR:             "10.1.0.0/16",
					HostSubnetLength: 0x9,
				},
			},
			DeprecatedHostSubnetLength: 0x0,
			ServiceNetworkCIDR:         "172.30.0.0/16",
			ExternalIPNetworkCIDRs: []string{
				"0.0.0.0/0",
			},
			IngressIPNetworkCIDR: "",
			VXLANPort:            0x0,
		},
		VolumeConfig: legacyconfigv1.MasterVolumeConfig{
			DynamicProvisioningEnabled: &dynamicProvisioningEnabled,
		},
		JenkinsPipelineConfig: legacyconfigv1.JenkinsPipelineConfig{
			AutoProvisionEnabled: (*bool)(nil),
			TemplateNamespace:    "",
			TemplateName:         "",
			ServiceName:          "",
			Parameters:           map[string]string(nil),
		},
		AuditConfig: legacyconfigv1.AuditConfig{
			Enabled:                  false,
			AuditFilePath:            "",
			MaximumFileRetentionDays: 0,
			MaximumRetainedFiles:     0,
			MaximumFileSizeMegabytes: 0,
			PolicyFile:               "",
			PolicyConfiguration: runtime.RawExtension{
				Raw:    []uint8(nil),
				Object: runtime.Object(nil),
			},
			LogFormat:         "",
			WebHookKubeConfig: "",
			WebHookMode:       "",
		},
		DisableOpenAPI: false,
	}

	t.Run("Decode MasterConfig", func(t *testing.T) {
		filename := "testdata/master_config.yaml"
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatalf("Cannot read file %s", filename)
		}
		masterConfig, err := MasterConfig(content)
		require.NoError(t, err)
		assert.Equal(t, masterConfig, expectedMasterConfig)
	})
}
