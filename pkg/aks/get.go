package cmd

import (
	"fmt"
	"github.com/C123R/kubectl-aks/util"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	getAksLong = `The AKS plugin's get subcommand download the kubernetes credentials from Azure 
and merge it with the default ~/.kube/config.`
	getAksExample = `  # Get the credentials AKS Cluster from current Azure Subscription.
  kubectl aks get -n foo-cluster
	`
	path string
)

// AksGetOptions provides information required to get the AKS context from Azure.
type AksGetOptions struct {
	userSpecifiedCluster string
	azureSession         *util.AzureSession
	genericclioptions.IOStreams
	args []string
}

// NewAKSGetOptions provides an instance of NamespaceOptions with default values
func NewAKSGetOptions(streams genericclioptions.IOStreams) *AksGetOptions {

	return &AksGetOptions{
		IOStreams: streams,
	}
}

// NewCmdAksGet provides a cobra command wrapping NamespaceOptions
func NewCmdAksGet(streams genericclioptions.IOStreams) *cobra.Command {

	o := NewAKSGetOptions(streams)

	cmd := &cobra.Command{
		Use:          "get CLUSTER_NAME",
		Short:        "Get Kubernetes cluster configuration.",
		Long:         getAksLong,
		Example:      getAksExample,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Get(); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", util.DefaultKubeConfig, "Path to write Kubeconfig")
	return cmd
}

// Complete sets all required information for kubernetes aks plugin
func (o *AksGetOptions) Complete(cmd *cobra.Command, args []string) error {

	o.args = args
	var err error
	// validating whether all required arguments are provided
	if len(o.args) == 0 {
		return cmd.Usage()
	}
	if len(o.args) > 0 {
		o.userSpecifiedCluster = args[0]
	}
	o.azureSession, err = util.NewSessionFromFile()
	if err != nil {
		return fmt.Errorf("error authenticating with azure,Error: %v", err)
	}
	return nil
}

// Get Kubernetes credentials for AKS cluster
func (o *AksGetOptions) Get() error {

	config, err := util.GetAKS(o.azureSession, o.userSpecifiedCluster)
	if err != nil {
		return fmt.Errorf("error getting kubernetes configuration for cluster %v,Error: %v", o.userSpecifiedCluster, err)
	}

	err = util.MergeConfig(config, path)
	if err != nil {
		return fmt.Errorf("error merging kubernetes configuration for cluster %v with %v ,Error: %v", o.userSpecifiedCluster, path, err)
	}

	fmt.Printf("Merged \"%v\" as current context in %v\n", o.userSpecifiedCluster, path)
	return nil
}
