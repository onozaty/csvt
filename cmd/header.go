package cmd

import (
	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newHeaderCmd() *cobra.Command {

	countCmd := &cobra.Command{
		Use:   "header",
		Short: "Show header",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputCsvPath, _ := cmd.Flags().GetString("input")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			columnNames, err := runHeader(format, inputCsvPath)

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

	return countCmd
}

func runHeader(format csv.Format, inputPath string) ([]string, error) {

	reader, close, err := setupInput(inputPath, format)
	if err != nil {
		return nil, err
	}
	defer close()

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
