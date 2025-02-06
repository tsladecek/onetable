package onetable

import (
	"fmt"
	"sort"
)

type HashIndex map[string]valueMetadata

type IndexHashTable struct {
	index HashIndex
}

func NewIndexHashTable() *IndexHashTable {
	index := make(HashIndex)
	return &IndexHashTable{index: index}
}

func (index *IndexHashTable) get(key string) *valueMetadata {
	v, ok := index.index[key]

	if !ok {
		return nil
	}

	return &v
}

func (index *IndexHashTable) insert(key string, valueMeta valueMetadata) error {
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
		v := index.get(k)
		if v == nil {
			return nil, fmt.Errorf("Found no value for key %s", k)
		}
		items[i] = &item{key: k, value: *v}
	}

	return items, nil
}
