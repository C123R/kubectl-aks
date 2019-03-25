package util

import (
	"context"
	"fmt"
	container "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

// AksCluster is an object representing details for AKS cluster
type AksCluster struct {
	ResourceGroup string
	K8sVersion    string
}

// GetAKS returns list of AKS clusters in resource group
func GetAKS(sess *AzureSession, name string) (string, error) {

	var err error
	var kubeconfig string

	var akslist AksCluster
	aksList, err := akslist.ListAKS(sess)

	if _, ok := aksList[name]; !ok {
		return kubeconfig, fmt.Errorf("invalid cluster name (%v), use `kubectl aks list` to get the correct list", name)
	}

	crClient := container.NewManagedClustersClient(sess.SubscriptionID)
	crClient.Authorizer = sess.Authorizer
	result, err := crClient.ListClusterUserCredentials(context.Background(), aksList[name].ResourceGroup, name)
	if err != nil {
		return kubeconfig, err
	}

	kubeconfig = string(*(*result.Kubeconfigs)[0].Value)

	return kubeconfig, err

}

// ListAKS returns list of AKS clusters in resource group
func (a *AksCluster) ListAKS(sess *AzureSession) (map[string]AksCluster, error) {

	mapOfAKSCluster := make(map[string]AksCluster)
	var err error
	crClient := container.NewManagedClustersClient(sess.SubscriptionID)
	crClient.Authorizer = sess.Authorizer

	for list, err := crClient.ListComplete(context.Background()); list.NotDone(); err = list.Next() {
		if err != nil {
			return mapOfAKSCluster, fmt.Errorf("error get the list of aks clusters: %v", err)
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

// ManageConfig is use to merge kubeconfiguration with existing config
func ManageConfig(config string, path string) error {

	var err error
	if _, err := os.Stat(path); os.IsNotExist(err) {

		return fmt.Errorf("Default/Provided path does not exist,\"%v\"", err)
	}
	file, _ := ioutil.TempFile("/tmp", "temp")

	// Delete temp file
	defer os.Remove(file.Name())

	tempFile := file.Name()
	// Write Kubernetes configuration for requested cluster in temporary file
	err = ioutil.WriteFile(tempFile, []byte(config), 0600)
	file.Sync()
	// handle this error
	if err != nil {
		// print it out
		fmt.Println(err)
	}

	rules := clientcmd.ClientConfigLoadingRules{
		Precedence: []string{clientcmd.RecommendedHomeFile, tempFile},
	}

	return err
}
