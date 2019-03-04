package util

import (
	"context"
	"fmt"
	container "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

// KubernertesConfig struct
type KubernertesConfig struct {
	APIVersion     string     `yaml:"apiVersion"`
	Clusters       []Clusters `yaml:"clusters"`
	Contexts       []Contexts `yaml:"contexts"`
	CurrentContext string     `yaml:"current-context"`
	Kind           string     `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []Users `yaml:"users"`
}

// Clusters config
type Clusters struct {
	Cluster struct {
		CertificateAuthorityData string `yaml:"certificate-authority-data"`
		Server                   string `yaml:"server"`
	} `yaml:"cluster"`
	Name string `yaml:"name"`
}

// Contexts config
type Contexts struct {
	Context struct {
		Cluster string `yaml:"cluster"`
		User    string `yaml:"user"`
	} `yaml:"context"`
	Name string `yaml:"name"`
}

// Users config
type Users struct {
	User struct {
		ClientCertificateData string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
		Token                 string `yaml:"token"`
	} `yaml:"user"`
	Name string `yaml:"name"`
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

func unmarshalYaml(path string) (*KubernertesConfig, error) {

	var k8sconfig KubernertesConfig

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	_ = yaml.Unmarshal(data, &k8sconfig)

	return &k8sconfig, err
}

func mergeConfig(temp string, existing string) ([]byte, error) {

	var file []byte
	existingFile, err := unmarshalYaml(existing)
	if err != nil {
		return file, fmt.Errorf("error: %v", err)
	}
	tempFile, err := unmarshalYaml(temp)
	if err != nil {
		return file, fmt.Errorf("error: %v", err)
	}

	for _, ctx := range tempFile.Contexts {
		if strings.HasPrefix(ctx.Context.User, "clusterAdmin") {
			adminName := fmt.Sprintf("%v-admin", ctx.Name)
			tempFile.CurrentContext, ctx.Name = adminName, adminName
		}
	}
	if tempFile.APIVersion == "" {
		return file, fmt.Errorf("Failed to configuration from %v", temp)

	}
	if existingFile.APIVersion == "" {
		existingFile = tempFile
	} else {

		mergeClusters(tempFile, existingFile)
		mergeContexts(tempFile, existingFile)
		mergeUsers(tempFile, existingFile)

		existingFile.CurrentContext = tempFile.CurrentContext
	}

	file, err = yaml.Marshal(existingFile)

	return file, err
}

func mergeClusters(temp *KubernertesConfig, existing *KubernertesConfig) {

	if len(temp.Clusters) != 0 {
		if len(existing.Clusters) == 0 {
			fmt.Println("its Empty")
			existing.Clusters = temp.Clusters
		}
		for _, i := range temp.Clusters {
			for key, j := range existing.Clusters {
				if i.Name == j.Name {
					if i == j {
						// Clusters with same name exist, deleting existing cluster and replace with the newly downloaded config
						existing.Clusters = func(s []Clusters, i int) []Clusters {
							s[i] = s[len(s)-1]
							return s[:len(s)-1]
						}(existing.Clusters, key)
					} else {
						fmt.Printf("A different object named %v already exists in %v", i.Name, "Clusters")
					}
				}
			}
			existing.Clusters = append(existing.Clusters, i)
		}
	}

}

func mergeContexts(temp *KubernertesConfig, existing *KubernertesConfig) {

	if len(temp.Contexts) != 0 {
		if len(existing.Contexts) == 0 {
			existing.Contexts = temp.Contexts
		}
		for _, i := range temp.Contexts {
			for key, j := range existing.Contexts {
				if i.Name == j.Name {
					if i == j {
						// Context with same name exist, deleting existing context and replace with the newly downloaded config
						existing.Contexts = func(s []Contexts, i int) []Contexts {
							s[i] = s[len(s)-1]
							return s[:len(s)-1]
						}(existing.Contexts, key)
					} else {
						fmt.Printf("A different object named %v already exists in %v", i.Name, "Context")
					}
				}
			}
			existing.Contexts = append(existing.Contexts, i)
		}
	}

}

func mergeUsers(temp *KubernertesConfig, existing *KubernertesConfig) {

	if len(temp.Users) != 0 {
		if len(existing.Users) == 0 {
			existing.Users = temp.Users
		}
		for _, i := range temp.Users {
			for key, j := range existing.Users {
				if i.Name == j.Name {
					if i == j {
						// User with same name exist, deleting existing user and replace with the newly downloaded config
						existing.Users = func(s []Users, i int) []Users {
							s[i] = s[len(s)-1]
							return s[:len(s)-1]
						}(existing.Users, key)
					} else {
						fmt.Printf("A different object named %v already exists in %v", i.Name, "User")
					}
				}
			}
			existing.Users = append(existing.Users, i)
		}
	}
}

// ManageConfig is use to merge kubeconfiguration with existing config
func ManageConfig(config string, path string) error {

	var err error
	if path == "-" {
		fmt.Println(config)
		return nil
	}
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

	// Merge configuration of temporary file with existing kubernetes configuration (default: ~/.kube/config)
	configFile, err := mergeConfig(tempFile, path)
	err = ioutil.WriteFile(path, configFile, 0600)
	if err != nil {
		fmt.Println(err)
	}

	return err
}
