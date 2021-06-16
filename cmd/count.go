package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/onozaty/csvt/csv"
	"github.com/onozaty/csvt/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newCountCmd() *cobra.Command {

	countCmd := &cobra.Command{
		Use:   "count",
		Short: "Count the number of records in CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			csvPath, _ := cmd.Flags().GetString("input")
			targetColumnName, _ := cmd.Flags().GetString("column")
			includeHeader, _ := cmd.Flags().GetBool("header")

			count, err := runCount(
				csvPath,
				CountOptions{
					targetColumnName: targetColumnName,
					includeHeader:    includeHeader,
				})

			if err != nil {
				return err
			}

			cmd.Printf("%d\n", count)

			return nil
		},
	}

	countCmd.Flags().StringP("input", "i", "", "CSV file path.")
	countCmd.MarkFlagRequired("input")
	countCmd.Flags().StringP("column", "c", "", "(optional) Name of the column to be counted. Only those with values will be counted.")
	countCmd.Flags().BoolP("header", "", false, "(optional) Counting including header. The default is to exclude header.")
	countCmd.Flags().SortFlags = false

	return countCmd
}

type CountOptions struct {
	targetColumnName string
	includeHeader    bool
}

func runCount(csvPath string, options CountOptions) (int, error) {

	file, err := os.Open(csvPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewCsvReader(file)
	return count(reader, options)
}

func count(reader csv.CsvReader, options CountOptions) (int, error) {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return 0, errors.Wrap(err, "failed to read the CSV file")
	}

	targetColumnIndex := -1
	if options.targetColumnName != "" {
		targetColumnIndex = util.IndexOf(columnNames, options.targetColumnName)
		if targetColumnIndex == -1 {
			return 0, fmt.Errorf("missing %s in the CSV file", options.targetColumnName)
		}
	}

	count := 0
	if options.includeHeader {
		count++
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, errors.Wrap(err, "failed to read the CSV file")
		}

		if targetColumnIndex == -1 || row[targetColumnIndex] != "" {
			count++
		}
	}

	return count, nil
}
