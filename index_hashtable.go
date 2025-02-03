package onetable

type HashIndex map[typeKey]valueMetadata

type IndexHashTable struct {
	index HashIndex
}

func NewIndexHashTable() *IndexHashTable {
	index := make(HashIndex)
	return &IndexHashTable{index: index}
}

func (h *IndexHashTable) get(key typeKey) *valueMetadata {
	v, ok := h.index[key]

	if !ok {
		return nil
	}

	return &v
}

func (h *IndexHashTable) insert(key typeKey, valueMeta valueMetadata) error {
	h.index[key] = valueMeta

	return nil
}

func (h *IndexHashTable) delete(key typeKey) error {
	delete(h.index, key)
	return nil
}
