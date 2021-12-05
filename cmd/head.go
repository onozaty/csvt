package cmd

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newHeadCmd() *cobra.Command {

	headCmd := &cobra.Command{
		Use:   "head",
		Short: "Show head few rows",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			number, _ := cmd.Flags().GetInt("number")

			// 表示件数は1以上
			if number <= 0 {
				return fmt.Errorf("number must be greater than or equal to 1")
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runHead(
				format,
				inputPath,
				number,
				cmd.OutOrStdout())
		},
	}

	headCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	headCmd.MarkFlagRequired("input")
	headCmd.Flags().IntP("number", "n", 10, "The number of records to show. If not specified, it will be the first 10 rows.")

	return headCmd
}

func runHead(format csv.Format, inputPath string, number int, writer io.Writer) error {

	reader, close, err := setupInput(inputPath, format)
	if err != nil {
		return err
	}
	defer close()

	return head(reader, number, writer)
}

func head(reader csv.CsvReader, number int, writer io.Writer) error {

	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the input CSV file")
	}

	table := tablewriter.NewWriter(writer)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader(columnNames)

	for i := 0; i < number; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the input CSV file")
		}

		table.Append(row)
	}

	table.Render()
	return nil
}
