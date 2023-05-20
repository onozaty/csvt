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

			inputPaths, _ := cmd.Flags().GetStringArray("input")
			outputPath, _ := cmd.Flags().GetString("output")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runConcat(format, inputPaths, outputPath)
		},
	}

	concatCmd.Flags().StringArrayP("input", "i", []string{}, "Input CSV files path.")
	concatCmd.MarkFlagRequired("input")
	concatCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	concatCmd.MarkFlagRequired("output")

	return concatCmd
}

func runConcat(format csv.Format, inputPaths []string, outputPath string) error {

	readers := []csv.CsvReader{}

	for _, inputPath := range inputPaths {
		reader, inputClose, err := setupInput(inputPath, format)
		if err != nil {
			return err
		}
		defer inputClose()

		readers = append(readers, reader)
	}

	writer, outputClose, err := setupOutput(outputPath, format)
	if err != nil {
		return err
	}
	defer outputClose()

	err = concat(readers, writer)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func concat(readers []csv.CsvReader, writer csv.CsvWriter) error {

	firstReader := readers[0]
	firstColumnNames, err := firstReader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the first CSV file")
	}

	// 1つ目の書き込み
	err = writer.Write(firstColumnNames)
	if err != nil {
		return err
	}

	for {
		row, err := firstReader.Read()
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

	// 2つ目以降
	count := 1
	for _, reader := range readers[1:] {
		count++

		columnNames, err := reader.Read()
		if err != nil {
			return errors.Wrapf(err, "failed to read CSV file (%d)", count)
		}

		if len(firstColumnNames) != len(columnNames) {
			return fmt.Errorf("number of columns does not match (%d)", count)
		}

		// 1つ目のCSVのカラム名とのカラム名のマッピングを作成
		columnIndexes := []int{}
		for _, firstColumnName := range firstColumnNames {

			columnIndex, err := getTargetColumnIndex(columnNames, firstColumnName)
			if err != nil {
				return errors.Wrapf(err, "no column corresponding in CSV file (%d)", count)
			}

			columnIndexes = append(columnIndexes, columnIndex)
		}

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return errors.Wrapf(err, "failed to read CSV file (%d)", count)
			}

			// 1つ目のCSVに合わせてカラム入れ替え
			swapedRow := []string{}
			for _, columnIndex := range columnIndexes {
				swapedRow = append(swapedRow, row[columnIndex])
			}

			err = writer.Write(swapedRow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
