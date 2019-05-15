package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"github.com/briandowns/spinner"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// AksCluster is an object representing details for AKS cluster
type AksCluster struct {
	ResourceGroup string
	K8sVersion    string
	Nodes         int32
}

// ValidAKSCluster is an object representing validated AKS Cluster
type ValidAKSCluster struct {
	ResourceGroup string
	K8sVersion    string
	Nodes         int32
	AKSInstance   containerservice.ManagedCluster
}

// GetCredentials returns list of AKS clusters in resource group
func (aksClient AKSClient) GetCredentials(name string) ([]byte, error) {

	aks, err := aksClient.validateAKSCluster(name)
	if err != nil {
		return nil, err
	}

	result, err := aksClient.ContainerService.ListClusterUserCredentials(context.Background(), aks.ResourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("unable to get the AKS clusters credentials for (%v), Error: %v", name, err)
	}

	return *(*result.Kubeconfigs)[0].Value, err

}

// UpgradeAKS returns list of AKS clusters in resource group
func (aksClient AKSClient) UpgradeAKS(ctx context.Context, name string, k8sVersion string) error {

	cluster, err := aksClient.validateAKSCluster(name)
	if err != nil {
		return err
	}

	k8sUpgradeVersions, _, err := aksClient.GetUpgrades(name)
	if err != nil {
		return err
	}

	if k8sVersion == cluster.K8sVersion {
		fmt.Printf("INFO: %v is currently running on version %v\n\n", name, k8sVersion)
		return nil
	} else if !stringSliceContains(k8sUpgradeVersions, k8sVersion) {
		return fmt.Errorf("Upgrade with version %v is either not allowed or not available for the %v", k8sVersion, name)
	} else if Confirmation() {

		*cluster.AKSInstance.KubernetesVersion = k8sVersion
		initalMsg := fmt.Sprintf("Upgrading %v to Kubernetes version %v ", name, k8sVersion)
		finalMSG := fmt.Sprintf("Successfully upgraded %v to Kubernetes version %v.\n", name, k8sVersion)

		s := getSpinner(initalMsg, finalMSG)

		s.Start()

		err = aksClient.CreateOrUpdate(context.Background(), name, cluster)
		if err != nil {

			s.FinalMSG = err.Error()
			s.Stop()
			fmt.Println()
			return err
		}
		s.Stop()

	} else {
		fmt.Println("Opeartion Cancelled")
	}
	return nil

}

// ScaleAKS returns list of AKS clusters in resource group
func (aksClient AKSClient) ScaleAKS(ctx context.Context, name string, count int32) error {

	cluster, err := aksClient.validateAKSCluster(name)
	if err != nil {
		return err
	}
	if count == cluster.Nodes {
		fmt.Printf("INFO: %v is currently running with %v nodes\n\n", name, count)
		return nil
	} else if Confirmation() {
		for _, agentPoolProfile := range *cluster.AKSInstance.AgentPoolProfiles {
			*agentPoolProfile.Count = count
		}

		initalMsg := fmt.Sprintf("Scaling %v to %d nodes ", name, count)
		finalMSG := fmt.Sprintf("Successfully scaled %v to %d nodes.\n", name, count)

		s := getSpinner(initalMsg, finalMSG)
		s.Start()

		// Making a request to update the AKS Cluster
		err = aksClient.CreateOrUpdate(context.Background(), name, cluster)
		if err != nil {
			s.FinalMSG = err.Error()
			s.Stop()
			fmt.Println()
			return err
		}
		s.Stop()
	} else {
		fmt.Println("Opeartion Cancelled")
	}
	return nil
}

// CreateOrUpdate creates or updates a managed cluster with the specified configuration for agents and Kubernetes
// version.
// Parameters:
// resourceGroupName - the name of the resource group.
// name - the name of AKS Cluster.
// object of ValidAKSCluster.
func (aksClient AKSClient) CreateOrUpdate(ctx context.Context, name string, cluster ValidAKSCluster) error {

	futureResponse, err := aksClient.ContainerService.CreateOrUpdate(ctx, cluster.ResourceGroup, name, cluster.AKSInstance)
	if err != nil {
		return err
	}

	// Setting one hour as DefaultPollingDuration
	aksClient.ContainerService.Client.PollingDuration = DefaultPollingDuration

	err = futureResponse.WaitForCompletionRef(ctx, aksClient.ContainerService.Client)
	if err != nil {
		return fmt.Errorf("Cannot get the AKS Cluster create or update future response: %v", err)
	}
	return nil
}

