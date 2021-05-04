package cmd

import (
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use: "csvt",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newJoinCmd())
	rootCmd.AddCommand(newCountCmd())

	return rootCmd
}

func Execute() {

	rootCmd := newRootCmd()
	cobra.CheckErr(rootCmd.Execute())
}
