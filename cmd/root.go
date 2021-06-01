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

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newJoinCmd())
	rootCmd.AddCommand(newCountCmd())
	rootCmd.AddCommand(newRemoveCmd())
	rootCmd.AddCommand(newChooseCmd())
	rootCmd.AddCommand(newHeaderCmd())
	rootCmd.AddCommand(newFilterCmd())
	rootCmd.AddCommand(newRenameCmd())

	return rootCmd
}

func Execute() {

	rootCmd := newRootCmd()
	cobra.CheckErr(rootCmd.Execute())
}
