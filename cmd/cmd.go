package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = newCobraCmdlineCommand(
	"apiserver", // TODO replace to your application name
	"My new server",
	"My new server built from GoQuick template",
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newCobraCmdlineCommand is a helper function to add new command-line command and parameters
func newCobraCmdlineCommand(use string, short string, long string) cobra.Command {
	return cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			bindViperToCobraCommands([]*cobra.Command{cmd})
		},
	}
}

func bindViperToCobraCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		copyViperToCobraFlags(cmd)
		if cmd.HasSubCommands() {
			bindViperToCobraCommands(cmd.Commands())
		}
	}
}

func copyViperToCobraFlags(cmd *cobra.Command) {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			if err := cmd.Flags().Set(f.Name, viper.GetString(f.Name)); err != nil {
				panic(err)
			}
		}
	})
}
