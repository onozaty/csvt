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

var joinCmd = &cobra.Command{
	Use: "join",
	RunE: func(cmd *cobra.Command, args []string) error {

		firstPath, _ := cmd.Flags().GetString("first")
		secondPath, _ := cmd.Flags().GetString("second")
		joinColumnName, _ := cmd.Flags().GetString("column")
		outputPath, _ := cmd.Flags().GetString("output")
		useFileTable, _ := cmd.Flags().GetBool("usingfile")

		return runJoin(firstPath, secondPath, joinColumnName, outputPath, useFileTable)
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)

	joinCmd.Flags().StringP("first", "1", "", "First CSV file path")
	joinCmd.MarkFlagRequired("first")
	joinCmd.Flags().StringP("second", "2", "", "Second CSV file path")
	joinCmd.MarkFlagRequired("second")
	joinCmd.Flags().StringP("column", "c", "", "Name of the column to use for the join")
	joinCmd.MarkFlagRequired("column")
	joinCmd.Flags().StringP("output", "o", "", "Output CSV file path")
	joinCmd.MarkFlagRequired("output")
	joinCmd.Flags().BoolP("usingfile", "", false, "Use temporary files for join (Use this when joining large files that will not fit in memory)")
	joinCmd.Flags().SortFlags = false
}

func runJoin(firstPath string, secondPath string, joinColumnName string, outputPath string, useFileTable bool) error {

	firstFile, err := os.Open(firstPath)
	if err != nil {
		return err
	}
	defer firstFile.Close()

	firstReader, err := csv.NewCsvReader(firstFile)
	if err != nil {
		return err
	}

	secondFile, err := os.Open(secondPath)
	if err != nil {
		return err
	}
	defer secondFile.Close()

	secondReader, err := csv.NewCsvReader(secondFile)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	out := csv.NewCsvWriter(outputFile)

	err = join(firstReader, secondReader, joinColumnName, out, useFileTable)

	out.Flush()

	return err
}

func join(first csv.CsvReader, second csv.CsvReader, joinColumnName string, out csv.CsvWriter, useFileTable bool) error {

	var secondTable csv.CsvTable
	var err error

	if useFileTable {
		secondTable, err = csv.LoadCsvFileTable(second, joinColumnName)
	} else {
		secondTable, err = csv.LoadCsvMemoryTable(second, joinColumnName)
	}
	if err != nil {
		return errors.Wrap(err, "failed to read the second CSV file")
	}
	defer secondTable.Close()

	firstColumnNames, err := first.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the first CSV file")
	}
	firstJoinColumnIndex := util.IndexOf(firstColumnNames, joinColumnName)
	if firstJoinColumnIndex == -1 {
		return fmt.Errorf("missing %s in the first CSV file", joinColumnName)
	}

	// 追加するものは、結合用のカラムを除く
	appendsecondColumnNames := util.Remove(secondTable.ColumnNames(), joinColumnName)
	outColumnNames := append(firstColumnNames, appendsecondColumnNames...)
	out.Write(outColumnNames)

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

		secondRow := make([]string, len(appendsecondColumnNames))

		for i, appendColumnName := range appendsecondColumnNames {
			if secondRowMap != nil {
				secondRow[i] = secondRowMap[appendColumnName]
			}
		}

		err = out.Write(append(firstRow, secondRow...))
		if err != nil {
			return err
		}
	}

	return nil
}
