package cmd

import (
	"io"
	"regexp"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newReplaceCmd() *cobra.Command {

	replaceCmd := &cobra.Command{
		Use:   "replace",
		Short: "Replace values in CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnNames, _ := cmd.Flags().GetStringArray("column")
			regexValue, _ := cmd.Flags().GetString("regex")
			replacement, _ := cmd.Flags().GetString("replacement")
			outputPath, _ := cmd.Flags().GetString("output")

			regex, err := regexp.Compile(regexValue)
			if err != nil {
				return errors.WithMessage(err, "regular expression specified in --regex is invalid")
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runReplace(
				format,
				inputPath,
				targetColumnNames,
				regex,
				replacement,
				outputPath)
		},
	}

	replaceCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	replaceCmd.MarkFlagRequired("input")
	replaceCmd.Flags().StringArrayP("column", "c", []string{}, "(optional) Name of the column to replace. If not specified, all columns are targeted.")
	replaceCmd.Flags().StringP("regex", "r", "", "The regular expression to replace.")
	replaceCmd.MarkFlagRequired("regex")
	replaceCmd.Flags().StringP("replacement", "t", "", "The string after replace.")
	replaceCmd.MarkFlagRequired("replacement")
	replaceCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	replaceCmd.MarkFlagRequired("output")

	return replaceCmd
}

func runReplace(format csv.Format, inputPath string, targetColumnNames []string, regex *regexp.Regexp, replacement string, outputPath string) error {

	reader, writer, close, err := setupInputOutput(inputPath, outputPath, format)
	if err != nil {
		return err
	}
	defer close()

	err = replace(reader, targetColumnNames, regex, replacement, writer)

	if err != nil {
		return err
	}

	return writer.Flush()
}

func replace(reader csv.CsvReader, targetColumnNames []string, regex *regexp.Regexp, replacement string, writer csv.CsvWriter) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	targetColumnIndexes, err := getTargetColumnsIndexes(columnNames, targetColumnNames)
	if err != nil {
		return err
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

		// 置換
		for _, targetColumnIndex := range targetColumnIndexes {
			row[targetColumnIndex] = regex.ReplaceAllString(row[targetColumnIndex], replacement)
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
