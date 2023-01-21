package cmd

import (
	"fmt"
	"io"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func newRenameCmd() *cobra.Command {

	renameCmd := &cobra.Command{
		Use:   "rename",
		Short: "Rename columns",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnNames, _ := cmd.Flags().GetStringArray("column")
			afterColumnNames, _ := cmd.Flags().GetStringArray("after")
			outputPath, _ := cmd.Flags().GetString("output")

			if len(targetColumnNames) != len(afterColumnNames) {
				return fmt.Errorf("the number of columns before and after the renaming is unmatched")
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runRename(
				format,
				inputPath,
				targetColumnNames,
				afterColumnNames,
				outputPath)
		},
	}

	renameCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	renameCmd.MarkFlagRequired("input")
	renameCmd.Flags().StringArrayP("column", "c", []string{}, "Name of column before renaming.")
	renameCmd.MarkFlagRequired("column")
	renameCmd.Flags().StringArrayP("after", "a", []string{}, "Name of column after renaming.")
	renameCmd.MarkFlagRequired("after")
	renameCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	renameCmd.MarkFlagRequired("output")

	return renameCmd
}

func runRename(format csv.Format, inputPath string, targetColumnNames []string, afterColumnNames []string, outputPath string) error {

	reader, writer, close, err := setupInputOutput(inputPath, outputPath, format)
	if err != nil {
		return err
	}
	defer close()

	err = rename(reader, targetColumnNames, afterColumnNames, writer)

	if err != nil {
		return err
	}

	return writer.Flush()
}

func rename(reader csv.CsvReader, targetColumnNames []string, afterColumnNames []string, writer csv.CsvWriter) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	// いったん対象カラムの位置を記憶して、変更後の名前に置き換え
	// (変更後の名前で探してしまわないように)
	targetColumnIndexes := []int{}
	for _, targetColumnName := range targetColumnNames {

		targetColumnIndex := slices.Index(columnNames, targetColumnName)
		if targetColumnIndex == -1 {
			return fmt.Errorf("missing %s in the CSV file", targetColumnName)
		}

		targetColumnIndexes = append(targetColumnIndexes, targetColumnIndex)
	}

	// 変更後に置き換え
	for i, targetColumnIndex := range targetColumnIndexes {
		columnNames[targetColumnIndex] = afterColumnNames[i]
	}

	err = writer.Write(columnNames)
	if err != nil {
		return err
	}

	// ヘッダ以外はそのまま書き込み
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the CSV file")
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
