package util

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
)

// AzureSession is an object representing session for subscription
type AzureSession struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
}

func getAzureAuth() (*map[string]interface{}, error) {

	defaultConfig, _ := homedir.Expand("~/.kube/azure.auth")
	if _, err := os.Stat(os.Getenv("AZURE_AUTH_LOCATION")); os.IsNotExist(err) {
		if _, err = os.Stat(defaultConfig); os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot get auth file: %v", err)
		}
		os.Setenv("AZURE_AUTH_LOCATION", defaultConfig)

	}

	data, err := ioutil.ReadFile(os.Getenv("AZURE_AUTH_LOCATION"))
	if err != nil {
		return nil, fmt.Errorf("cannot open the auth file: %v", err)
	}
	contents := make(map[string]interface{})
	err = json.Unmarshal(data, &contents)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal: %v", err)
	}

	return &contents, err
}

// NewSessionFromFile returns Azure Session Object
func NewSessionFromFile() (*AzureSession, error) {

	authInfo, err := getAzureAuth()

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	authorizer, err := auth.NewAuthorizerFromCLI()

	if err != nil {
		authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)

		if err != nil {
			return nil, fmt.Errorf("can't initialize authorizer: %v", err)
		}
	}

	sess := AzureSession{
		SubscriptionID: (*authInfo)["subscriptionId"].(string),
		Authorizer:     authorizer,
	}

	return &sess, nil
}
