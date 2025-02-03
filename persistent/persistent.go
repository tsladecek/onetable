package persistent

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

type Offset int

type valueMetadata struct {
	offset Offset
	length int
}

const TOMBSTONE int = -1

type Index interface {
	get(key string) *valueMetadata
	insert(key string, value valueMetadata) error
	delete(key string) error
}

const (
	dataFileName  string = "data.ot"
	indexFileName string = "index.ot"
)

type Persistent struct {
	Path  string
	Index Index
	// TODO Filesystem lock
	lock      sync.Mutex
	offset    Offset
	dataPath  string
	indexPath string
}

func (p *Persistent) fillIndex(indexPath string) error {
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

		key := record[0]
		var offset Offset
		var offsetRaw, length int

		if offsetRaw, err = strconv.Atoi(record[1]); err != nil {
			return fmt.Errorf("Invalid record at line %d. Offset %s is not an integer", idx, record[1])
		}

		offset = Offset(offsetRaw)

		if length, err = strconv.Atoi(record[2]); err != nil {
			return fmt.Errorf("Invalid record at line %d. Length %s is not an integer", idx, record[2])
		}

		if length == TOMBSTONE {
			p.Index.delete(key)
			continue
		}

		p.Index.insert(key, valueMetadata{offset: offset, length: length})
	}
	return nil
}

func (p *Persistent) loadData() error {
	dataPath := path.Join(p.Path, dataFileName)
	indexPath := path.Join(p.Path, indexFileName)

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

	err := p.fillIndex(indexPath)
	if err != nil {
		panic(err.Error())
	}

	p.dataPath = dataPath
	p.indexPath = indexPath

	dataFile, err := os.Stat(dataPath)
	if err != nil {
		return err
	}

	p.offset = Offset(dataFile.Size())

	return nil
}

func New(folderPath string, index Index) (*Persistent, error) {
	p := &Persistent{Path: folderPath, Index: index}

	// check if path to data folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err != nil {
			panic(err.Error())
		}
	}

	// if there is data at dataPath, populate the inmemory index
	err := p.loadData()

	if err != nil {
		panic(err.Error())
	}

	return p, nil
}

func validateKey(key string) error {
	if strings.Contains(key, "\n") || strings.Contains(key, ",") {
		return errors.New("Invalid key. Contains one of forbidden characters: '\\n' or ','")
	}

	return nil
}

func (p *Persistent) writeValue(value []byte) error {
	f, err := os.OpenFile(p.dataPath, os.O_APPEND|os.O_WRONLY, 0644)
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

func (p *Persistent) writeKey(key string, valueMeta valueMetadata) error {
	f, err := os.OpenFile(p.indexPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{
		key, strconv.Itoa(int(valueMeta.offset)), strconv.Itoa(valueMeta.length)})

	w.Flush()

	return nil
}

func (p *Persistent) Insert(key string, value []byte) error {
	err := validateKey(key)
	if err != nil {
		return err
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	err = p.writeValue(value)
	if err != nil {
		return err
	}

	valueMeta := valueMetadata{offset: p.offset, length: len(value)}

	err = p.writeKey(key, valueMeta)
	if err != nil {
		return err
	}

	p.Index.insert(key, valueMeta)
	p.offset = p.offset + Offset(len(value))

	return nil
}

func (p *Persistent) Get(key string) ([]byte, error) {
	valueMeta := p.Index.get(key)

	if valueMeta == nil {
		return nil, nil
	}

	f, err := os.Open(p.dataPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	b := make([]byte, valueMeta.length)
	f.ReadAt(b, int64(valueMeta.offset))

	return b, nil
}

func (p *Persistent) Delete(key string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.writeKey(key, valueMetadata{offset: -1, length: TOMBSTONE})
	return p.Index.delete(key)
}
