package cmd

import (
	"fmt"
	"github.com/C123R/kubectl-aks/util"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

var (
	listAKSLong = `The AKS plugin's list subcommand get the list of AKS clusters from current
Azure subscription.
	`
)

// NewCmdAksList provides a cobra command wrapping NamespaceOptions
func NewCmdAksList(streams genericclioptions.IOStreams) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List AKS cluster from the current Azure Subscrption.",
		Long:  listAKSLong,

		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("list subcommend does not support any arguments")
			}
			if err := getList(); err != nil {
				return fmt.Errorf("error getting list of AKS Clusters,Error: %v", err)
			}
			return nil
		},
	}
	return cmd
}

func getList() error {

	aksClient, err := util.NewAKSClient()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	aksList, err := aksClient.ListAKS()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	fmt.Printf("%v\t\t\t%v\t\t%v\t\t%v\n", "NAME", "VERSION", "NODES", "RESOURCE GROUP")
	for key, value := range aksList {
		fmt.Fprintf(os.Stdout, "%v\t\t%v\t\t%v\t\t%v\n", key, value.K8sVersion, value.Nodes, value.ResourceGroup)
	}

	return nil

}
