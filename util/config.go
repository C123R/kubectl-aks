package util

import (
	"context"
	"fmt"
	container "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"os"
	"strings"
)

// AksCluster is an object representing details for AKS cluster
type AksCluster struct {
	ResourceGroup string
	K8sVersion    string
}

// GetAKS returns list of AKS clusters in resource group
func GetAKS(sess *AzureSession, name string) ([]byte, error) {

	var err error

	var akslist AksCluster
	aksList, err := akslist.ListAKS(sess)

	if _, ok := aksList[name]; !ok {
		return nil, fmt.Errorf("invalid cluster name (%v), use `kubectl aks list` to get the correct list", name)
	}

	crClient := container.NewManagedClustersClient(sess.SubscriptionID)
	crClient.Authorizer = sess.Authorizer
	result, err := crClient.ListClusterUserCredentials(context.Background(), aksList[name].ResourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("unable to get the AKS clusters credentials for (%v), Error: %v", name, err)
	}

	return *(*result.Kubeconfigs)[0].Value, err

}

// ListAKS returns list of AKS clusters in resource group
func (a *AksCluster) ListAKS(sess *AzureSession) (map[string]AksCluster, error) {

	mapOfAKSCluster := make(map[string]AksCluster)
	var err error
	crClient := container.NewManagedClustersClient(sess.SubscriptionID)
	crClient.Authorizer = sess.Authorizer

	for list, err := crClient.ListComplete(context.Background()); list.NotDone(); err = list.Next() {
		if err != nil {
			return mapOfAKSCluster, fmt.Errorf("error getting the list of aks clusters: %v", err)
		}

		clusterName := *list.Value().Name
		rg := strings.Split(*list.Value().NodeResourceGroup, "_")[1]
		version := *list.Value().KubernetesVersion

		mapOfAKSCluster[clusterName] = makeMapOfCluster(rg, version)

	}
	return mapOfAKSCluster, err

}

func makeMapOfCluster(rg string, version string) AksCluster {

	return AksCluster{
		ResourceGroup: rg,
		K8sVersion:    version,
	}
}

// MergeConfig is use to merge kubernetes configuration with default kube config ~/.kube/config
func MergeConfig(config []byte, path string) error {

	var err error
	file, err := ioutil.TempFile("/tmp", "temp")
	if err != nil {
		return err
	}
	// Delete temp file
	defer os.Remove(file.Name())

	tempFile := file.Name()
	// Write Kubernetes configuration for requested cluster in temporary file
	err = ioutil.WriteFile(tempFile, config, 0600)
	file.Sync()
	if err != nil {
		return err
	}

	tempConfig, err := encodeConfig([]string{tempFile})
	if err != nil {
		return err
	}
	mergedConfig, err := encodeConfig([]string{DefaultKubeConfig, tempFile})
	if err != nil {
		return err
	}

	// Merge new cluster as current context in mergedConfig
	mergedConfig.CurrentContext = tempConfig.CurrentContext
	err = clientcmd.WriteToFile(*mergedConfig, path)

	if err != nil {
		return err
	}
	return err
}

func encodeConfig(precedence []string) (*clientcmdapi.Config, error) {

	var err error
	rules := clientcmd.ClientConfigLoadingRules{
		Precedence: precedence,
	}

	config, err := rules.Load()
	if err != nil {
		return nil, err
	}
	return config, err
}
