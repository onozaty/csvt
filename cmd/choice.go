package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/onozaty/csvt/csv"
	"github.com/onozaty/csvt/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newChoiceCmd() *cobra.Command {

	choiceCmd := &cobra.Command{
		Use:   "choice",
		Short: "Choice columns from CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnNames, _ := cmd.Flags().GetStringArray("column")
			outputPath, _ := cmd.Flags().GetString("output")

			return runChoice(
				inputPath,
				targetColumnNames,
				outputPath)
		},
	}

	choiceCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	choiceCmd.MarkFlagRequired("input")
	choiceCmd.Flags().StringArrayP("column", "c", []string{}, "Name of the column to choice.")
	choiceCmd.MarkFlagRequired("columns")
	choiceCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	choiceCmd.MarkFlagRequired("output")
	choiceCmd.Flags().SortFlags = false

	return choiceCmd
}

func runChoice(inputPath string, targetColumnNames []string, outputPath string) error {

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

	err = choice(reader, targetColumnNames, writer)

	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}

func choice(reader csv.CsvReader, choiceColumnNames []string, writer csv.CsvWriter) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	choiceColumnIndexes := []int{}
	for _, choiceColumnName := range choiceColumnNames {

		choiceColumnIndex := util.IndexOf(columnNames, choiceColumnName)
		if choiceColumnIndex == -1 {
			return fmt.Errorf("missing %s in the CSV file", choiceColumnName)
		}

		choiceColumnIndexes = append(choiceColumnIndexes, choiceColumnIndex)
	}

	// 指定されたカラムのみに絞るフィルタを定義
	filter := func(row []string) []string {

		filtered := []string{}

		for i, item := range row {

			if util.Contains(choiceColumnIndexes, i) {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	err = writer.Write(filter(columnNames))
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

		err = writer.Write(filter(row))
		if err != nil {
			return err
		}
	}

	return nil
}
