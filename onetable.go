package onetable

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

type typeOffset int

type ValueMetadata interface {
	Offset() typeOffset
	Length() int
}

type valueMetadata struct {
	offset typeOffset
	length int
}

func (v valueMetadata) Offset() typeOffset {
	return v.offset
}

func (v valueMetadata) Length() int {
	return v.length
}

type item struct {
	key   string
	value ValueMetadata
}

const tombstone int = -1

type Index interface {
	get(key string) (ValueMetadata, bool)
	insert(key string, value ValueMetadata) error
	delete(key string) error
	between(fromKey string, toKey string) ([]*item, error)
}

const (
	dataFileName  string = "data.ot"
	indexFileName string = "index.ot"
)

type OneTable struct {
	Path      string
	Index     Index
	lock      sync.Mutex
	offset    typeOffset
	dataPath  string
	indexPath string
}

func (o *OneTable) fillIndex(indexPath string) error {
	f, err := os.Open(indexPath)
	if err != nil {
		return err
	}

	defer f.Close()

	r := csv.NewReader(f)

	idx := 1
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if len(record) != 3 {
			return fmt.Errorf("Invalid record at line %d. Does not contain 3 separated fields", idx)
		}

		idx += 1

		key := string(record[0])
		var offset typeOffset
		var offsetRaw, length int

		if offsetRaw, err = strconv.Atoi(record[1]); err != nil {
			return fmt.Errorf("Invalid record at line %d. Offset %s is not an integer", idx, record[1])
		}

		offset = typeOffset(offsetRaw)

		if length, err = strconv.Atoi(record[2]); err != nil {
			return fmt.Errorf("Invalid record at line %d. Length %s is not an integer", idx, record[2])
		}

		if length == tombstone {
			o.Index.delete(key)
			continue
		}

		o.Index.insert(key, valueMetadata{offset: offset, length: length})
	}
	return nil
}

func (o *OneTable) loadData() error {
	dataPath := path.Join(o.Path, dataFileName)
	indexPath := path.Join(o.Path, indexFileName)

	// if data exists and index does not, panic
	_, dataFileErr := os.Stat(dataPath)
	_, indexFileErr := os.Stat(indexPath)

	if os.IsExist(dataFileErr) && os.IsNotExist(indexFileErr) {
		return errors.New("Index does not exist for data")
	}

	// if data file does not exist, create new files
	if os.IsNotExist(dataFileErr) {
		if _, err := os.Stat(indexPath); os.IsExist(err) {
			os.Remove(indexPath)
		}

		err := os.WriteFile(dataPath, []byte{}, 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(indexPath, []byte{}, 0644)
		if err != nil {
			return err
		}
	}

	err := o.fillIndex(indexPath)
	if err != nil {
		panic(err.Error())
	}

	o.dataPath = dataPath
	o.indexPath = indexPath

	dataFile, err := os.Stat(dataPath)
	if err != nil {
		return err
	}

	o.offset = typeOffset(dataFile.Size())

	return nil
}

func New(folderPath string, index Index) (*OneTable, error) {
	o := &OneTable{Path: folderPath, Index: index}

	// check if path to data folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err != nil {
			panic(err.Error())
		}
	}

	// if there is data at dataPath, populate the inmemory index
	err := o.loadData()

	if err != nil {
		panic(err.Error())
	}

	return o, nil
}

func validateKey(key string) error {
	if strings.Contains(string(key), "\n") || strings.Contains(string(key), ",") {
		return errors.New("Invalid key. Contains one of forbidden characters: '\\n' or ','")
	}

	return nil
}

func (o *OneTable) writeValue(value []byte) error {
	f, err := os.OpenFile(o.dataPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(value)

	if err != nil {
		return err
	}

	return nil
}

func (o *OneTable) writeKey(key string, valueMeta valueMetadata) error {
	f, err := os.OpenFile(o.indexPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{
		string(key),
		strconv.Itoa(int(valueMeta.offset)),
		strconv.Itoa(valueMeta.length),
	})

	w.Flush()

	return nil
}

func (o *OneTable) Insert(key string, value []byte) error {
	err := validateKey(key)
	if err != nil {
		return err
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	err = o.writeValue(value)
	if err != nil {
		return err
	}

	valueMeta := valueMetadata{offset: o.offset, length: len(value)}

	err = o.writeKey(key, valueMeta)
	if err != nil {
		return err
	}

	o.Index.insert(key, valueMeta)
	o.offset = o.offset + typeOffset(len(value))

	return nil
}

func (o *OneTable) Get(key string) ([]byte, error) {
	valueMeta, found := o.Index.get(key)

	if !found {
		return nil, nil
	}

	f, err := os.Open(o.dataPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	b := make([]byte, valueMeta.Length())
	f.ReadAt(b, int64(valueMeta.Offset()))

	return b, nil
}

func (o *OneTable) Delete(key string) error {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.writeKey(key, valueMetadata{offset: -1, length: tombstone})
	return o.Index.delete(key)
}
