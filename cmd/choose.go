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

func newChooseCmd() *cobra.Command {

	chooseCmd := &cobra.Command{
		Use:   "choose",
		Short: "Choose columns from CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			inputPath, _ := cmd.Flags().GetString("input")
			targetColumnNames, _ := cmd.Flags().GetStringArray("column")
			outputPath, _ := cmd.Flags().GetString("output")

			return runChoose(
				inputPath,
				targetColumnNames,
				outputPath)
		},
	}

	chooseCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	chooseCmd.MarkFlagRequired("input")
	chooseCmd.Flags().StringArrayP("column", "c", []string{}, "Name of the column to choose.")
	chooseCmd.MarkFlagRequired("columns")
	chooseCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	chooseCmd.MarkFlagRequired("output")
	chooseCmd.Flags().SortFlags = false

	return chooseCmd
}

func runChoose(inputPath string, targetColumnNames []string, outputPath string) error {

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

	err = choose(reader, targetColumnNames, writer)

	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}

func choose(reader csv.CsvReader, chooseColumnNames []string, writer csv.CsvWriter) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	chooseColumnIndexes := []int{}
	for _, chooseColumnName := range chooseColumnNames {

		chooseColumnIndex := util.IndexOf(columnNames, chooseColumnName)
		if chooseColumnIndex == -1 {
			return fmt.Errorf("missing %s in the CSV file", chooseColumnName)
		}

		chooseColumnIndexes = append(chooseColumnIndexes, chooseColumnIndex)
	}

	// 指定されたカラムのみに絞るフィルタを定義
	filter := func(row []string) []string {

		filtered := []string{}

		for i, item := range row {

			if util.Contains(chooseColumnIndexes, i) {
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