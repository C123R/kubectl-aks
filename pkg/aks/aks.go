package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	aksLong = `The AKS plugin manages AKS Clusters.`
)

// NewCmdAks provides a cobra command wrapping AksOptions
func NewCmdAks(streams genericclioptions.IOStreams) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "aks",
		Short:        "Manage Kubernetes Clusters from Kubectl.",
		Long:         aksLong,
		SilenceUsage: true,
	}
	cmd.AddCommand(NewCmdAksGet(streams))
	cmd.AddCommand(NewCmdAksList(streams))
	cmd.AddCommand(NewCmdAksScale(streams))
	cmd.AddCommand(NewCmdAksUpgrade(streams))
	cmd.AddCommand(NewCmdAksGetUpgrades(streams))
	return cmd
}
