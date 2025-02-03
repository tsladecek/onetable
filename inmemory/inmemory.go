package inmemory

import (
	"fmt"
)

type InMemory struct {
	data map[string][]byte
}

func New() *InMemory {
	data := make(map[string][]byte)

	return &InMemory{data: data}
}

func (t InMemory) Get(key string) ([]byte, error) {
	value, inside := t.data[key]

	if inside == false {
		return nil, fmt.Errorf("Value for Key %s not found", key)
	}

	return value, nil
}

func (t InMemory) Insert(key string, value []byte) error {
	t.data[key] = value
	return nil
}

func (t InMemory) Delete(key string) error {
	delete(t.data, key)
	return nil
}
