package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"
	"golang.org/x/text/encoding"
)

func getFlagBaseCsvFormat(f *pflag.FlagSet) (csv.Format, error) {

	return getFlagCsvFormat(f, "delim", "quote", "sep", "allquote", "encoding", "bom")
}

func getFlagCsvFormat(f *pflag.FlagSet, delimName string, quoteName string, sepName string, allquoteName string, encodingName string, bomName string) (csv.Format, error) {

	format := csv.Format{}
	if v, err := getFlagRune(f, delimName); err != nil {
		return format, err
	} else {
		format.Delimiter = v
	}
	if v, err := getFlagRune(f, quoteName); err != nil {
		return format, err
	} else {
		format.Quote = v
	}
	if v, err := getFlagEscapedString(f, sepName); err != nil {
		return format, err
	} else {
		format.RecordSeparator = v
	}
	if v, err := f.GetBool(allquoteName); err != nil {
		return format, err
	} else {
		format.AllQuotes = v
	}
	if v, err := getFlagEncoding(f, encodingName); err != nil {
		return format, err
	} else {
		format.Encoding = v
	}
	if v, err := f.GetBool(bomName); err != nil {
		return format, err
	} else {
		format.WithBom = v
	}
	return format, nil
}

func getFlagEscapedString(f *pflag.FlagSet, name string) (string, error) {

	str, _ := f.GetString(name)

	if !strings.Contains(str, `\`) {
		return str, nil
	}

	// \nのように指定されているものを、スケープ文字として扱えるように
	unq, err := strconv.Unquote(`"` + str + `"`)
	if err != nil {
		return "", errors.Wrapf(err, "Could not parse value %s of flag %s", str, name)
	}

	return unq, nil
}

func getFlagRune(f *pflag.FlagSet, name string) (rune, error) {

	str, err := getFlagEscapedString(f, name)
	if err != nil {
		return 0, err
	}

	rs := []rune(str)

	if len(rs) == 0 {
		return 0, nil
	}

	if len(rs) != 1 {
		return 0, fmt.Errorf("flag %s should be specified with a single character", name)
	}

	return rs[0], nil
}

func getFlagEncoding(f *pflag.FlagSet, name string) (encoding.Encoding, error) {

	str, _ := f.GetString(name)

	if str == "" {
		return nil, nil
	}

	return csv.Encoding(str)
}

func getTargetColumnsIndexes(allColumnNames []string, targetColumnNames []string) ([]int, error) {

	if len(targetColumnNames) == 0 {
		// 指定カラムが無かった場合、全てのカラムが対象

		targetColumnIndexes := []int{}

		for i := range allColumnNames {
			targetColumnIndexes = append(targetColumnIndexes, i)
		}

		return targetColumnIndexes, nil

	} else {

		targetColumnIndexes := []int{}
		for _, targetColumnName := range targetColumnNames {

			targetColumnIndex, err := getTargetColumnIndex(allColumnNames, targetColumnName)
			if err != nil {
				return nil, err
			}

			targetColumnIndexes = append(targetColumnIndexes, targetColumnIndex)
		}

		return targetColumnIndexes, nil
	}
}

func getTargetColumnIndex(allColumnNames []string, targetColumnName string) (int, error) {

	targetColumnIndex := slices.Index(allColumnNames, targetColumnName)
	if targetColumnIndex == -1 {
		return -1, fmt.Errorf("missing %s in the CSV file", targetColumnName)
	}

	return targetColumnIndex, nil
}

func setupInput(inputPath string, format csv.Format) (csv.CsvReader, func(), error) {

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewCsvReader(inputFile, format)

	close := func() {
		inputFile.Close()
	}

	return reader, close, nil
}

func setupOutput(outputPath string, format csv.Format) (csv.CsvWriter, func(), error) {

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewCsvWriter(outputFile, format)

	close := func() {
		outputFile.Close()
	}

	return writer, close, nil
}

func setupInputOutput(inputPath string, outputPath string, format csv.Format) (csv.CsvReader, csv.CsvWriter, func(), error) {

	reader, inputClose, err := setupInput(inputPath, format)
	if err != nil {
		return nil, nil, nil, err
	}

	writer, outputClose, err := setupOutput(outputPath, format)
	if err != nil {
		inputClose()
		return nil, nil, nil, err
	}

	allClose := func() {
		inputClose()
		outputClose()
	}

	return reader, writer, allClose, nil
}
