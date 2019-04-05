package util

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/mitchellh/go-homedir"
	"os"
)

// AzureSession is an object representing session for subscription
type AzureSession struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
}

func getAzureAuth() (auth.FileSettings, error) {

	var s auth.FileSettings
	defaultConfig, _ := homedir.Expand("~/.kube/azure.auth")
	if _, err := os.Stat(os.Getenv("AZURE_AUTH_LOCATION")); os.IsNotExist(err) {
		if _, err = os.Stat(defaultConfig); os.IsNotExist(err) {
			return s, fmt.Errorf("cannot get auth file: %v", err)
		}
		os.Setenv("AZURE_AUTH_LOCATION", defaultConfig)
	}
	return auth.GetSettingsFromFile()
}

// NewSessionFromFile returns Azure Session Object
func NewSessionFromFile() (*AzureSession, error) {

	var sess AzureSession

	settings, err := getAzureAuth()
	if err != nil {
		return nil, fmt.Errorf("Error getting environment variables from Azure Auth file, Error :%v", err)
	}
	authorizer, err := auth.NewAuthorizerFromCLI()
	sess = AzureSession{
		SubscriptionID: settings.GetSubscriptionID(),
		Authorizer:     authorizer,
	}
	if err != nil {

		authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
		sess = AzureSession{
			SubscriptionID: settings.GetSubscriptionID(),
			Authorizer:     authorizer,
		}
		if err != nil {
			return nil, fmt.Errorf("can't initialize authorizer: %v", err)
		}
	}
	return &sess, nil
}
