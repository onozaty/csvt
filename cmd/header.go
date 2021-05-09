package cmd

import (
	"os"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newHeaderCmd() *cobra.Command {

	countCmd := &cobra.Command{
		Use:   "header",
		Short: "Show the header of CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			inputCsvPath, _ := cmd.Flags().GetString("input")

			columnNames, err := runHeader(inputCsvPath)

			if err != nil {
				return err
			}

			for _, columnName := range columnNames {
				cmd.Printf("%s\n", columnName)
			}

			return nil
		},
	}

	countCmd.Flags().StringP("input", "i", "", "CSV file path.")
	countCmd.MarkFlagRequired("input")
	countCmd.Flags().SortFlags = false

	return countCmd
}

func runHeader(inputCsvPath string) ([]string, error) {

	file, err := os.Open(inputCsvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := csv.NewCsvReader(file)
	if err != nil {
		return nil, err
	}

	return header(reader)
}

func header(reader csv.CsvReader) ([]string, error) {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the CSV file")
	}

	return columnNames, nil
}
