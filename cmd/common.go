package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

func getFlagBaseCsvFormat(f *pflag.FlagSet) (csv.Format, error) {

	return getFlagCsvFormat(f, "delim", "quote", "sep", "allquote")
}

func getFlagCsvFormat(f *pflag.FlagSet, delimName string, quoteName string, sepName string, allquoteName string) (csv.Format, error) {

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
