package cmd

import (
	"context"
	"fmt"
	"github.com/C123R/kubectl-aks/util"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	upgradeAKSLong = `Upgrade a managed Kubernetes cluster to a newer version.`
)

// NewAKSUpgradeOptions provides an instance of AksUpgradeOptions with default values
func NewAKSUpgradeOptions(streams genericclioptions.IOStreams) *AksUpgradeOptions {

	return &AksUpgradeOptions{
		IOStreams: streams,
	}
}

// AksUpgradeOptions provides information required to get the AKS context from Azure.
type AksUpgradeOptions struct {
	userSpecifiedCluster string
	aksClient            util.AKSClient
	genericclioptions.IOStreams
	args    []string
	version string
}

// NewCmdAksUpgrade provides a cobra command wrapping AksUpgradeOptions
func NewCmdAksUpgrade(streams genericclioptions.IOStreams) *cobra.Command {

	o := NewAKSUpgradeOptions(streams)

	cmd := &cobra.Command{
		Use:          "upgrade",
		Short:        "Upgrade a managed Kubernetes cluster to a newer version.",
		Long:         upgradeAKSLong,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Upgrade(); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&o.version, "version", "v", "", "Version of Kubernetes to upgrade the cluster to")
	return cmd
}

// Complete sets all required information for kubernetes aks plugin
func (o *AksUpgradeOptions) Complete(cmd *cobra.Command, args []string) error {

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
	if o.version == "" {
		cmd.Usage()
		fmt.Println()
		return fmt.Errorf("You must specify the version: --version/-v")
	}
	o.aksClient, err = util.NewAKSClient()
	if err != nil {
		return fmt.Errorf("error authenticating with azure,Error: %v", err)
	}
	return nil
}

// Upgrade Kubernetes credentials for AKS cluster
func (o *AksUpgradeOptions) Upgrade() error {

	fmt.Println("Kubernetes may be unavailable during cluster upgrades.")
	if util.Confirmation() {
		err := o.aksClient.UpgradeAKS(context.Background(), o.userSpecifiedCluster, o.version)
		if err != nil {
			return err
		}
		return nil
	}
	fmt.Println("Opeartion Cancelled")
	return nil

}
