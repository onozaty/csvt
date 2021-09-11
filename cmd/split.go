package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newSplitCmd() *cobra.Command {

	splitCmd := &cobra.Command{
		Use:   "split",
		Short: "Split into multiple CSV files",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			maxRows, _ := cmd.Flags().GetInt("rows")
			outputBasePath, _ := cmd.Flags().GetString("output")

			// 最大行数は1以上
			if maxRows <= 0 {
				return fmt.Errorf("rows must be greater than or equal to 1")
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runSplit(
				format,
				inputPath,
				maxRows,
				outputBasePath)
		},
	}

	splitCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	splitCmd.MarkFlagRequired("input")
	splitCmd.Flags().IntP("rows", "r", 0, "Maximum number of rows.")
	splitCmd.MarkFlagRequired("rows")
	splitCmd.Flags().StringP("output", "o", "", "Output CSV file base path. If you specify \"output.csv\", the file will be output as \"output-1.csv\" \"output-2.csv\" ...")
	splitCmd.MarkFlagRequired("output")

	return splitCmd
}

func runSplit(format csv.Format, inputPath string, maxRows int, outputBasePath string) error {

	reader, close, err := setupInput(inputPath, format)
	if err != nil {
		return err
	}
	defer close()

	// EOFを取り出せるReaderでラップ
	splitReader := &splitReader{r: reader}

	columnNames, err := splitReader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the input CSV file")
	}

	num := 1

	for {
		outputPath, err := makeOutputPath(outputBasePath, num)
		if err != nil {
			return err
		}

		err = splitOne(format, splitReader, maxRows, columnNames, outputPath)
		if err != nil {
			return err
		}

		if splitReader.EOF() {
			break
		}

		num++
	}

	return nil
}

func makeOutputPath(outputBasePath string, num int) (string, error) {

	// outputBathPath: dir/output.csv, num: 1
	//   -> dir/output-1.csv
	// outputBathPath: dir/output, num: 1
	//   -> dir/output-1

	ext := filepath.Ext(outputBasePath)
	outputPath := outputBasePath[:len(outputBasePath)-len(ext)] + "-" + strconv.Itoa(num) + ext

	dir := filepath.Dir(outputPath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {

		err = os.Mkdir(dir, 0755)
		if err != nil {
			return "", err
		}
	}

	return outputPath, nil
}

func splitOne(format csv.Format, reader csv.CsvReader, maxRows int, columnNames []string, outputPath string) error {

	writer, close, err := setupOutput(outputPath, format)
	if err != nil {
		return err
	}
	defer close()

	err = writer.Write(columnNames)
	if err != nil {
		return err
	}

	for i := 0; i < maxRows; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read the input CSV file")
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

type splitReader struct {
	r         csv.CsvReader
	isEof     bool
	peekedErr error
	peekedRow []string
}

func (r *splitReader) EOF() bool {

	if !r.isEof && r.peekedRow == nil {
		// 先読みしてEOFか確認
		row, err := r.r.Read()
		if err == io.EOF {
			r.isEof = true
		} else {
			r.peekedErr = err
			r.peekedRow = row
		}
	}

	return r.isEof
}

func (r *splitReader) Read() ([]string, error) {
	if r.peekedRow != nil || r.peekedErr != nil {
		// 先読みしている情報があれば、そちらを返す
		tempRow := r.peekedRow
		tempErr := r.peekedErr
		r.peekedRow = nil
		r.peekedErr = nil

		return tempRow, tempErr
	}

	row, err := r.r.Read()
	if err == io.EOF {
		r.isEof = true
	}

	return row, err
}
