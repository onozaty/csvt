package cmd

import (
	"io"
	"os"

	"github.com/onozaty/csvt/csv"
	"github.com/spf13/cobra"
)

func newTransformCmd() *cobra.Command {

	transformCmd := &cobra.Command{
		Use:   "transform",
		Short: "Transform the format of CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {

			inputPath, _ := cmd.Flags().GetString("input")
			inputFormat, err := getFlagBaseCsvFormat(cmd.Flags())
			if err != nil {
				return err
			}

			outputPath, _ := cmd.Flags().GetString("output")
			outputFormat, err := getFlagCsvFormat(cmd.Flags(), "out-delim", "out-quote", "out-sep", "out-allquote", "out-encoding", "out-bom")
			if err != nil {
				return err
			}

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runTransform(
				inputPath,
				inputFormat,
				outputPath,
				outputFormat)
		},
	}

	transformCmd.Flags().StringP("input", "i", "", "Input CSV file path.")
	transformCmd.MarkFlagRequired("input")
	transformCmd.Flags().StringP("output", "o", "", "Output CSV file path.")
	transformCmd.MarkFlagRequired("output")
	transformCmd.Flags().StringP("out-delim", "", "", "(optional) Output CSV delimiter. The default is ','")
	transformCmd.Flags().StringP("out-quote", "", "", "(optional) Output CSV quote. The default is '\"'")
	transformCmd.Flags().StringP("out-sep", "", "", "(optional) Output CSV record separator. The default is CRLF.")
	transformCmd.Flags().BoolP("out-allquote", "", false, "(optional) Always quote output CSV fields. The default is to quote only the necessary fields.")
	transformCmd.Flags().StringP("out-encoding", "", "", "(optional) Output CSV encoding. The default is utf-8. Supported encodings: utf-8, shift_jis, euc-jp")
	transformCmd.Flags().BoolP("out-bom", "", false, "(optional) Output CSV with BOM.")

	return transformCmd
}

func runTransform(inputPath string, inputFormat csv.Format, outputPath string, outputFormat csv.Format) error {

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	reader := csv.NewCsvReader(inputFile, inputFormat)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	writer := csv.NewCsvWriter(outputFile, outputFormat)

	// フォーマットは設定済みなので、そのままコピーするだけ
	if err := copy(reader, writer); err != nil {
		return err
	}

	return writer.Flush()
}

func copy(reader csv.CsvReader, writer csv.CsvWriter) error {

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
