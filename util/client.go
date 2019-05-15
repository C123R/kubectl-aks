package util

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/mitchellh/go-homedir"
)

// AKSClient is an object representing session for subscription
type AKSClient struct {
	SubscriptionID   string
	Authorizer       autorest.Authorizer
	ContainerService containerservice.ManagedClustersClient
}

func getAzureAuth() (auth.FileSettings, error) {

	var s auth.FileSettings
	defaultConfig, _ := homedir.Expand("~/.kube/azure.auth")
	if _, err := os.Stat(os.Getenv("AZURE_AUTH_LOCATION")); os.IsNotExist(err) {
		if _, err = os.Stat(defaultConfig); os.IsNotExist(err) {
			return s, fmt.Errorf("Cannot get the Azure Auth: %v", err)
		}
		os.Setenv("AZURE_AUTH_LOCATION", defaultConfig)
	}
	return auth.GetSettingsFromFile()
}

// NewAKSClient returns Azure Session Object
func NewAKSClient() (AKSClient, error) {

	var aksClient AKSClient

	settings, err := getAzureAuth()
	if err != nil {
		return aksClient, err
	}
	authorizer, err := auth.NewAuthorizerFromCLI()
	crService := containerservice.NewManagedClustersClient(settings.GetSubscriptionID())
	crService.Authorizer = authorizer

	aksClient = AKSClient{
		SubscriptionID:   settings.GetSubscriptionID(),
		Authorizer:       authorizer,
		ContainerService: crService,
	}

	if err != nil {
		authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
		aksClient = AKSClient{
			SubscriptionID:   settings.GetSubscriptionID(),
			Authorizer:       authorizer,
			ContainerService: crService,
		}
		if err != nil {
			return aksClient, fmt.Errorf("can't initialize authorizer: %v", err)
		}
	}

	return aksClient, nil
}
