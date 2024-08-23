package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "helm-dev",
		Short: "Development tool for the development of Helm.",
		Long:  `A tool to help with the development of Helm. Meant to peak inside of Helm resources.`,
	}

	rootCmd.AddCommand(NewInspectResourceCmd())

	return rootCmd
}
