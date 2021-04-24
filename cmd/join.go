package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use: "join",
	RunE: func(cmd *cobra.Command, args []string) error {

		firstPath, _ := cmd.Flags().GetString("first")
		secondPath, _ := cmd.Flags().GetString("second")
		joinColumnName, _ := cmd.Flags().GetString("output")
		outputPath, _ := cmd.Flags().GetString("column")

		return execute(firstPath, secondPath, joinColumnName, outputPath)
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
	joinCmd.Flags().SortFlags = false
}

type CsvTable interface {
	find(key string) map[string]string
	joinColumnName() string
	columnNames() []string
}

type MemoryTable struct {
	JoinColumnName string
	ColumnNames    []string
	Rows           map[string][]string
}

func execute(firstPath string, secondPath string, joinColumnName string, outputPath string) error {

	firstFile, err := os.Open(firstPath)
	if err != nil {
		return err
	}
	defer firstFile.Close()

	firstReader, err := newCsvReader(firstFile)
	if err != nil {
		return err
	}

	secondFile, err := os.Open(secondPath)
	if err != nil {
		return err
	}
	defer secondFile.Close()

	secondReader, err := newCsvReader(secondFile)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	out := csv.NewWriter(outputFile)

	err = join(firstReader, secondReader, joinColumnName, out)

	out.Flush()

	return err
}

func join(first *csv.Reader, second *csv.Reader, joinColumnName string, out *csv.Writer) error {

	secondTable, err := loadCsvTable(second, joinColumnName)
	if err != nil {
		return errors.Wrap(err, "Failed to read the second CSV file.")
	}

	firstColumnNames, err := first.Read()
	if err != nil {
		return errors.Wrap(err, "Failed to read the first CSV file.")
	}
	firstJoinColumnIndex := indexOf(firstColumnNames, joinColumnName)
	if firstJoinColumnIndex == -1 {
		return fmt.Errorf("%s is not found.", joinColumnName)
	}

	// 追加するものは、結合用のカラムを除く
	appendsecondColumnNames := remove(secondTable.columnNames(), joinColumnName)
	outColumnNames := append(firstColumnNames, appendsecondColumnNames...)
	out.Write(outColumnNames)

	// 基準となるCSVを読み込みながら、結合用のカラムの値をキーとしてもう片方のCSVから値を取得
	for {
		firstRow, err := first.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "Failed to read the first CSV file.")
		}

		secondRowMap := secondTable.find(firstRow[firstJoinColumnIndex])
		secondRow := make([]string, len(appendsecondColumnNames))

		for i, appendColumnName := range appendsecondColumnNames {
			if secondRowMap != nil {
				secondRow[i] = secondRowMap[appendColumnName]
			}
		}

		out.Write(append(firstRow, secondRow...))
	}

	return nil
}

var utf8bom = []byte{0xEF, 0xBB, 0xBF}

func newCsvReader(file *os.File) (*csv.Reader, error) {

	br := bufio.NewReader(file)
	mark, err := br.Peek(len(utf8bom))
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(mark, utf8bom) {
		// BOMがあれば読み飛ばす
		br.Discard(len(utf8bom))
	}

	return csv.NewReader(br), nil
}

func (t *MemoryTable) find(key string) map[string]string {

	row := t.Rows[key]

	if row == nil {
		return nil
	}

	rowMap := make(map[string]string)
	for i := 0; i < len(t.ColumnNames); i++ {
		rowMap[t.ColumnNames[i]] = row[i]
	}

	return rowMap
}

func (t *MemoryTable) joinColumnName() string {

	return t.JoinColumnName
}

func (t *MemoryTable) columnNames() []string {

	return t.ColumnNames
}

func loadCsvTable(reader *csv.Reader, joinColumnName string) (CsvTable, error) {

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	primaryColumnIndex := indexOf(headers, joinColumnName)
	if primaryColumnIndex == -1 {
		return nil, fmt.Errorf("%s is not found.", joinColumnName)
	}

	rows := make(map[string][]string)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// 格納前に既にあるか確認
		// -> 重複して存在した場合はエラーに
		_, has := rows[row[primaryColumnIndex]]
		if has {
			return nil, fmt.Errorf("%s is duplicated.", row[primaryColumnIndex])
		}

		rows[row[primaryColumnIndex]] = row
	}

	return &MemoryTable{
		JoinColumnName: joinColumnName,
		ColumnNames:    headers,
		Rows:           rows,
	}, nil
}

func indexOf(strings []string, search string) int {

	for i, v := range strings {
		if v == search {
			return i
		}
	}
	return -1
}

func remove(strings []string, search string) []string {

	result := []string{}
	for _, v := range strings {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
