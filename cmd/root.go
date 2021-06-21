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

	rootCmd.PersistentFlags().StringP("delim", "", "", "(optional) CSV delimiter. The default is ','")
	rootCmd.PersistentFlags().StringP("quote", "", "", "(optional) CSV quote. The default is '\"'")
	rootCmd.PersistentFlags().StringP("sep", "", "", "(optional) CSV record separator. The default is CRLF.")
	rootCmd.PersistentFlags().BoolP("allquote", "", false, "(optional) Always quote CSV fields. The default is to quote only the necessary fields.")
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newJoinCmd())
	rootCmd.AddCommand(newCountCmd())
	rootCmd.AddCommand(newRemoveCmd())
	rootCmd.AddCommand(newChooseCmd())
	rootCmd.AddCommand(newHeaderCmd())
	rootCmd.AddCommand(newFilterCmd())
	rootCmd.AddCommand(newRenameCmd())
	rootCmd.AddCommand(newTransformCmd())

	for _, c := range rootCmd.Commands() {
		c.Flags().SortFlags = false
		c.InheritedFlags().SortFlags = false
	}

	return rootCmd
}

func Execute() {

	rootCmd := newRootCmd()
	cobra.CheckErr(rootCmd.Execute())
}
