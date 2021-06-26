# csvt

[![GitHub license](https://img.shields.io/github/license/onozaty/csvt)](https://github.com/onozaty/csvt/blob/main/LICENSE)
[![Test](https://github.com/onozaty/csvt/actions/workflows/test.yaml/badge.svg)](https://github.com/onozaty/csvt/actions/workflows/test.yaml)

`csvt` is a command line tool for processing CSV.

`csvt` consists of multiple subcommands.

* [choose](#choose) Choose columns from CSV file.
* [count](#count) Count the number of records in CSV file.
* [filter](#filter) Filter rows of CSV file.
* [header](#header) Show the header of CSV file.
* [join](#join) Join CSV files.
* [remove](#remove) Remove columns from CSV file.
* [rename](#rename) Rename columns from CSV file.
* [transform](#transform) Transform the format of CSV file.

## Common flags

Flags related to the CSV format are available in each subcommand as common flags.

```
Global Flags:
      --delim string   (optional) CSV delimiter. The default is ','
      --quote string   (optional) CSV quote. The default is '"'
      --sep string     (optional) CSV record separator. The default is CRLF.
      --allquote       (optional) Always quote CSV fields. The default is to quote only the necessary fields.
```

For example, when dealing with TSV files, change the delimiter to a tab as shown below.

```
$ csvt count -i INPUT --delim "\t"
```

## choose

Create a new CSV file by choosing columns from the input CSV file.

### Usage

```
$ csvt choose -i INPUT -c COLUMN1 -c COLUMN2 -o OUTPUT
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

## count

Count the number of records in CSV file.

### Usage

```
$ csvt count -i INPUT
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

## filter

Create a new CSV file by filtering the input CSV file to rows that match the conditions.

### Usage

```
$ csvt filter -i INPUT -c COLUMN -o OUTPUT
```

```
Usage:
  csvt filter [flags]

Flags:
  -i, --input string    Input CSV file path.
  -c, --column string   Name of the column to use for filtering. If neither --equal nor --regex is specified, it will filter by those with values.
  -o, --output string   Output CSV file path.
      --equal string    (optional) Filter by matching value.
      --regex string    (optional) Filter by regular expression.
  -h, --help            help for filter
```

### Example

The contents of `input.csv`.

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,2
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
1,"Taro, Yamada",10,2
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
1,"Taro, Yamada",10,2
4,Jun,22,2
```

Regular expressions can also be used.  
Use `--regex` to specify a regular expression.

```
$ csvt filter -i input.csv -c Name --regex [Yy]amada -o output.csv 
```

```
UserID,Name,Age,CompanyID
1,"Taro, Yamada",10,2
3,yamada,30,
```


## header

Show the header of CSV file.

### Usage

```
$ csvt header -i INPUT
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

## join

Join CSV files.  

Using the first CSV file as a base, join the contents of the second CSV file to create a new CSV file.  
It is similar to EXCEL's VLOOKUP.

### Usage

```
$ csvt join -1 INPUT1 -2 INPUT2 -c COLUMN -o OUTPUT
```

```
Usage:
  csvt join [flags]

Flags:
  -1, --first string     First CSV file path.
  -2, --second string    Second CSV file path.
  -c, --column string    Name of the column to use for joining.
      --column2 string   (optional) Name of the column to use for joining in the second CSV file. Specify if different from the first CSV file.
  -o, --output string    Output CSV file path.
      --usingfile        (optional) Use temporary files for joining. Use this when joining large files that will not fit in memory.
      --norecord         (optional) No error even if there is no record corresponding to sencod CSV.
  -h, --help             help for join
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

If the column name in the second CSV file is different from that in the first CSV file, specify it with `--column2`.

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
$ csvt remove -i INPUT -c COLUMN1 -c COLUMN2 -o OUTPUT
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
$ csvt rename -i INPUT -c BEFORE_COLUMN1 -a AFTER_COLUMN1 -c BEFORE_COLUMN2 -a AFTER_COLUMN2 -o OUTPUT
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

## transform

Transform the format of CSV file.

### Usage

```
$ csvt transform -i INPUT -o OUTPUT --out-delim DELIMITER --out-quote QUOTE --out-sep SEPARATOR --out-allquote
```

```
Usage:
  csvt transform [flags]

Flags:
  -i, --input string       Input CSV file path.
  -o, --output string      Output CSV file path.
      --out-delim string   (optional) Output CSV delimiter. The default is ','
      --out-quote string   (optional) Output CSV quote. The default is '"'
      --out-sep string     (optional) Output CSV record separator. The default is CRLF.
      --out-allquote       (optional) Always quote output CSV fields. The default is to quote only the necessary fields.
  -h, --help               help for transform
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

## Install

You can download the binary from the following.

* https://github.com/onozaty/csvt/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)
