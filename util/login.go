package util

import (
	"context"
	"encoding/json"
	"fmt"
	container "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
)

// AzureSession is an object representing session for subscription
type AzureSession struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
}

// AksCluster is an object representing details for AKS cluster
type AksCluster struct {
	ResourceGroup string
	K8sVersion    string
}

func makeMapOfCluster(rg string, version string) AksCluster {

	return AksCluster{
		ResourceGroup: rg,
		K8sVersion:    version,
	}
}

// ReadJSON returns json content of sdk auth file
func ReadJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Wrap(err, "Can't open the file")
	}

	contents := make(map[string]interface{})
	err = json.Unmarshal(data, &contents)

	if err != nil {
		err = errors.Wrap(err, "Can't unmarshal file")
	}

	return &contents, err
}

// NewSessionFromFile returns Azure Session Object
func NewSessionFromFile() (*AzureSession, error) {

	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)

	if err != nil {
		return nil, fmt.Errorf("can't initialize authorizer: %v", err)
	}

	authInfo, err := ReadJSON(os.Getenv("AZURE_AUTH_LOCATION"))
	if err != nil {
		return nil, fmt.Errorf("can't get auth file: %v", err)
	}

	sess := AzureSession{
		SubscriptionID: (*authInfo)["subscriptionId"].(string),
		Authorizer:     authorizer,
	}

	return &sess, nil
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
