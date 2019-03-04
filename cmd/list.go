// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List AKS cluster from the current Azure Subscrption",
	Long:  "",
	Run:   list,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func list(cmd *cobra.Command, args []string) {

	sess, err := util.NewSessionFromFile()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	var akslist util.AksCluster
	aksList, err := akslist.ListAKS(sess)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%v\t\t\t%v\t\t%v\t\n", "NAME", "KUBERNETES VERSION", "RESOURCE GROUP")
	for key, value := range aksList {

		fmt.Fprintf(os.Stdout, "%v\t\t%v\t\t\t\t%v\t\n", key, value.K8sVersion, value.ResourceGroup)
	}

}
