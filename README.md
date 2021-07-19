# csvt

[![GitHub license](https://img.shields.io/github/license/onozaty/csvt)](https://github.com/onozaty/csvt/blob/main/LICENSE)
[![Test](https://github.com/onozaty/csvt/actions/workflows/test.yaml/badge.svg)](https://github.com/onozaty/csvt/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/onozaty/csvt/branch/main/graph/badge.svg?token=VSU64LAK8P)](https://codecov.io/gh/onozaty/csvt)

`csvt` is a command line tool for processing CSV.

`csvt` consists of multiple subcommands.

* [add](#add) Add column.
* [choose](#choose) Choose columns.
* [concat](#concat) Concat CSV files.
* [count](#count) Count the number of records.
* [exclude](#exclude) Exclude rows by included in another CSV file.
* [filter](#filter) Filter rows by condition.
* [header](#header) Show header.
* [include](#include) Filter rows by included in another CSV file.
* [join](#join) Join CSV files.
* [remove](#remove) Remove columns.
* [rename](#rename) Rename columns.
* [replace](#replace) Replace values.
* [sort](#sort) Sort rows.
* [slice](#slice) Slice specified range of rows.
* [transform](#transform) Transform format.
* [unique](#unique) Extract unique rows.

## Common flags

Flags related to the CSV format are available in each subcommand as common flags.

```
Global Flags:
      --delim string      (optional) CSV delimiter. The default is ','
      --quote string      (optional) CSV quote. The default is '"'
      --sep string        (optional) CSV record separator. The default is CRLF.
      --allquote          (optional) Always quote CSV fields. The default is to quote only the necessary fields.
      --encoding string   (optional) CSV encoding. The default is utf-8. Supported encodings: utf-8, shift_jis, euc-jp
      --bom               (optional) CSV with BOM. When reading, the BOM will be automatically removed without this flag.
```

For example, when dealing with TSV files, change the delimiter to a tab as shown below.

```
$ csvt count -i INPUT --delim "\t"
```

## add

Create a new CSV file by adding column to the input CSV file.

The following values can be set for the new column.

* Fixed value.
* Same value as another column.
* Value by template. As a template engine, [text/template](https://pkg.go.dev/text/template) will be used.

### Usage

```
csvt add -i INPUT -c ADD_COLUMN [--value VALUE | --template TEMPLATE | --copy-column FROM_COLUMN] -o OUTPUT
```

```
Usage:
  csvt add [flags]

Flags:
  -i, --input string         Input CSV file path.
  -c, --column string        Name of the column to add.
      --value string         (optional) Fixed value to set for the added column.
      --template string      (optional) Template for the value to be set for the added column.
      --copy-column string   (optional) Name of the column from which the value is copied.
  -o, --output string        Output CSV file path.
  -h, --help                 help for add
```

### Example

The contents of `input.csv`.

```
col1,col2
1,a
2,b
3,c
```

Add "col3" as a new column. Set "x" as a fixed value.

```
$ csvt add -i input.csv -c col3 --value x -o output.csv
```

The contents of the created `output.csv`.

```
col1,col2,col3
1,a,x
2,b,x
3,c,x
```

Add "col1x" by copying "col1".

```
$ csvt add -i input.csv -c col1x --copy-column col1 -o output.csv
```

```
col1,col2,col1x
1,a,1
2,b,2
3,c,3
```

Use the template to add a column that combines the values of "col1" and "col2".

```
$ csvt add -i input.csv -c col3 --template "{{.col1}}-{{.col2}}" -o output.csv
```

```
col1,col2,col3
1,a,1-a
2,b,2-b
3,c,3-c
```

Please refer to the following for template syntax.

* https://pkg.go.dev/text/template

## choose

Create a new CSV file by choosing columns from the input CSV file.

### Usage

```
csvt choose -i INPUT -c COLUMN1 ... -o OUTPUT
```

```
Usage:
  csvt choose [flags]

Flags:
  -i, --input string         Input CSV file path.
  -c, --column stringArray   Name of the column to choose.
  -o, --output string        Output CSV file path.
  -h, --help                 help for choose
```

### Example

The contents of `input.csv`.

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,2
2,Hanako,21,1
3,Smith,30,2
4,Jun,22,4
```

Create `output.csv` by choosing "Name" and "Age" from `input.csv`.

```
$ csvt choose -i input.csv -c Name -c Age -o output.csv
```

The contents of the created `output.csv`.

```
Name,Age
"Taro, Yamada",10
Hanako,21
Smith,30
Jun,22
```

## concat

Create a new CSV file by concat the two CSV files.  
Check the column names and concat them into the same column.

### Usage

```
csvt concat -1 INPUT1 -2 INPUT2 -o OUTPUT
```

```
Usage:
  csvt concat [flags]

Flags:
  -1, --first string    First CSV file path.
  -2, --second string   Second CSV file path.
  -o, --output string   Output CSV file path.
  -h, --help            help for concat
```

### Example

The contents of `input1.csv`.

```
ID,Name
1,name1
2,name2
```

The contents of `input2.csv`.

```
Name,ID
name3,3
name4,4
```

Concat `input1.csv` and `input2.csv`.

```
$ csvt concat -1 input1.csv -2 input2.csv -o output.csv
```

The contents of the created `output.csv`.

```
ID,Name
1,name1
2,name2
3,name3
4,name4
```

## count

Count the number of records in CSV file.

### Usage

```
csvt count -i INPUT [-c COLUMN] [--header]
```

```
Usage:
  csvt count [flags]

Flags:
  -i, --input string    CSV file path.
  -c, --column string   (optional) Name of the column to be counted. Only those with values will be counted.
      --header          (optional) Counting including header. The default is to exclude header.
  -h, --help            help for count
```

### Example

The contents of `input.csv`.

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,2
2,Hanako,21,1
3,Smith,30,
4,Jun,22,4
```

Count the number of records.

```
$ csvt count -i input.csv
4
```

Count the number of lines, including headers.

```
$ csvt count -i input.csv --header
5
```

Counts the number of records for which a value exists in "CompanyID".

```
$ csvt count -i input.csv -c CompanyID
3
```

## exclude

Create a new CSV file by exclude on the rows included in another CSV file.

### Usage

```
csvt exclude -i INPUT -c COLUMN -a ANOTHER [--column-another COLUMN2] -o OUTPUT
```

```
Usage:
  csvt exclude [flags]

Flags:
  -i, --input string            Input CSV file path.
  -c, --column string           Name of the column to use for exclude.
  -a, --another string          Another CSV file path. Exclude by included in this CSV file.
      --column-another string   (optional) Name of the column to use for exclude in the another CSV file. Specify if different from the input CSV file.
  -o, --output string           Output CSV file path.
  -h, --help                    help for exclude
```

### Example

The contents of `input.csv`.

```
col1,col2
1,A
2,B
3,C
4,D
```

The contents of `another.csv`.

```
col1,col3
2,2
3,2
```

Exclude by "col1" values in `another.csv`.

```
$ csvt exclude -i input.csv -c col1 -a another.csv -o output.csv
```

The contents of the created `output.csv`.

```
col1,col2
1,A
4,D
```

## filter

Create a new CSV file by filtering the input CSV file to rows that match the conditions.

### Usage

```
csvt filter -i INPUT [[-c COLUMN1] ...] [--equal VALUE | --regex REGEX | --equal-column COLUMN] [--not] -o OUTPUT
```

```
Usage:
  csvt filter [flags]

Flags:
  -i, --input string          Input CSV file path.
  -c, --column stringArray    (optional) Name of the column to use for filtering. If not specified, all columns are targeted.
      --equal string          (optional) Filter by matching value. If neither --equal nor --regex nor --equal-column is specified, it will filter by those with values.
      --regex string          (optional) Filter by regular expression.
      --equal-column string   (optional) Filter by other column value.
      --not                   (optional) Filter by non-matches.
  -o, --output string         Output CSV file path.
  -h, --help                  help for filter
```

### Example

The contents of `input.csv`.

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,1
2,Hanako,21,1
3,yamada,30,
4,Jun,22,2
```

Create `output.csv` by filter by non-empty values of "CompanyID" from `input.csv`.

```
$ csvt filter -i input.csv -c CompanyID -o output.csv
```

The contents of the created `output.csv`.

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,1
2,Hanako,21,1
4,Jun,22,2
```

You can also filter by matching the specified value.  
Specify a value by using `--equal`.

```
$ csvt filter -i input.csv -c CompanyID --equal 2 -o output.csv 
```

```
UserID,Name,Age,CompanyID
4,Jun,22,2
```

You can use `--not` to invert the filtering target.

```
$ csvt filter -i input.csv -c CompanyID --equal 2 --not -o output.csv 
```

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,1
2,Hanako,21,1
3,yamada,30,
```

You can also filter by matching with other column.
The column can be specified with `--equal-column`.

```
$ csvt filter -i input.csv -c UserID --equal-column CompanyID -o output.csv 
```

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,1
```

Regular expressions can also be used.  
Use `--regex` to specify a regular expression.

```
$ csvt filter -i input.csv -c Name --regex [Yy]amada -o output.csv 
```

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,1
3,yamada,30,
```

Please refer to the following for the syntax of regular expressions.

* https://pkg.go.dev/regexp/syntax

## header

Show the header of CSV file.

### Usage

```
csvt header -i INPUT
```

```
Usage:
  csvt header [flags]

Flags:
  -i, --input string   CSV file path.
  -h, --help           help for header
```

### Example

The contents of `input.csv`.

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,2
2,Hanako,21,1
3,Smith,30,2
4,Jun,22,4
```

```
$ csvt header -i input.csv
UserID
Name
Age
CompanyID
```

## include

Create a new CSV file by filtering on the rows included in another CSV file.

### Usage

```
csvt include -i INPUT -c COLUMN -a ANOTHER [--column-another COLUMN2] -o OUTPUT
```

```
Usage:
  csvt include [flags]

Flags:
  -i, --input string            Input CSV file path.
  -c, --column string           Name of the column to use for filtering.
  -a, --another string          Another CSV file path. Filter by included in this CSV file.
      --column-another string   (optional) Name of the column to use for filtering in the another CSV file. Specify if different from the input CSV file.
  -o, --output string           Output CSV file path.
  -h, --help                    help for include
```

### Example

The contents of `input.csv`.

```
col1,col2
1,A
2,B
3,C
4,D
```

The contents of `another.csv`.

```
col1,col3
2,2
3,2
```

Filter by "col1" values in `another.csv`.

```
$ csvt include -i input.csv -c col1 -a another.csv -o output.csv
```

The contents of the created `output.csv`.

```
col1,col2
2,B
3,C
```

## join

Join CSV files.  

Using the first CSV file as a base, join the contents of the second CSV file to create a new CSV file.  
It is similar to EXCEL's VLOOKUP.

### Usage

```
csvt join -1 INPUT1 -2 INPUT2 -c COLUMN [--column2 COLUMN2] -o OUTPUT [--usingfile] [--norecord]
```

```
Usage:
  csvt join [flags]

Flags:
  -1, --first string           First CSV file path.
  -2, --second string          Second CSV file path.
  -c, --column string          Name of the column to use for joining.
      --column-second string   (optional) Name of the column to use for joining in the second CSV file. Specify if different from the first CSV file.
  -o, --output string          Output CSV file path.
      --usingfile              (optional) Use temporary files for joining. Use this when joining large files that will not fit in memory.
      --norecord               (optional) No error even if there is no record corresponding to sencod CSV.
  -h, --help                   help for join
```

### Example

The contents of `input1.csv`.

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,2
2,Hanako,21,1
3,Smith,30,2
4,Jun,22,4
```

The contents of `input2.csv`.

```
CompanyID,CompanyName
1,CompanyA
2,CompanyB
3,CompanyC
4,"AAA Inc"
```

Join by "CompanyID".

```
$ csvt join -1 input1.csv -2 input2.csv -c CompanyID -o output.csv
```

The contents of the created `output.csv`.

```
UserID,Name,Age,CompanyID,CompanyName
1,"Taro, Yamada",10,2,CompanyB
2,Hanako,21,1,CompanyA
3,Smith,30,2,CompanyB
4,Jun,22,4,AAA Inc
```

If the `input2.csv` looks like the following and there is no corresponding "CompanyID", an error will occur.  

```
CompanyID,CompanyName
1,CompanyA
2,CompanyB
```

If you don't want to raise an error even if there is no value, specify `--norecord`

```
$ csvt join -1 input1.csv -2 input2.csv -c CompanyID -o output.csv --norecord
```

If the column name in the second CSV file is different from that in the first CSV file, specify it with `--column-second`.

```
$ csvt join -1 input1.csv -2 input2.csv -c CompanyID --column2 ID -o output.csv
```

If the second CSV file you specify is so large that it would take up too much memory on your PC, specify `--usingfile`.  
If you specify --usingfile, it will use a temporary file for joining instead of memory.

```
$ csvt join -1 input1.csv -2 input2.csv -c CompanyID -o output.csv --usingfile
```

## remove

Create a new CSV file by remove columns from the input CSV file.

### Usage

```
csvt remove -i INPUT -c COLUMN1 ... -o OUTPUT
```

```
Usage:
  csvt remove [flags]

Flags:
  -i, --input string         Input CSV file path.
  -c, --column stringArray   Name of the column to remove.
  -o, --output string        Output CSV file path.
  -h, --help                 help for remove
```

### Example

The contents of `input.csv`.

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,2
2,Hanako,21,1
3,Smith,30,2
4,Jun,22,4
```

Create `output.csv` by remove "Name" and "Age" from `input.csv`.

```
$ csvt remove -i input.csv -c Name -c Age -o output.csv
```

The contents of the created `output.csv`.

```
UserID,CompanyID
1,2
2,1
3,2
4,4
```

## rename

Create a new CSV file by rename columns from the input CSV file.

### Usage

```
csvt rename -i INPUT -c BEFORE_COLUMN1 ... -a AFTER_COLUMN1 ... -o OUTPUT
```

```
Usage:
  csvt rename [flags]

Flags:
  -i, --input string         Input CSV file path.
  -c, --column stringArray   Name of column before renaming.
  -a, --after stringArray    Name of column after renaming.
  -o, --output string        Output CSV file path.
  -h, --help                 help for rename
```

### Example

The contents of `input.csv`.

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,2
2,Hanako,21,1
```

Create `output.csv` by renmae "UserID" to "ID" and "CompanyID" to "Company" from `input.csv`.

```
$ csvt rename -i input.csv -c UserID -a ID -c CompanyID -a Company -o output.csv
```

The contents of the created `output.csv`.

```
ID,Name,Age,Company
1,"Taro, Yamada",10,2
2,Hanako,21,1
```

## replace

Create a new CSV file by replacing the values in the input CSV file.  
Regular expression are used for replace.

### Usage

```
csvt replace -i INPUT [[-c COLUMN1] ...] -r REGEX -t REPLACEMENT -o OUTPUT
```

```
Usage:
  csvt replace [flags]

Flags:
  -i, --input string         Input CSV file path.
  -c, --column stringArray   (optional) Name of the column to replace. If not specified, all columns are targeted.
  -r, --regex string         The regular expression to replace.
  -t, --replacement string   The string after replace.
  -o, --output string        Output CSV file path.
  -h, --help                 help for replace
```

### Example

The contents of `input.csv`.

```
col1,col2,col3
aa,abc,a1
bb,aabb,99
```

Create `output.csv` by replacing `a` with `x` in all columns.

```
$ csvt replace -i input.csv -r a -t x -o output.csv
```

The contents of the created `output.csv`.

```
col1,col2,col3
xx,xbc,x1
bb,xxbb,99
```

You can specify the target column with `-c`.  
Replacing `a` with `x` in "col1" and "col2".

```
$ csvt replace -i input.csv -c col1 -c col2 -r a -t x -o output.csv
```

```
col1,col2,col3
xx,xbc,a1
bb,xxbb,99
```

You can also use the capture group as `-t`.

```
$ csvt replace -i input.csv -c col3 -r ".*?([0-9]+)" -t "#$1" x -o output.csv
```

```
col1,col2,col3
aa,abc,#1
bb,aabb,#99
```

Please refer to the following for the syntax of regular expressions.

* https://golang.org/pkg/regexp/syntax/

## sort

Creates a new CSV file from the input CSV file by sorting by the values in the specified columns.

### Usage

```
csvt sort -i INPUT -c COLUMN1 ... [--desc] [--number] -o OUTPUT [--usingfile]
```

```
Usage:
  csvt sort [flags]

Flags:
  -i, --input string         Input CSV file path.
  -c, --column stringArray   Name of the column to use for sorting.
      --desc                 (optional) Sort in descending order. The default is ascending order.
      --number               (optional) Sorts as a number. The default is to sort as a string.
  -o, --output string        Output CSV file path.
      --usingfile            (optional) Use temporary files for sorting. Use this when sorting large files that will not fit in memory.
  -h, --help                 help for sort
```

### Example

The contents of `input.csv`.

```
col1,col2
02,a
10,b
01,a
11,c
20,b
```

Sort by "col1".

```
$ csvt sort -i input.csv -c col1 -o output.csv
```

The contents of the created `output.tsv`.

```
col1,col2
01,a
02,a
10,b
11,c
20,b
```

By default, it is sorted as a string.
For example, it could look like this

```
col1
1
12
123
2
21
3
```

If you want to sort as a number, specify `--number`.

```
$ csvt sort -i input.csv -c col1 --number -o output.csv
```

```
col1
1
2
3
12
21
123
```

## slice

Create a new CSV file by slicing the specified range of rows from the input CSV file.

### Usage

```
csvt slice -i INPUT [-s START] [-e END] -o OUTPUT
```

```
Usage:
  csvt slice [flags]

Flags:
  -i, --input string    Input CSV file path.
  -s, --start int       The number of the starting row. If not specified, it will be the first row. (default 1)
  -e, --end int         The number of the end row. If not specified, it will be the last row. (default 2147483647)
  -o, --output string   Output CSV file path.
  -h, --help            help for slice
```

### Example

The contents of `input.csv`.

```
ID,Name
1,name1
2,name2
3,name3
4,name4
5,name5
```

Slice the second through fourth records.

```
$ csvt slice -i input.csv -s 2 -e 4 -o output.csv
```

The contents of the created `output.tsv`.

```
ID,Name
2,name2
3,name3
4,name4
```

The `-s` and `-e` can be omitted.
If you want to extract the first row, it is sufficient to specify only `-e`, as shown below.

```
$ csvt slice -i input.csv -e 1 -o output.csv
```

```
ID,Name
1,name1
```

## transform

Transform the format of CSV file.

### Usage

```
csvt transform -i INPUT -o OUTPUT [--out-delim DELIMITER] [--out-quote QUOTE] [--out-sep SEPARATOR] [--out-allquote] [out-encoding ENCODING] [--out-bom]
```

```
Usage:
  csvt transform [flags]

Flags:
  -i, --input string          Input CSV file path.
  -o, --output string         Output CSV file path.
      --out-delim string      (optional) Output CSV delimiter. The default is ','
      --out-quote string      (optional) Output CSV quote. The default is '"'
      --out-sep string        (optional) Output CSV record separator. The default is CRLF.
      --out-allquote          (optional) Always quote output CSV fields. The default is to quote only the necessary fields.
      --out-encoding string   (optional) Output CSV encoding. The default is utf-8. Supported encodings: utf-8, shift_jis, euc-jp
      --out-bom               (optional) Output CSV with BOM.
  -h, --help                  help for transform
```

### Example

The contents of `input.csv`.

```
UserID,Name
1,"Taro, Yamada"
2,Hanako
```

Create `output.tsv` by transform `input.csv` to a TSV file.

```
$ csvt transform -i input.csv -o output.tsv --out-delim "\t"
```

The contents of the created `output.tsv`.

```
UserID  Name
1 Taro, Yamada
2 Hanako
```

Use common flag `--delim` to transform TSV file back to CSV file.

```
$ csvt transform -i output.tsv -o output2.csv --delim "\t"
```

## unique

Extracts unique records using the value of a specified columns.

### Usage

```
csvt unique -i INPUT -c COLUMN1 ... -o OUTPUT
```

```
Usage:
  csvt unique [flags]

Flags:
  -i, --input string         Input CSV file path.
  -c, --column stringArray   Name of the column to use for extract unique rows.
  -o, --output string        Output CSV file path.
  -h, --help                 help for unique
```

### Example

The contents of `input.csv`.

```
col1,col2
1,2
2,1
1,1
1,2
```

Extract the unique row in "col1".

```
$ csvt unique -i input.csv -c col1 -o output.tsv
```

The contents of the created `output.tsv`.

```
col1,col2
1,2
2,1
```

You can also specify multiple columns.  
Extract unique rows with "col1" and "col2".

```
$ csvt unique -i input.csv -c col1 -c col2 -o output.tsv
```

```
col1,col2
1,2
2,1
1,1
```

## Install

You can download the binary from the following.

* https://github.com/onozaty/csvt/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)
