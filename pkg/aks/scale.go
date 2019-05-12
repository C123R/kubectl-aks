package cmd

import (
	"context"
	"fmt"
	"github.com/C123R/kubectl-aks/util"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	scaleAKSLong = `Scale the node pool in a managed Kubernetes cluster.`
)

// NewAKSScaleOptions provides an instance of AksScaleOptions with default values
func NewAKSScaleOptions(streams genericclioptions.IOStreams) *AksScaleOptions {

	return &AksScaleOptions{
		IOStreams: streams,
	}
}

// AksScaleOptions provides information required to get the AKS context from Azure.
type AksScaleOptions struct {
	userSpecifiedCluster string
	aksClient            util.AKSClient
	genericclioptions.IOStreams
	args  []string
	count int32
}

// NewCmdAksScale provides a cobra command wrapping NamespaceOptions
func NewCmdAksScale(streams genericclioptions.IOStreams) *cobra.Command {

	o := NewAKSScaleOptions(streams)
	cmd := &cobra.Command{
		Use:   "scale",
		Short: "Scale the node pool in a managed Kubernetes cluster.",
		Long:  scaleAKSLong,

		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Scale(); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().Int32VarP(&o.count, "count", "c", 0, "Number of nodes in the Kubernetes node pool")
	return cmd
}

// Complete sets all required information for kubernetes aks plugin
func (o *AksScaleOptions) Complete(cmd *cobra.Command, args []string) error {

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
	if o.count == 0 {
		cmd.Usage()
		fmt.Println()
		return fmt.Errorf("You must specify the count of nodes: --count/-c")
	}

	o.aksClient, err = util.NewAKSClient()
	if err != nil {
		return fmt.Errorf("error authenticating with azure,Error: %v", err)
	}
	return nil
}

// Scale Kubernetes credentials for AKS cluster
func (o *AksScaleOptions) Scale() error {

	if util.Confirmation() {
		err := o.aksClient.ScaleAKS(context.Background(), o.userSpecifiedCluster, o.count)
		if err != nil {
			return err
		}
		return nil
	}
	fmt.Println("Opeartion Cancelled")
	return nil
}
