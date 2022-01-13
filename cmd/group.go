package cmd

import (
	"io"
	_sort "sort"
	"strconv"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newGroupCmd() *cobra.Command {

	gcountCmd := &cobra.Command{
		Use:   "group",
		Short: "Aggregate by group",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnName, _ := cmd.Flags().GetString("column")
			outputPath, _ := cmd.Flags().GetString("output")
			countColumnName, _ := cmd.Flags().GetString("count-column")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runGroupCount(
				format,
				inputPath,
				targetColumnName,
				countColumnName,
				outputPath)
		},
	}

	gcountCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	gcountCmd.MarkFlagRequired("input")
	gcountCmd.Flags().StringP("column", "c", "", "Name of the column to use for grouping.")
	gcountCmd.MarkFlagRequired("column")
	gcountCmd.Flags().StringP("count-column", "", "COUNT", "(optional) Column name for the number of records.")
	gcountCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	gcountCmd.MarkFlagRequired("output")

	return gcountCmd
}

func runGroupCount(format csv.Format, inputPath string, targetColumnName string, countColumnName string, outputPath string) error {

	reader, writer, close, err := setupInputOutput(inputPath, outputPath, format)
	if err != nil {
		return err
	}
	defer close()

	err = groupCount(reader, targetColumnName, countColumnName, writer)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func groupCount(reader csv.CsvReader, targetColumnName string, countColumnName string, writer csv.CsvWriter) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	targetColumnIndex, err := getTargetColumnIndex(columnNames, targetColumnName)
	if err != nil {
		return err
	}

	counter := map[string]int{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the CSV file")
		}

		val := row[targetColumnIndex]
		counter[val] = counter[val] + 1
	}

	if err := writer.Write([]string{targetColumnName, countColumnName}); err != nil {
		return err
	}

	// グループ化した値でソートして出力
	keys := []string{}
	for k := range counter {
		keys = append(keys, k)
	}
	_sort.Strings(keys)

	for _, k := range keys {
		if err := writer.Write([]string{k, strconv.Itoa(counter[k])}); err != nil {
			return err
		}
	}

	return nil
}
