package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	aksLong = `The AKS plugin is use to get the credentials of the Kubernetes cluster using kubectl CLI.`
)

// NewCmdAks provides a cobra command wrapping NamespaceOptions
func NewCmdAks(streams genericclioptions.IOStreams) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "aks SUBCOMMAND",
		Short:        "Manage Kubernetes Clusters from Kubectl.",
		Long:         aksLong,
		SilenceUsage: true,
		//RunE:         Help(),
	}
	cmd.AddCommand(NewCmdAksGet(streams))
	cmd.AddCommand(NewCmdAksList(streams))
	return cmd
}
