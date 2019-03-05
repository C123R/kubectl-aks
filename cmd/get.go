// Copyright © 2019 NAME HERE cizer.ciz@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/kubectl-aks/util"
	"github.com/mitchellh/go-homedir"
	"os"
	//	log "github.com/sirupsen/logrus"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"time"
)

// name Flag to accept cluster name
var name string

// path Flag to accept path for the Kubeconfig
var path string

// getCmd represents the get command

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Kubernetes credentials from Azure and add it to ~/.kube/config",
	Long: `
The AKS plugin's get command download the kubernetes credentials from Azure
and merge it with the default ~/.kube/config.

For example:

	$ kubectl aks get -n foo-cluster

	You can get the list of AKS cluster in Azure Subscription using "kubectl aks list".
`,
	RunE: get,
}

func init() {

	rootCmd.AddCommand(getCmd)
	defaultConfig, _ := homedir.Expand("~/.kube/config")
	getCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the AKS Cluster (required)")
	getCmd.MarkFlagRequired("name")
	getCmd.Flags().StringVarP(&path, "path", "p", defaultConfig, "Path to write Kubeconfig")

}

func get(cmd *cobra.Command, args []string) error {

	sess, err := util.NewSessionFromFile()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Start()
	s.Prefix = fmt.Sprintf("Getting credentials for AKS cluster \"%v\"", name)
	s.Writer = os.Stderr
	s.FinalMSG = fmt.Sprintf("Merged \"%v\" as current context in %v\n", name, path)

	kubeconfig, err := util.GetAKS(sess, name)
	if err != nil {
		fmt.Println()
		return fmt.Errorf("%v", err)
	}
	err = util.ManageConfig(kubeconfig, path)
	if err != nil {
		fmt.Println()
		return err
	}
	s.Stop()
	return err

}
