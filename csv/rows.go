package csv

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/onozaty/csvt/util"
)

type CsvSortedRows interface {
	Count() int
	Row(index int) ([]string, error)
	Close() error
}

type memorySortedRows struct {
	rows [][]string
}

func (t *memorySortedRows) Count() int {

	return len(t.rows)
}

func (t *memorySortedRows) Row(index int) ([]string, error) {

	return t.rows[index], nil
}

func (t *memorySortedRows) Close() error {

	// リソースは保持しないので何もしない
	return nil
}

func LoadMemorySortedRows(reader CsvReader, useColumnNames []string, compare func(item1 string, item2 string) (int, error)) (CsvSortedRows, error) {

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	useColumnIndexes := []int{}
	for _, useColumnName := range useColumnNames {

		useColumnIndex := util.IndexOf(headers, useColumnName)
		if useColumnIndex == -1 {
			return nil, fmt.Errorf("%s is not found", useColumnName)
		}

		useColumnIndexes = append(useColumnIndexes, useColumnIndex)
	}

	rows := [][]string{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	var sortError error
	// ソート
	sort.SliceStable(rows, func(i, j int) bool {

		if sortError != nil {
			// エラーが起きているときは以降の比較は行わない
			return false
		}

		n := 0

		for _, useColumnIndex := range useColumnIndexes {

			n, sortError = compare(rows[i][useColumnIndex], rows[j][useColumnIndex])
			if sortError != nil {
				return false
			}

			if n != 0 {
				break
			}
		}

		return n < 0
	})

	if sortError != nil {
		return nil, sortError
	}

	return &memorySortedRows{
		rows: rows,
	}, nil
}

type fileSortedRows struct {
	sortedIndexies []int
	dbPath         string
	db             *bolt.DB
}

func (t *fileSortedRows) Count() int {
	return len(t.sortedIndexies)
}

func (t *fileSortedRows) Row(index int) ([]string, error) {

	// 既にDBを開いている場合は、使いまわす
	// (CsvTableのClose時に閉じている)
	if t.db == nil {
		db, err := bolt.Open(t.dbPath, 0600, nil)
		if err != nil {
			return nil, err
		}
		t.db = db
	}

	row := make([]string, 0)

	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("csvRows"))

		v := b.Get([]byte(strconv.Itoa(t.sortedIndexies[index])))
		if v != nil {
			json.Unmarshal(v, &row)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(row) == 0 {
		return nil, nil
	}

	return row, nil
}

func (t *fileSortedRows) Close() error {

	if t.db != nil {
		err := t.db.Close()
		if err != nil {
			return err
		}
	}

	return os.Remove(t.dbPath)
}

type SortSource struct {
	index int
	items []string
}

func LoadFileSortedRows(reader CsvReader, useColumnNames []string, compare func(item1 string, item2 string) (int, error)) (CsvSortedRows, error) {

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	useColumnIndexes := []int{}
	for _, useColumnName := range useColumnNames {

		useColumnIndex := util.IndexOf(headers, useColumnName)
		if useColumnIndex == -1 {
			return nil, fmt.Errorf("%s is not found", useColumnName)
		}

		useColumnIndexes = append(useColumnIndexes, useColumnIndex)
	}

	dbFile, err := os.CreateTemp("", "csvdb")
	if err != nil {
		return nil, err
	}
	defer dbFile.Close()

	db, err := bolt.Open(dbFile.Name(), 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sortSources := []SortSource{}
	rowIndex := 0
	eof := false

	for !eof {

		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("csvRows"))
			if err != nil {
				return err
			}

			// 1トランザクションで大量の書き込みを行うと速度が落ちるため
			// 分割してコミットを行う
			for i := 0; i < 10000; i++ {
				row, err := reader.Read()
				if err == io.EOF {
					eof = true
					break
				}
				if err != nil {
					return err
				}

				items := []string{}
				for _, useColumnIndex := range useColumnIndexes {
					items = append(items, row[useColumnIndex])
				}

				sortSources = append(sortSources, SortSource{
					index: rowIndex,
					items: items,
				})

				rowJson, err := json.Marshal(row)
				if err != nil {
					return err
				}

				err = b.Put([]byte(strconv.Itoa(rowIndex)), []byte(rowJson))
				if err != nil {
					return err
				}

				rowIndex++
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	var sortError error
	// ソート
	sort.SliceStable(sortSources, func(i, j int) bool {

		if sortError != nil {
			// エラーが起きているときは以降の比較は行わない
			return false
		}

		n := 0

		for index := range useColumnIndexes {

			n, sortError = compare(sortSources[i].items[index], sortSources[j].items[index])
			if sortError != nil {
				return false
			}

			if n != 0 {
				break
			}
		}

		return n < 0
	})

	if sortError != nil {
		return nil, sortError
	}

	sortedIndexies := []int{}
	for _, sortSource := range sortSources {
		sortedIndexies = append(sortedIndexies, sortSource.index)
	}

	return &fileSortedRows{
		sortedIndexies: sortedIndexies,
		dbPath:         dbFile.Name(),
	}, nil
}

func CompareString(item1 string, item2 string) (int, error) {
	if item1 == item2 {
		return 0, nil
	}
	if item1 < item2 {
		return -1, nil
	}
	return 1, nil
}

func CompareNumber(item1 string, item2 string) (int, error) {

	num1, err := strconv.Atoi(item1)
	if err != nil {
		return 0, err
	}

	num2, err := strconv.Atoi(item2)
	if err != nil {
		return 0, err
	}

	return num1 - num2, nil
}
