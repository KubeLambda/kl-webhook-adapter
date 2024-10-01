package cmd

import (
	"kl-webhook-adapter/internal/infra"

	"github.com/spf13/cobra"
)

var (
	deployment string
	serverCmd  = &cobra.Command{
		Use:   "run --deployment={local|dev|prod|...}",
		Short: "Run my service",
		Long: "Run my service and use configuration settings specified in 'deployment' flag. Name of deployment " +
			"corresponds to filename in config directory, e.g. 'run --deployment=local' means that service " +
			"will be started with config values loaded from configs/local.yaml file",
		Run: func(cmd *cobra.Command, args []string) {
			infra.Start(deployment)
		},
	}
)

func init() {
	serverCmd.Flags().StringVar(&deployment, "deployment", "",
		"deployment environment for API server, e.g. local, prod (it should match your configuration filename)")
	_ = serverCmd.MarkFlagRequired("deployment")
	rootCmd.AddCommand(serverCmd)
}
