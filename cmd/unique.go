package cmd

import (
	"io"
	"strings"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newUniqueCmd() *cobra.Command {

	uniqueCmd := &cobra.Command{
		Use:   "unique",
		Short: "Extract unique rows",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnNames, _ := cmd.Flags().GetStringArray("column")
			outputPath, _ := cmd.Flags().GetString("output")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runUnique(
				format,
				inputPath,
				targetColumnNames,
				outputPath)
		},
	}

	uniqueCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	uniqueCmd.MarkFlagRequired("input")
	uniqueCmd.Flags().StringArrayP("column", "c", []string{}, "Name of the column to use for extract unique rows.")
	uniqueCmd.MarkFlagRequired("column")
	uniqueCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	uniqueCmd.MarkFlagRequired("output")

	return uniqueCmd
}

func runUnique(format csv.Format, inputPath string, targetColumnNames []string, outputPath string) error {

	reader, writer, close, err := setupInputOutput(inputPath, outputPath, format)
	if err != nil {
		return err
	}
	defer close()

	err = unique(reader, targetColumnNames, writer)

	if err != nil {
		return err
	}

	return writer.Flush()
}

func unique(reader csv.CsvReader, targetColumnNames []string, writer csv.CsvWriter) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	targetColumnIndexes, err := getTargetColumnsIndexes(columnNames, targetColumnNames)
	if err != nil {
		return err
	}

	// キー作成用の関数
	const keyConcatChar = "\x00"
	makeKey := func(row []string) string {
		key := ""
		for i, targetColumnIndex := range targetColumnIndexes {

			if i != 0 {
				// 複数カラムの場合、結合して
				key += keyConcatChar
			}

			// 区切り文字と項目内の値を区別するためエスケープ
			key += strings.ReplaceAll(row[targetColumnIndex], keyConcatChar, keyConcatChar+keyConcatChar)
		}

		return key
	}

	err = writer.Write(columnNames)
	if err != nil {
		return err
	}

	// 重複チェック用のmap(valueは利用しないので一律0を入れる)
	keyMap := make(map[string]int)

	// ヘッダ以外
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the CSV file")
		}

		key := makeKey(row)

		_, has := keyMap[key]
		if !has {
			// 重複していない行なので書き込み
			err = writer.Write(row)
			if err != nil {
				return err
			}

			keyMap[key] = 0
		}
	}

	return nil
}
