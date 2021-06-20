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

func newRemoveCmd() *cobra.Command {

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove columns from CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnNames, _ := cmd.Flags().GetStringArray("column")
			outputPath, _ := cmd.Flags().GetString("output")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runRemove(
				format,
				inputPath,
				targetColumnNames,
				outputPath)
		},
	}

	removeCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	removeCmd.MarkFlagRequired("input")
	removeCmd.Flags().StringArrayP("column", "c", []string{}, "Name of the column to remove.")
	removeCmd.MarkFlagRequired("column")
	removeCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	removeCmd.MarkFlagRequired("output")

	return removeCmd
}

func runRemove(format csv.Format, inputPath string, targetColumnNames []string, outputPath string) error {

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	reader := csv.NewCsvReader(inputFile, format)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	writer := csv.NewCsvWriter(outputFile, format)

	err = remove(reader, targetColumnNames, writer)

	if err != nil {
		return err
	}

	return writer.Flush()
}

func remove(reader csv.CsvReader, removeColumnNames []string, writer csv.CsvWriter) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	removeColumnIndexes := []int{}
	for _, removeColumnName := range removeColumnNames {

		removeColumnIndex := util.IndexOf(columnNames, removeColumnName)
		if removeColumnIndex == -1 {
			return fmt.Errorf("missing %s in the CSV file", removeColumnName)
		}

		removeColumnIndexes = append(removeColumnIndexes, removeColumnIndex)
	}

	// 指定したカラム以外に絞るフィルタを定義
	filter := func(row []string) []string {

		filtered := []string{}

		for i, item := range row {

			if !util.Contains(removeColumnIndexes, i) {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	err = writer.Write(filter(columnNames))
	if err != nil {
		return err
	}

	// ヘッダ以外
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the CSV file")
		}

		err = writer.Write(filter(row))
		if err != nil {
			return err
		}
	}

	return nil
}
