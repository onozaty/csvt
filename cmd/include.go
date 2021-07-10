package cmd

import (
	"fmt"
	"io"

	"github.com/onozaty/csvt/csv"
	"github.com/onozaty/csvt/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newIncludeCmd() *cobra.Command {

	includeCmd := &cobra.Command{
		Use:   "include",
		Short: "Filter rows by included in another CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnName, _ := cmd.Flags().GetString("column")
			anotherPath, _ := cmd.Flags().GetString("another")
			anotherColumnName, _ := cmd.Flags().GetString("column-another")
			outputPath, _ := cmd.Flags().GetString("output")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runInclude(
				format,
				inputPath,
				targetColumnName,
				anotherPath,
				outputPath,
				IncludeOptions{
					anotherColumnName: anotherColumnName,
				})
		},
	}

	includeCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	includeCmd.MarkFlagRequired("input")
	includeCmd.Flags().StringP("column", "c", "", "Name of the column to use for filtering.")
	includeCmd.MarkFlagRequired("column")
	includeCmd.Flags().StringP("another", "a", "", "Another CSV file path. Filter by included in this CSV file.")
	includeCmd.MarkFlagRequired("another")
	includeCmd.Flags().StringP("column-another", "", "", "(optional) Name of the column to use for filtering in the another CSV file. Specify if different from the input CSV file.")
	includeCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	includeCmd.MarkFlagRequired("output")

	return includeCmd
}

type IncludeOptions struct {
	anotherColumnName string
}

func runInclude(format csv.Format, inputPath string, targetColumnName string, anotherPath string, outputPath string, options IncludeOptions) error {

	reader, writer, close, err := setupInputOutput(inputPath, outputPath, format)
	if err != nil {
		return err
	}
	defer close()

	anotherReader, anotherClose, err := setupInput(anotherPath, format)
	if err != nil {
		return err
	}
	defer anotherClose()

	err = include(reader, targetColumnName, anotherReader, writer, options)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func include(reader csv.CsvReader, targetColumnName string, anotherReader csv.CsvReader, writer csv.CsvWriter, options IncludeOptions) error {

	inputTargetColumnName := targetColumnName
	anotherTargetColumnName := targetColumnName
	if options.anotherColumnName != "" {
		anotherTargetColumnName = options.anotherColumnName
	}

	inputColumnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the input CSV file")
	}
	inputTargetColumnIndex := util.IndexOf(inputColumnNames, inputTargetColumnName)
	if inputTargetColumnIndex == -1 {
		return fmt.Errorf("missing %s in the input CSV file", inputTargetColumnName)
	}

	anotherItemSet, err := csv.LoadItemSet(anotherReader, anotherTargetColumnName)
	if err != nil {
		return errors.Wrap(err, "failed to read the another CSV file")
	}

	err = writer.Write(inputColumnNames)
	if err != nil {
		return err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the input CSV file")
		}

		// 比較対象のCSV内に存在した場合は出力
		if anotherItemSet.Contains(row[inputTargetColumnIndex]) {

			err = writer.Write(row)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
