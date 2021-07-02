package cmd

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newFilterCmd() *cobra.Command {

	filterCmd := &cobra.Command{
		Use:   "filter",
		Short: "Filter rows of CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnNames, _ := cmd.Flags().GetStringArray("column")
			outputPath, _ := cmd.Flags().GetString("output")
			equalValue, _ := cmd.Flags().GetString("equal")
			regexValue, _ := cmd.Flags().GetString("regex")

			if equalValue != "" && regexValue != "" {
				return fmt.Errorf("not allowed to specify both --equal and --regex")
			}

			var regex *regexp.Regexp = nil
			if regexValue != "" {
				regex, err = regexp.Compile(regexValue)
				if err != nil {
					return errors.WithMessage(err, "regular expression specified in --regex is invalid")
				}
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runFilter(
				format,
				inputPath,
				targetColumnNames,
				outputPath,
				FilterOptions{
					equalValue: equalValue,
					regex:      regex,
				})
		},
	}

	filterCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	filterCmd.MarkFlagRequired("input")
	filterCmd.Flags().StringArrayP("column", "c", []string{}, "(optional) Name of the column to use for filtering. If not specified, all columns are targeted.")
	filterCmd.Flags().StringP("equal", "", "", "(optional) Filter by matching value. If neither --equal nor --regex is specified, it will filter by those with values.")
	filterCmd.Flags().StringP("regex", "", "", "(optional) Filter by regular expression. If neither --equal nor --regex is specified, it will filter by those with values.")
	filterCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	filterCmd.MarkFlagRequired("output")

	return filterCmd
}

type FilterOptions struct {
	equalValue string
	regex      *regexp.Regexp
}

func runFilter(format csv.Format, inputPath string, targetColumnNames []string, outputPath string, options FilterOptions) error {

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

	err = filter(reader, targetColumnNames, writer, options)

	if err != nil {
		return err
	}

	return writer.Flush()
}

func filter(reader csv.CsvReader, targetColumnNames []string, writer csv.CsvWriter, options FilterOptions) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	targetColumnIndexes, err := getTargetColumnIndexes(columnNames, targetColumnNames)
	if err != nil {
		return err
	}

	// 行を絞るフィルタを定義
	filter := func(row []string) bool {

		// 対象のカラムを順次比較していく
		for _, targetColumnIndex := range targetColumnIndexes {
			value := row[targetColumnIndex]

			if options.equalValue != "" {
				if value == options.equalValue {
					return true
				}
			} else if options.regex != nil {
				if options.regex.MatchString(value) {
					return true
				}
			} else {
				if value != "" {
					return true
				}
			}
		}

		return false
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
