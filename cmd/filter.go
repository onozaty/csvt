package cmd

import (
	"fmt"
	"io"
	"regexp"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newFilterCmd() *cobra.Command {

	filterCmd := &cobra.Command{
		Use:   "filter",
		Short: "Filter rows by condition",
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
			nonMatch, _ := cmd.Flags().GetBool("not")
			equalColumnName, _ := cmd.Flags().GetString("equal-column")

			optionCondCount := 0
			if equalValue != "" {
				optionCondCount++
			}
			if regexValue != "" {
				optionCondCount++
			}
			if equalColumnName != "" {
				optionCondCount++
			}
			if optionCondCount >= 2 {
				return fmt.Errorf("not allowed to specify both --equal and --regex and --equal-column")
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
					equalValue:      equalValue,
					regex:           regex,
					equalColumnName: equalColumnName,
					nonMatch:        nonMatch,
				})
		},
	}

	filterCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	filterCmd.MarkFlagRequired("input")
	filterCmd.Flags().StringArrayP("column", "c", []string{}, "(optional) Name of the column to use for filtering. If not specified, all columns are targeted.")
	filterCmd.Flags().StringP("equal", "", "", "(optional) Filter by matching value. If neither --equal nor --regex nor --equal-column is specified, it will filter by those with values.")
	filterCmd.Flags().StringP("regex", "", "", "(optional) Filter by regular expression.")
	filterCmd.Flags().StringP("equal-column", "", "", "(optional) Filter by other column value.")
	filterCmd.Flags().BoolP("not", "", false, "(optional) Filter by non-matches.")
	filterCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	filterCmd.MarkFlagRequired("output")

	return filterCmd
}

type FilterOptions struct {
	equalValue      string
	regex           *regexp.Regexp
	equalColumnName string
	nonMatch        bool
}

func runFilter(format csv.Format, inputPath string, targetColumnNames []string, outputPath string, options FilterOptions) error {

	reader, writer, close, err := setupInputOutput(inputPath, outputPath, format)
	if err != nil {
		return err
	}
	defer close()

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

	targetColumnIndexes, err := getTargetColumnsIndexes(columnNames, targetColumnNames)
	if err != nil {
		return err
	}

	equalColumnIndex := -1
	if options.equalColumnName != "" {
		equalColumnIndex, err = getTargetColumnIndex(columnNames, options.equalColumnName)
		if err != nil {
			return err
		}
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
			} else if options.equalColumnName != "" {
				if value == row[equalColumnIndex] {
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

		filterd := filter(row)
		if options.nonMatch {
			// 一致しなかったもので絞る場合、反転させる
			filterd = !filterd
		}

		if !filterd {
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
