package onetable

import (
	"fmt"
	"sort"
)

type HashIndex map[string]ValueMetadata

type IndexHashTable struct {
	index HashIndex
}

func NewIndexHashTable() *IndexHashTable {
	index := make(HashIndex)
	return &IndexHashTable{index: index}
}

func (index *IndexHashTable) get(key string) (ValueMetadata, bool) {
	v, found := index.index[key]
	return v, found
}

func (index *IndexHashTable) insert(key string, valueMeta ValueMetadata) error {
	index.index[key] = valueMeta

	return nil
}

func (index *IndexHashTable) delete(key string) error {
	delete(index.index, key)
	return nil
}

func (index *IndexHashTable) between(fromKey string, toKey string) ([]*item, error) {
	keys := []string{}

	for k := range index.index {
		if k >= fromKey && k <= toKey {
			keys = append(keys, k)
		}
	}

	items := make([]*item, len(keys))
	sort.Strings(keys)

	for i, k := range keys {
		v, found := index.get(k)
		if !found {
			return nil, fmt.Errorf("Found no value for key %s", k)
		}
		items[i] = &item{key: k, value: v}
	}

	return items, nil
}
