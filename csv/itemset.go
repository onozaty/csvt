package csv

import (
	"fmt"
	"io"

	"github.com/onozaty/csvt/util"
)

type ItemSet struct {
	items map[string]struct{}
}

// 入れておく値は何でも良い
var itemValue = struct{}{}

func (hashset *ItemSet) Add(item string) {
	hashset.items[item] = itemValue
}

func (hashset *ItemSet) Contains(item string) bool {
	_, contains := hashset.items[item]
	return contains
}

func (hashset *ItemSet) Count() int {
	return len(hashset.items)
}

func NewItemSet() *ItemSet {
	return &ItemSet{
		items: make(map[string]struct{}),
	}
}

func LoadItemSet(reader CsvReader, targetColumnName string) (*ItemSet, error) {

	columnNames, err := reader.Read()
	if err != nil {
		return nil, err
	}
	targetColumnIndex := util.IndexOf(columnNames, targetColumnName)
	if targetColumnIndex == -1 {
		return nil, fmt.Errorf("%s is not found", targetColumnName)
	}

	itemSet := NewItemSet()

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		itemSet.Add(row[targetColumnIndex])
	}

	return itemSet, nil
}
