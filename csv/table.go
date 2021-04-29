package csv

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/boltdb/bolt"
	"github.com/onozaty/csvt/util"
)

type CsvTable interface {
	Find(key string) (map[string]string, error)
	JoinColumnName() string
	ColumnNames() []string
	Close() error
}

type MemoryTable struct {
	joinColumnName string
	columnNames    []string
	rows           map[string][]string
}

func (t *MemoryTable) Find(key string) (map[string]string, error) {

	row := t.rows[key]

	if row == nil {
		return nil, nil
	}

	rowMap := make(map[string]string)
	for i := 0; i < len(t.columnNames); i++ {
		rowMap[t.columnNames[i]] = row[i]
	}

	return rowMap, nil
}

func (t *MemoryTable) JoinColumnName() string {

	return t.joinColumnName
}

func (t *MemoryTable) ColumnNames() []string {

	return t.columnNames
}

func (t *MemoryTable) Close() error {

	// リソースは保持しないので何もしない
	return nil
}

func LoadCsvMemoryTable(reader CsvReader, joinColumnName string) (CsvTable, error) {

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	primaryColumnIndex := util.IndexOf(headers, joinColumnName)
	if primaryColumnIndex == -1 {
		return nil, fmt.Errorf("%s is not found", joinColumnName)
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
			return nil, fmt.Errorf("%s:%s is duplicated", joinColumnName, row[primaryColumnIndex])
		}

		rows[row[primaryColumnIndex]] = row
	}

	return &MemoryTable{
		joinColumnName: joinColumnName,
		columnNames:    headers,
		rows:           rows,
	}, nil
}

type FileTable struct {
	joinColumnName string
	columnNames    []string
	dbPath         string
	db             *bolt.DB
}

func (t *FileTable) Find(key string) (map[string]string, error) {

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

		v := b.Get([]byte(key))
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

	rowMap := make(map[string]string)
	for i := 0; i < len(t.columnNames); i++ {
		rowMap[t.columnNames[i]] = row[i]
	}

	return rowMap, nil
}

func (t *FileTable) JoinColumnName() string {

	return t.joinColumnName
}

func (t *FileTable) ColumnNames() []string {

	return t.columnNames
}

func (t *FileTable) Close() error {

	if t.db != nil {
		err := t.db.Close()
		if err != nil {
			return err
		}
	}

	return os.Remove(t.dbPath)
}

func LoadCsvFileTable(reader CsvReader, joinColumnName string) (CsvTable, error) {

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	primaryColumnIndex := util.IndexOf(headers, joinColumnName)
	if primaryColumnIndex == -1 {
		return nil, fmt.Errorf("%s is not found", joinColumnName)
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

				key := row[primaryColumnIndex]

				// 格納前に既にあるか確認
				// -> 重複して存在した場合はエラーに
				v := b.Get([]byte(key))
				if v != nil {
					return fmt.Errorf("%s:%s is duplicated", joinColumnName, key)
				}

				rowJson, err := json.Marshal(row)
				if err != nil {
					return err
				}

				err = b.Put([]byte(key), []byte(rowJson))
				if err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return &FileTable{
		joinColumnName: joinColumnName,
		columnNames:    headers,
		dbPath:         dbFile.Name(),
	}, nil
}
