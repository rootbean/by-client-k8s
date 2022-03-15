package clientk8s

import (
	"log"
	"os"
	"sync"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var doOnce sync.Once
var k8sClientSet *kubernetes.Clientset

// Generic
type Metav1TypeMeta struct {
	Kind       string
	APIVersion string
}

type Metav1ObjectMeta struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
}

type Metav1CreateOptions struct {
	// Valid values are:
	// - Ignore: ignores unknown/duplicate fields.
	// - Warn: responds with a warning for each
	// unknown/duplicate field, but successfully serves the request.
	// - Strict: fails the request on unknown/duplicate fields.
	// +optional
	FieldValidation string
}

// RbacV1
type Rbacv1PolicyRule struct {
	Verbs           []string
	APIGroups       []string
	Resources       []string
	ResourceNames   []string
	NonResourceURLs []string
}

type Rbacv1Subject struct {
	Kind      string
	APIGroup  string
	Name      string
	Namespace string
}

type Rbacv1RoleRef struct {
	APIGroup string
	Kind     string
	Name     string
}

// PVC

type PersistentVolumeAccessMode struct {
	ReadWriteOnce bool
	ReadOnlyMany  bool
	ReadWriteMany bool
}

// Regcred
type SecretTypeStruct struct {
	// SecretTypeBasicAuth,
	// SecretTypeBootstrapToken,
	// SecretTypeDockerConfigJson,
	// SecretTypeDockercfg,
	// SecretTypeOpaque,
	// SecretTypeSSHAuth,
	// SecretTypeServiceAccountToken,
	// SecretTypeTLS
	SecretType string
}

// apps

type ContainerPortStruct struct {
	Name          string
	HostPort      int32
	ContainerPort int32
	// TCP
	// UDP
	Protocol string
	HostIP   string
}

type EnvFromSourceStruct struct {
	Prefix       string
	ConfigMapRef string
	SecretRef    string
}

type EnvVarStruct struct {
	Name  string
	Value string
}

type VolumeMountStruct struct {
	Name      string
	ReadOnly  bool
	MountPath string
	SubPath   string
}

type VolumeDevicesStruct struct {
	Name       string
	DevicePath string
}

type ResourceListStruct struct {
	resourcesLimitsCPU      string
	resourcesLimitsMemory   string
	resourcesRequestsCPU    string
	resourcesRequestsMemory string
}

type DeploymentContainerStruct struct {
	ContainerName          string
	ContainerImage         string
	ContainerCommand       []string
	ContainerArgs          []string
	ContainerWorkingDir    string
	ContainerPorts         []ContainerPortStruct
	ContainerEnvFrom       []EnvFromSourceStruct
	ContainerEnvVar        []EnvVarStruct
	ContainerVolumeMounts  []VolumeMountStruct
	ContainerVolumeDevices []VolumeDevicesStruct
	ContainerResource      ResourceListStruct
}

func init() {
	doOnce.Do(func() {
		if k8sClientSet == nil {

			k8sClientSet, err = createClient()

			if err != nil {
				log.Fatalf("Error get client: %v", err)
			}

		}
	})
}

var err error

// createClient ...
func createClient() (*kubernetes.Clientset, error) {

	isCluster := os.Getenv("CLIENT_K8S_RUN_IN_CLUSTER")

	var config *rest.Config

	if isCluster == "cluster" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("CLIENT_K8S_KUBECONFIG"))
	}

	if err != nil {
		log.Fatalf("Error connection in ClusterConfig: %s", err)
		return nil, err

	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error setting new config: %v", err)
		return nil, err

	}

	return clientset, nil

}
