package cmd

import (
	"fmt"
	"io"
	"math"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newSliceCmd() *cobra.Command {

	sliceCmd := &cobra.Command{
		Use:   "slice",
		Short: "Slice specified range of rows",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			start, _ := cmd.Flags().GetInt("start")
			end, _ := cmd.Flags().GetInt("end")
			outputPath, _ := cmd.Flags().GetString("output")

			// 開始は1以上
			if start <= 0 {
				return fmt.Errorf("start must be greater than or equal to 1")
			}

			// 開始と終了が逆
			if start > end {
				return fmt.Errorf("end must be greater than or equal to start")
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runSlice(
				format,
				inputPath,
				start,
				end,
				outputPath)
		},
	}

	sliceCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	sliceCmd.MarkFlagRequired("input")
	sliceCmd.Flags().IntP("start", "s", 1, "The number of the starting row. If not specified, it will be the first row.")
	sliceCmd.Flags().IntP("end", "e", math.MaxInt32, "The number of the end row. If not specified, it will be the last row.")
	sliceCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	sliceCmd.MarkFlagRequired("output")

	return sliceCmd
}

func runSlice(format csv.Format, inputPath string, start int, end int, outputPath string) error {

	reader, writer, close, err := setupInputOutput(inputPath, outputPath, format)
	if err != nil {
		return err
	}
	defer close()

	err = slice(reader, start, end, writer)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func slice(reader csv.CsvReader, start int, end int, writer csv.CsvWriter) error {

	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the input CSV file")
	}

	err = writer.Write(columnNames)
	if err != nil {
		return err
	}

	currentRowNum := 0

	for currentRowNum <= end {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the input CSV file")
		}

		currentRowNum++

		if currentRowNum >= start && currentRowNum <= end {
			err = writer.Write(row)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
