package cmd

import (
	"fmt"
	"io"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newConcatCmd() *cobra.Command {

	concatCmd := &cobra.Command{
		Use:   "concat",
		Short: "Concat CSV files",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			firstPath, _ := cmd.Flags().GetString("first")
			secondPath, _ := cmd.Flags().GetString("second")
			outputPath, _ := cmd.Flags().GetString("output")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runConcat(format, firstPath, secondPath, outputPath)
		},
	}

	concatCmd.Flags().StringP("first", "1", "", "First CSV file path.")
	concatCmd.MarkFlagRequired("first")
	concatCmd.Flags().StringP("second", "2", "", "Second CSV file path.")
	concatCmd.MarkFlagRequired("second")
	concatCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	concatCmd.MarkFlagRequired("output")

	return concatCmd
}

func runConcat(format csv.Format, firstPath string, secondPath string, outputPath string) error {

	firstReader, firstClose, err := setupInput(firstPath, format)
	if err != nil {
		return err
	}
	defer firstClose()

	secondReader, secondClose, err := setupInput(secondPath, format)
	if err != nil {
		return err
	}
	defer secondClose()

	writer, outputClose, err := setupOutput(outputPath, format)
	if err != nil {
		return err
	}
	defer outputClose()

	err = concat(firstReader, secondReader, writer)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func concat(first csv.CsvReader, second csv.CsvReader, writer csv.CsvWriter) error {

	firstColumnNames, err := first.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the first CSV file")
	}

	secondColumnNames, err := second.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the second CSV file")
	}

	if len(firstColumnNames) != len(secondColumnNames) {
		return fmt.Errorf("number of columns does not match")
	}

	// 1つ目のCSVのカラム名と2つ目のCSVのカラム名のマッピングを作成
	secondColumnIndexes := []int{}
	for _, firstColumnName := range firstColumnNames {

		secondColumnIndex, err := getTargetColumnIndex(secondColumnNames, firstColumnName)
		if err != nil {
			return errors.Wrap(err, "no column corresponding to the second CSV file")
		}

		secondColumnIndexes = append(secondColumnIndexes, secondColumnIndex)
	}

	// 1つ目の書き込み
	writer.Write(firstColumnNames)
	for {
		row, err := first.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the first CSV file")
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	// 2つ目の書き込み
	for {
		row, err := second.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the second CSV file")
		}

		// 1つ目のCSVに合わせてカラム入れ替え
		swapedRow := []string{}
		for _, secondColumnIndex := range secondColumnIndexes {
			swapedRow = append(swapedRow, row[secondColumnIndex])
		}

		err = writer.Write(swapedRow)
		if err != nil {
			return err
		}
	}

	return nil
}