// get the Spinner object for the long running process
func getSpinner(initialMsg string, finalMsg string) *spinner.Spinner {

	spinner := spinner.New(spinner.CharSets[4], 100*time.Millisecond)
	spinner.Prefix = initialMsg
	spinner.FinalMSG = finalMsg
	return spinner

}

// GetUpgrades get the available upgrade version and current version
func (aksClient AKSClient) GetUpgrades(name string) ([]string, string, error) {

	var k8sUpgradeVersions []string
	var currentVersion string

	cluster, err := aksClient.validateAKSCluster(name)
	if err != nil {
		return k8sUpgradeVersions, currentVersion, err
	}

	instance, err := aksClient.ContainerService.GetUpgradeProfile(context.Background(), cluster.ResourceGroup, name)
	if err != nil {
		return k8sUpgradeVersions, currentVersion, fmt.Errorf("unable to get the available upgrades for (%v), Error: %v", name, err)
	} else if instance.ControlPlaneProfile.Upgrades == nil {
		return k8sUpgradeVersions, *instance.ControlPlaneProfile.KubernetesVersion, nil
	}

	return *instance.ControlPlaneProfile.Upgrades, *instance.ControlPlaneProfile.KubernetesVersion, nil
}

// ListAKS returns list of AKS clusters in resource group
func (aksClient AKSClient) ListAKS() (map[string]AksCluster, error) {

	mapOfAKSCluster := make(map[string]AksCluster)

	for list, err := aksClient.ContainerService.ListComplete(context.Background()); list.NotDone(); err = list.Next() {
		if err != nil {
			return mapOfAKSCluster, fmt.Errorf("error getting the list of aks clusters: %v", err)
		}
		clusterName := *list.Value().Name
		rg := strings.Split(*list.Value().NodeResourceGroup, "_")[1]
		version := *list.Value().KubernetesVersion
		nodes := (*(*list.Value().AgentPoolProfiles)[0].Count)

		mapOfAKSCluster[clusterName] = makeMapOfCluster(rg, version, nodes)
	}
	return mapOfAKSCluster, nil

}

// Validate AKS Cluster with a specified name.
// Parameters:
// name - the name of the AKS Cluster.
// Retunrs Object of ValidAKSCluster(empty) with error, if no such AKS Cluster in Azure or AKS Cluster is in Failed State
func (aksClient AKSClient) validateAKSCluster(name string) (ValidAKSCluster, error) {

	aksList, err := aksClient.ListAKS()

	var validCluster ValidAKSCluster

	if _, ok := aksList[name]; !ok {
		return validCluster, fmt.Errorf("Invalid Cluster name (%v), use `kubectl aks list` to get the correct list", name)
	}
	instance, err := aksClient.ContainerService.Get(context.Background(), aksList[name].ResourceGroup, name)

	if err != nil {
		return validCluster, fmt.Errorf("Unable to get the AKS instance for %v in %v", name, aksList[name].ResourceGroup)
	} else if *instance.ProvisioningState == "Failed" {
		return validCluster, fmt.Errorf("AKS Cluster %v is currently in a failed state", name)
	}

	validCluster = ValidAKSCluster{
		Nodes:         aksList[name].Nodes,
		ResourceGroup: aksList[name].ResourceGroup,
		K8sVersion:    aksList[name].K8sVersion,
		AKSInstance:   instance,
	}

	return validCluster, err

}

func makeMapOfCluster(rg string, version string, nodes int32) AksCluster {

	return AksCluster{
		ResourceGroup: rg,
		K8sVersion:    version,
		Nodes:         nodes,
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

func stringSliceContains(s []string, v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}

// Confirmation asks for user input before proceeding
func Confirmation() bool {

	var input string
	fmt.Printf("Do you want to continue with this operation? [y|n]: ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		panic(err)
	}

	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	if input == "y" || input == "yes" {
		return true
	}
	return false

}
