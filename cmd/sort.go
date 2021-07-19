package cmd

import (
	"github.com/onozaty/csvt/csv"
	"github.com/spf13/cobra"
)

func newSortCmd() *cobra.Command {

	sortCmd := &cobra.Command{
		Use:   "sort",
		Short: "Sort rows",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnNames, _ := cmd.Flags().GetStringArray("column")
			sortDescending, _ := cmd.Flags().GetBool("desc")
			asNumber, _ := cmd.Flags().GetBool("number")
			useFileRows, _ := cmd.Flags().GetBool("usingfile")
			outputPath, _ := cmd.Flags().GetString("output")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runSort(
				format,
				inputPath,
				targetColumnNames,
				outputPath,
				SortOptions{
					sortDescending: sortDescending,
					asNumber:       asNumber,
					useFileRows:    useFileRows,
				})
		},
	}

	sortCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	sortCmd.MarkFlagRequired("input")
	sortCmd.Flags().StringArrayP("column", "c", []string{}, "Name of the column to use for sorting.")
	sortCmd.MarkFlagRequired("column")
	sortCmd.Flags().BoolP("desc", "", false, "(optional) Sort in descending order. The default is ascending order.")
	sortCmd.Flags().BoolP("number", "", false, "(optional) Sorts as a number. The default is to sort as a string.")
	sortCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	sortCmd.MarkFlagRequired("output")
	sortCmd.Flags().BoolP("usingfile", "", false, "(optional) Use temporary files for sorting. Use this when sorting large files that will not fit in memory.")

	return sortCmd
}

type SortOptions struct {
	sortDescending bool
	asNumber       bool
	useFileRows    bool
}

func runSort(format csv.Format, inputPath string, targetColumnNames []string, outputPath string, options SortOptions) error {

	reader, writer, close, err := setupInputOutput(inputPath, outputPath, format)
	if err != nil {
		return err
	}
	defer close()

	err = sort(reader, targetColumnNames, writer, options)

	if err != nil {
		return err
	}

	return writer.Flush()
}

func sort(reader csv.CsvReader, targetColumnNames []string, writer csv.CsvWriter, options SortOptions) error {

	var compare func(item1 string, item2 string) (int, error)

	if options.asNumber {
		compare = csv.CompareNumber
	} else {
		compare = csv.CompareString
	}

	if options.sortDescending {
		compare = csv.Descending(compare)
	}

	var sortedRows csv.CsvSortedRows
	var err error
	if options.useFileRows {
		sortedRows, err = csv.LoadCsvFileSortedRows(reader, targetColumnNames, compare)
	} else {
		sortedRows, err = csv.LoadCsvMemorySortedRows(reader, targetColumnNames, compare)
	}
	if err != nil {
		return err
	}

	err = writer.Write(sortedRows.ColumnNames())
	if err != nil {
		return err
	}

	for i := 0; i < sortedRows.Count(); i++ {

		row, err := sortedRows.Row(i)
		if err != nil {
			return err
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
