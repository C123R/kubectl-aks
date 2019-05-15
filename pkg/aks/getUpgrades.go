package cmd

import (
	"fmt"

	"github.com/C123R/kubectl-aks/util"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	//"os"
)

var (
	getUpgradesAKSLong = `Get the upgrade versions available for a managed Kubernetes cluster.`
)

// NewAKSGetUpgradesOptions provides an instance of AksGetUpgradesOptions with default values
func NewAKSGetUpgradesOptions(streams genericclioptions.IOStreams) *AksGetUpgradesOptions {

	return &AksGetUpgradesOptions{
		IOStreams: streams,
	}
}

// AksGetUpgradesOptions provides information required to get the AKS context from Azure.
type AksGetUpgradesOptions struct {
	userSpecifiedCluster string
	aksClient            util.AKSClient
	genericclioptions.IOStreams
	args []string
}

// NewCmdAksGetUpgrades provides a cobra command wrapping NamespaceOptions
func NewCmdAksGetUpgrades(streams genericclioptions.IOStreams) *cobra.Command {

	o := NewAKSGetUpgradesOptions(streams)
	cmd := &cobra.Command{
		Use:   "get-upgrades",
		Short: "Get the upgrade versions available for a managed Kubernetes cluster.",
		Long:  scaleAKSLong,

		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.GetUpgrades(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

// Complete sets all required information for kubernetes aks plugin
func (o *AksGetUpgradesOptions) Complete(cmd *cobra.Command, args []string) error {

	o.args = args
	var err error
	// validating whether all required arguments are provided
	if len(o.args) == 0 {
		cmd.Usage()
		fmt.Println()
		return fmt.Errorf("You must specify the name of AKS cluster. Use \"kubectl aks list\" to get the list of AKS clusters")
	}
	if len(o.args) > 0 {
		o.userSpecifiedCluster = args[0]
	}
	o.aksClient, err = util.NewAKSClient()
	if err != nil {
		return fmt.Errorf("error authenticating with azure,Error: %v", err)
	}
	return nil
}

// GetUpgrades Kubernetes credentials for AKS cluster
func (o *AksGetUpgradesOptions) GetUpgrades() error {

	k8sUpgradeVersions, currentVersion, err := o.aksClient.GetUpgrades(o.userSpecifiedCluster)
	if err != nil {
		return err
	}
	if len(k8sUpgradeVersions) == 0 {
		fmt.Printf("Currently there are no new upgrades avaialble, %v is upto date [%v].\n", o.userSpecifiedCluster, currentVersion)
	} else {
		fmt.Printf("Current Version: %v\n\n", currentVersion)
		fmt.Printf("List of avaliable upgrades for %v:\n", o.userSpecifiedCluster)
		for _, v := range k8sUpgradeVersions {
			fmt.Println(v)
		}
	}
	return nil
}
