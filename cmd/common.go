package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/onozaty/csvt/csv"
	"github.com/onozaty/csvt/util"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
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
	if v, err := getFlagString(f, sepName); err != nil {
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

func getFlagString(f *pflag.FlagSet, name string) (string, error) {

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

	str, err := getFlagString(f, name)
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

func getTargetColumnIndexes(allColumnNames []string, targetColumnNames []string) ([]int, error) {

	if len(targetColumnNames) == 0 {
		// 指定カラムが無かった場合、全てのカラムが対象

		targetColumnIndexes := []int{}

		for i, _ := range allColumnNames {
			targetColumnIndexes = append(targetColumnIndexes, i)
		}

		return targetColumnIndexes, nil

	} else {

		targetColumnIndexes := []int{}
		for _, targetColumnName := range targetColumnNames {

			targetColumnIndex := util.IndexOf(allColumnNames, targetColumnName)
			if targetColumnIndex == -1 {
				return nil, fmt.Errorf("missing %s in the CSV file", targetColumnName)
			}

			targetColumnIndexes = append(targetColumnIndexes, targetColumnIndex)
		}

		return targetColumnIndexes, nil
	}
}
