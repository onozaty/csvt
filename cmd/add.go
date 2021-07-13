package cmd

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/onozaty/csvt/csv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add column",
		RunE: func(cmd *cobra.Command, args []string) error {

			format, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			inputPath, _ := cmd.Flags().GetString("input")
			addColumnName, _ := cmd.Flags().GetString("column")
			staticValue, err := getFlagEscapedString(cmd.Flags(), "value") // バックスラッシュ記法を使いたい項目
			if err != nil {
				return err
			}
			templString, err := getFlagEscapedString(cmd.Flags(), "template") // バックスラッシュ記法を使いたい項目
			if err != nil {
				return err
			}
			copyColumnName, _ := cmd.Flags().GetString("copy-column")
			outputPath, _ := cmd.Flags().GetString("output")

			specifyValueCount := 0
			if staticValue != "" {
				specifyValueCount++
			}
			if templString != "" {
				specifyValueCount++
			}
			if copyColumnName != "" {
				specifyValueCount++
			}
			if specifyValueCount >= 2 {
				return fmt.Errorf("not allowed to specify both --value and --template and --copy-column")
			}

			var templ *template.Template = nil

			if templString != "" {
				templ, err = template.New("template").Parse(templString)
				if err != nil {
					return errors.WithMessage(err, "--template is invalid")
				}
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runAdd(
				format,
				inputPath,
				addColumnName,
				outputPath,
				AddOptions{
					staticValue:    staticValue,
					template:       templ,
					copyColumnName: copyColumnName,
				})
		},
	}

	addCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	addCmd.MarkFlagRequired("input")
	addCmd.Flags().StringP("column", "c", "", "Name of the column to add.")
	addCmd.MarkFlagRequired("column")
	addCmd.Flags().StringP("value", "", "", "(optional) Fixed value to set for the added column.")
	addCmd.Flags().StringP("template", "", "", "(optional) Template for the value to be set for the added column.")
	addCmd.Flags().StringP("copy-column", "", "", "(optional) Name of the column from which the value is copied.")
	addCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	addCmd.MarkFlagRequired("output")

	return addCmd
}

type AddOptions struct {
	staticValue    string
	template       *template.Template
	copyColumnName string
}

func runAdd(format csv.Format, inputPath string, addColumnName string, outputPath string, options AddOptions) error {

	reader, writer, close, err := setupInputOutput(inputPath, outputPath, format)
	if err != nil {
		return err
	}
	defer close()

	err = add(reader, addColumnName, writer, options)

	if err != nil {
		return err
	}

	return writer.Flush()
}

func add(reader csv.CsvReader, addColumnName string, writer csv.CsvWriter, options AddOptions) error {

	// ヘッダ
	columnNames, err := reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read the CSV file")
	}

	copyColumnIndex := -1
	if options.copyColumnName != "" {
		copyColumnIndex, err = getTargetColumnIndex(columnNames, options.copyColumnName)
		if err != nil {
			return err
		}
	}

	err = writer.Write(append(columnNames, addColumnName))
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

		addColumnValue := ""
		if options.staticValue != "" {

			addColumnValue = options.staticValue

		} else if options.template != nil {

			addColumnValue, err = executeTemplate(options.template, columnNames, row)
			if err != nil {
				return err
			}

		} else if options.copyColumnName != "" {

			addColumnValue = row[copyColumnIndex]
		}

		err = writer.Write(append(row, addColumnValue))
		if err != nil {
			return err
		}
	}

	return nil
}

func executeTemplate(templ *template.Template, columnNames []string, row []string) (string, error) {

	// テンプレートにはヘッダ名をキーとしたmapで
	rowMap := make(map[string]string)
	for i, columnName := range columnNames {
		rowMap[columnName] = row[i]
	}

	w := new(strings.Builder)
	err := templ.Execute(w, rowMap)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}
