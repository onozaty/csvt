package cmd

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/onozaty/csvt/csv"
	"github.com/onozaty/csvt/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newFilterCmd() *cobra.Command {

	filterCmd := &cobra.Command{
		Use:   "filter",
		Short: "Filter rows of CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnName, _ := cmd.Flags().GetString("column")
			outputPath, _ := cmd.Flags().GetString("output")
			equalValue, _ := cmd.Flags().GetString("equal")
			regexValue, _ := cmd.Flags().GetString("regex")

			if equalValue != "" && regexValue != "" {
				return fmt.Errorf("not allowed to specify both --equal and --regex")
			}

			var regex *regexp.Regexp = nil
			var err error = nil
			if regexValue != "" {
				regex, err = regexp.Compile(regexValue)
				if err != nil {
					return errors.WithMessage(err, "regular expression specified in --regex is invalid")
				}
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runFilter(
				inputPath,
				targetColumnName,
				outputPath,
				FilterOptions{
					equalValue: equalValue,
					regex:      regex,
				})
		},
	}

	filterCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	filterCmd.MarkFlagRequired("input")
	filterCmd.Flags().StringP("column", "c", "", "Name of the column to use for filtering. If neither --equal nor --regex is specified, it will filter by those with values.")
	filterCmd.MarkFlagRequired("column")
	filterCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	filterCmd.MarkFlagRequired("output")
	filterCmd.Flags().StringP("equal", "", "", "(optional) Filter by matching value.")
	filterCmd.Flags().StringP("regex", "", "", "(optional) Filter by regular expression.")
	filterCmd.Flags().SortFlags = false

	return filterCmd
}

type FilterOptions struct {
	equalValue string
	regex      *regexp.Regexp
}

func runFilter(inputPath string, targetColumnName string, outputPath string, options FilterOptions) error {

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	reader, err := csv.NewCsvReader(inputFile)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	writer := csv.NewCsvWriter(outputFile)

	err = filter(reader, targetColumnName, writer, options)

	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}

func filter(reader csv.CsvReader, targetColumnName string, writer csv.CsvWriter, options FilterOptions) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	targetColumnIndex := util.IndexOf(columnNames, targetColumnName)
	if targetColumnIndex == -1 {
		return fmt.Errorf("missing %s in the CSV file", targetColumnName)
	}

	// 行を絞るフィルタを定義
	filter := func(row []string) bool {

		value := row[targetColumnIndex]

		if options.equalValue != "" {
			return value == options.equalValue
		}

		if options.regex != nil {
			return options.regex.MatchString(value)
		}

		return value != ""
	}

	err = writer.Write(columnNames)
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

		if !filter(row) {
			// 一致しないものは出力しない
			continue
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
