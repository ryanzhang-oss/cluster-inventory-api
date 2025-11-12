package main

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"

	"sigs.k8s.io/cluster-inventory-api/apis/v1alpha1"
	"sigs.k8s.io/cluster-inventory-api/pkg/credentials"
)

// The example below showcases how to use Azure's kubelogin exec plugin to sign into an
// AKS cluster using the workload identity method.
//
// As the method requires cluster-specific information such as tenant ID and client ID,
// the example also demonstrates how to pass in additional command-line arguments
// and/or environment variables to the exec plugin via the ClusterProfile API using the
// reserved extensions, as defined in KEP 5339.
func main() {
	providers := []credentials.Provider{
		{
			Name: "aks-workload-identity",
			ExecConfig: &clientcmdapi.ExecConfig{
				Command: "kubelogin",
				Args: []string{
					"get-token",
					"--login",
					"workloadidentity",
					"--tenant-id",
					"72f988bf-86f1-41af-91ab-2d7cd011db47",
					"--authority-host",
					"https://login.microsoftonline.com/",
					"--federated-token-file",
					// The well-known path where AKS mounts the service account token as a projected volume.
					//
					// This is an application-specific information and it can be configured to
					// a different path if needed.
					//"/var/run/secrets/azure/tokens/azure-identity-token",
					"/home/zhangryan/github/cluster-inventory-api/azure-idenity-token",
					"--server-id",
					"6dae42f8-4368-4678-94ff-3960e28e3630",
				},
				Env: []clientcmdapi.ExecEnvVar{
					{
						Name:  "AZURE_CLIENT_ID",
						Value: "120a6d5d-3277-4a86-8b76-065929547177",
					},
				},
				APIVersion:         "client.authentication.k8s.io/v1beta1",
				ProvideClusterInfo: false,
				InteractiveMode:    clientcmdapi.NeverExecInteractiveMode,
			},
			AdditionalCLIArgEnvVarExtensionFlag: credentials.AdditionalCLIArgEnvVarExtensionFlagAllow,
		},
	}
	cps := credentials.New(providers)
	profile := &v1alpha1.ClusterProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bravelion",
			Namespace: "fleet-system",
		},
		Spec: v1alpha1.ClusterProfileSpec{
			DisplayName: "bravelion",
			ClusterManager: v1alpha1.ClusterManager{
				Name: "kubefleet",
			},
		},
		Status: v1alpha1.ClusterProfileStatus{
			Version: v1alpha1.ClusterVersion{
				Kubernetes: "v1.32.7",
			},
			CredentialProviders: []v1alpha1.CredentialProvider{
				{
					Name: "aks-workload-identity",
					Cluster: clientcmdapiv1.Cluster{
						Server:                   "https://query-rout-kubecon-2025-dem-d712bf-wg0xnuz5.hcp.westus2.azmk8s.io:443",
						CertificateAuthorityData: []byte("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUU2VENDQXRHZ0F3SUJBZ0lSQUxUa04xN0tIUmlvRWZSWHlMdUFEV0V3RFFZSktvWklodmNOQVFFTEJRQXcKRFRFTE1Ba0dBMVVFQXhNQ1kyRXdJQmNOTWpVeE1UQTFNakV5TmpNNFdoZ1BNakExTlRFeE1EVXlNVE0yTXpoYQpNQTB4Q3pBSkJnTlZCQU1UQW1OaE1JSUNJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBZzhBTUlJQ0NnS0NBZ0VBCnFGTENQR2x0ZjhIanVVbi9kZEVpUGxaRENadjdKQUdER0ZWajlqS0JEd0ozZHJPMTVaekpXcEE1R3ZzaTBla3MKNThsa1gxR1czMDFyWXZRclVsOUhOYUc0M2hJRzBuaGw5SldxNUg0bWF1QndHM0t4SUpYcDV1Qk44bUFGQkFCQwp5cjZoOG9aUEtkYUh0dzliTytueWxTRjc3bzVTZ0NVWE5GbDhnTy8yUEhoOVl3VzhOWDhYRGdSUUdSTlV5RjJwCmdCUCtvdnhZbTN0aFpuTU5rbExjYU0zR0tvWHdlSXoxa2c1d3FIOFZ4OWF0L01QY3ozWDY3djJJdDFKRCtKMjMKY1N2ZnptUUtaYjF0Q2xXWGF5NHNEekZRZjM3QmNveDFZdSsxQlk4WDlVeEZiY1BjUjJyT3VwUUtqRlhkRm5mSQpMaEN3NzFaQ1BnK21pTEtldW9kdzBQTHJzYnRraktMNlhkQ0ZDTlRIRGN1OTFqczJyUW8za2dSRGRuOTV4S0xBCkFWeEhUeFR2a0xmTjJPM0YybjhIYmxIa3h3Nk9neHlVcFdvb3A1cHQ3eGlVbTBSV3l2L3Vkb3JRYnVyS29uQk4KUkFYaWVpZjQ1aktNR1hiR1NRRjQ5SE9laUpuMHFaTjZOSEJIODFTNmFPUjYvdmY3cDJ6WXdUSG1Ta0tIQURzbQo2eGwvSllVZlhwa0tJL1V0MVNDZEYzTDk2VGp0dlIzWWFiZWx4Y0NCY1h5NFV0K29jMXpsWE1taTNlb0kyWFNICnpVQkpkcDh5NEhXbGlFSDRML2JzbjE4QWYvMEZKcFNlcDlSN2hrcjJxanpsMHhEREZGeFI3UmIwbS82Wng1aVEKYXFMZTBmV1kwbURHRTVrOGxmbGwya3pLVWNZNktIMnJURjVTOVc0TGJJRUNBd0VBQWFOQ01FQXdEZ1lEVlIwUApBUUgvQkFRREFnS2tNQThHQTFVZEV3RUIvd1FGTUFNQkFmOHdIUVlEVlIwT0JCWUVGTUxIcXg4NzRIL2RRVmRoCmFMOXBId0hWOTQ2Y01BMEdDU3FHU0liM0RRRUJDd1VBQTRJQ0FRQm43Zk1iQ0FqYTRlU21aaC94OUk3OVB4U2sKbE9WTEI1U1htcnExaEZlU2VFUmdDcUt1enhYWkZXRC9jL3pkYjZwcVVJVko3QUFZV0ZyN0R0dEFXUGpqWWl6KwpRRVFlOWpHTkx5NXlaN1pNcndOZDh4QkZ3Slc4REhHQXhOYkV4Q2dKbzB0emRCSXR3N0lLOXU5ampJK1MrUS9CCnhvbGd2SkZGc042SlpTd1VoRHhxQ2RCVzBSRkpjNjhsVStMdDRSdlk0aXQxaDRNc3ljQW16ZzBMS1F6aXZHTDEKa1ZaSDlXa2JmVDhYaE00Qkg5NmloTXVOWTI4MzZFK1JUMlVLWEQ2N1I1cGF2bWFlNDlGQThmemlUVVFiZ3NwSgpqdHMzZWNEbG1Pc0gyNTdLdmNCQkVaSTVLVE9hdjNqYW4wcnZOdnN6WStUQyt5Q0hqRmo5ODJvYnZQRWdPMlcrCko0ZDVKVitXTit4VXg4OEY1ejhBMFVSOUh4Wkt0cEYwVFJ1U0NGQ0JmMnkxNTcyUm91ZHM1R1NsODJnSW5maEkKYkJHdkJFajRtSDlyVkdBbjY2YmR4elVNNWkzclJ2NUpnUVFESnEzRVFzSE1va0p6UXpSdUZka3IxMzlPMUlYbgo5UE9QOGFTQjYycDY5WGRWWGJRRjZYbHE2NUFMS3pEUmkrOFNKVmM1Ti90d1JtU1JhNUlKd2VrZVp2ek16OS96CnRzdTFuUjUrQzhZcnF6cURhQUc3RXhTNTYxU0Z4ZURVWWdTYlJ6a0RlN3VHcXhOSlY4VDROOUxmUVM4b2Y2aFgKWDZBUTQxTFg0Uk1NU1dEbU0rbmx5U0hFc3NWL0YzblpzdHEvamdob1lZMlBzQktzK21Yd0h6VTVXZFFzM1d4eQpLbW9MelFKbEwwS0IrNGVadEE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="),
					},
				},
			},
		},
	}

	restConfig, err := cps.BuildConfigFromCP(profile)
	if err != nil {
		log.Fatalf("Failed to prepare REST config: %v", err)
	}

	// The generated REST config can be used to build a Kubernetes client.
	//
	// It will invoke the kubelogin plugin as follows:
	//
	// kubelogin get-token \
	//     --login workloadidentity \
	//     --federated-token-file /var/run/secrets/tokens/azure-identity-token \
	//     --tenant-id TENANT_ID \
	//     --client-id CLIENT_ID \
	//     --authority-host https://login.microsoftonline.com/
	log.Printf("Prepared REST config:\n%+v", restConfig)
	log.Printf("CLI Args: %s", restConfig.ExecProvider.Args)
	log.Printf("Env Vars: %+v", restConfig.ExecProvider.Env)

	// Build a Kubernetes clientset from the REST config.
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	log.Printf("Successfully created Kubernetes client: %+v", kubeClient)
	namespaces, err := kubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list namespaces: %v", err)
	}
	log.Printf("Got namespaces : %+v", namespaces.Items)
}
