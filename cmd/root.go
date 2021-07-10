package cmd

import (
	"fmt"

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
	rootCmd.PersistentFlags().StringP("encoding", "", "", "(optional) CSV encoding. The default is utf-8. Supported encodings: utf-8, shift_jis, euc-jp")
	rootCmd.PersistentFlags().BoolP("bom", "", false, "(optional) CSV with BOM. When reading, the BOM will be automatically removed without this flag.")
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
	rootCmd.AddCommand(newReplaceCmd())
	rootCmd.AddCommand(newUniqueCmd())
	rootCmd.AddCommand(newIncludeCmd())
	rootCmd.AddCommand(newExcludeCmd())
	rootCmd.AddCommand(newConcatCmd())
	rootCmd.AddCommand(newSliceCmd())

	for _, c := range rootCmd.Commands() {
		// フラグ以外は受け付けないように
		c.Args = func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("only flags can be specified")
			}
			return nil
		}
		c.Flags().SortFlags = false
		c.InheritedFlags().SortFlags = false
	}

	return rootCmd
}

func Execute() {

	rootCmd := newRootCmd()
	cobra.CheckErr(rootCmd.Execute())
}
