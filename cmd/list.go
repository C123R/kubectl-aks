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
	"github.com/C123R/kubectl-aks/util"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List AKS cluster from the current Azure Subscrption",
	Long:  "",
	RunE:  list,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {

	sess, err := util.NewSessionFromFile()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Start()
	s.Prefix = fmt.Sprintf("Getting List of AKS clusters for your current Subscription")
	s.FinalMSG = fmt.Sprintf("%v\t\t\t%v\t\t%v\n", "NAME", "VERSION", "RESOURCE GROUP")

	var akslist util.AksCluster
	aksList, err := akslist.ListAKS(sess)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	s.Stop()
	for key, value := range aksList {

		fmt.Fprintf(os.Stdout, "%v\t\t%v\t\t%v\n", key, value.K8sVersion, value.ResourceGroup)
	}

	return err

}
