package cmd

import (
	"fmt"
	"io"

	"github.com/onozaty/csvt/csv"
	"github.com/onozaty/csvt/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newJoinCmd() *cobra.Command {

	joinCmd := &cobra.Command{
		Use:   "join",
		Short: "Join CSV files",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			firstPath, _ := cmd.Flags().GetString("first")
			secondPath, _ := cmd.Flags().GetString("second")
			joinColumnName, _ := cmd.Flags().GetString("column")
			outputPath, _ := cmd.Flags().GetString("output")

			secondJoinColumnName, _ := cmd.Flags().GetString("column-second")
			useFileTable, _ := cmd.Flags().GetBool("usingfile")
			noRecordNoError, _ := cmd.Flags().GetBool("norecord")
			joinOptions := JoinOptions{
				secondJoinColumnName: secondJoinColumnName,
				useFileTable:         useFileTable,
				noRecordNoError:      noRecordNoError,
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runJoin(format, firstPath, secondPath, joinColumnName, outputPath, joinOptions)
		},
	}

	joinCmd.Flags().StringP("first", "1", "", "First CSV file path.")
	joinCmd.MarkFlagRequired("first")
	joinCmd.Flags().StringP("second", "2", "", "Second CSV file path.")
	joinCmd.MarkFlagRequired("second")
	joinCmd.Flags().StringP("column", "c", "", "Name of the column to use for joining.")
	joinCmd.MarkFlagRequired("column")
	joinCmd.Flags().StringP("column-second", "", "", "(optional) Name of the column to use for joining in the second CSV file. Specify if different from the first CSV file.")
	joinCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	joinCmd.MarkFlagRequired("output")
	joinCmd.Flags().BoolP("usingfile", "", false, "(optional) Use temporary files for joining. Use this when joining large files that will not fit in memory.")
	joinCmd.Flags().BoolP("norecord", "", false, "(optional) No error even if there is no record corresponding to sencod CSV.")

	return joinCmd
}

type JoinOptions struct {
	secondJoinColumnName string
	useFileTable         bool
	noRecordNoError      bool
}

func runJoin(format csv.Format, firstPath string, secondPath string, joinColumnName string, outputPath string, options JoinOptions) error {

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

	err = join(firstReader, secondReader, joinColumnName, writer, options)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func join(first csv.CsvReader, second csv.CsvReader, joinColumnName string, writer csv.CsvWriter, options JoinOptions) error {

	firstJoinColumnName := joinColumnName
	secondJoinColumnName := joinColumnName
	if options.secondJoinColumnName != "" {
		secondJoinColumnName = options.secondJoinColumnName
	}

	var secondTable csv.CsvTable
	var err error

	if options.useFileTable {
		secondTable, err = csv.LoadCsvFileTable(second, secondJoinColumnName)
	} else {
		secondTable, err = csv.LoadCsvMemoryTable(second, secondJoinColumnName)
	}
	if err != nil {
		return errors.Wrap(err, "failed to read the second CSV file")
	}
	defer secondTable.Close()

	firstColumnNames, err := first.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the first CSV file")
	}
	firstJoinColumnIndex := util.IndexOf(firstColumnNames, firstJoinColumnName)
	if firstJoinColumnIndex == -1 {
		return fmt.Errorf("missing %s in the first CSV file", firstJoinColumnName)
	}

	// 追加するものは、結合用のカラムを除く
	appendsecondColumnNames := util.Remove(secondTable.ColumnNames(), secondJoinColumnName)
	outColumnNames := append(firstColumnNames, appendsecondColumnNames...)
	err = writer.Write(outColumnNames)
	if err != nil {
		return err
	}

	// 基準となるCSVを読み込みながら、結合用のカラムの値をキーとしてもう片方のCSVから値を取得
	for {
		firstRow, err := first.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the first CSV file")
		}

		secondRowMap, err := secondTable.Find(firstRow[firstJoinColumnIndex])
		if err != nil {
			return errors.Wrap(err, "failed to find the second CSV file")
		}

		if secondRowMap == nil && !options.noRecordNoError {
			// 対応するレコードが無かった場合にエラーに
			return fmt.Errorf(
				"%s was not found in the second CSV file\nif you don't want to raise an error, use the 'norecord' option",
				firstRow[firstJoinColumnIndex])
		}

		secondRow := make([]string, len(appendsecondColumnNames))

		for i, appendColumnName := range appendsecondColumnNames {
			if secondRowMap != nil {
				secondRow[i] = secondRowMap[appendColumnName]
			}
		}

		err = writer.Write(append(firstRow, secondRow...))
		if err != nil {
			return err
		}
	}

	return nil
}
