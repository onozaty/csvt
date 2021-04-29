package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "csvt",
}

func Execute() {

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		_ = rootCmd.Help()
	}
	rootCmd.SilenceErrors = true

	cobra.CheckErr(rootCmd.Execute())
}
