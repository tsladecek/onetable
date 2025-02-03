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

type valuePosition struct {
	offset int
	length int
}

type Index map[string]valuePosition

const (
	DataFileName  string = "data.ot"
	IndexFileName string = "index.ot"
)

type Persistent struct {
	Path string
	// TODO Filesystem lock
	lock      sync.Mutex
	index     Index
	offset    int
	dataPath  string
	indexPath string
}

func (p *Persistent) checkPath(pth string) error {
	if _, err := os.Stat(pth); os.IsNotExist(err) {
		return err
	}

	return nil
}

// Populate in memory index data structure from the index log.
// Index file contains newline separated lines
// in following format: {key: string},{offset: int64},{length: int}
func (p *Persistent) fillIndex() error {
	index := make(Index)
	p.index = index

	f, err := os.Open(p.indexPath)
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
		var offset, length int

		if offset, err = strconv.Atoi(record[1]); err != nil {
			return fmt.Errorf("Invalid record at line %d. Offset %s is not an integer", idx, record[1])
		}

		if length, err = strconv.Atoi(record[2]); err != nil {
			return fmt.Errorf("Invalid record at line %d. Length %s is not an integer", idx, record[2])
		}

		p.index[key] = valuePosition{offset: offset, length: length}
	}
	return nil
}

func (p *Persistent) loadData() error {
	dataPath := path.Join(p.Path, DataFileName)
	indexPath := path.Join(p.Path, IndexFileName)

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

	p.fillIndex()

	p.dataPath = dataPath
	p.indexPath = indexPath

	dataFile, err := os.Stat(dataPath)
	if err != nil {
		return err
	}

	p.offset = int(dataFile.Size())

	return nil
}

func New(folderPath string) (*Persistent, error) {
	p := &Persistent{Path: folderPath}

	// check if path to data folder exists
	err := p.checkPath(folderPath)
	if err != nil {
		panic(err.Error())
	}

	// if there is data at dataPath, populate the inmemory index
	err = p.loadData()

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

func (p *Persistent) Insert(key string, value []byte) error {
	err := validateKey(key)
	if err != nil {
		return err
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	f, err := os.OpenFile(p.dataPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(value)

	if err != nil {
		return err
	}

	fidx, err := os.OpenFile(p.indexPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fidx.Close()

	w := csv.NewWriter(fidx)
	w.Write([]string{key, strconv.Itoa(int(p.offset)), strconv.Itoa(len(value))})
	w.Flush()

	p.index[key] = valuePosition{offset: p.offset, length: len(value)}
	p.offset = p.offset + len(value)

	return nil
}

func (p *Persistent) Get(key string) ([]byte, error) {
	valueMetadata, inside := p.index[key]

	if !inside {
		return nil, nil
	}

	f, err := os.Open(p.dataPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	b := make([]byte, valueMetadata.length)
	f.ReadAt(b, int64(valueMetadata.offset))

	return b, nil
}
