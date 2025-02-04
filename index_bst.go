package onetable

type BSTNode struct {
	key   string
	value valueMetadata
	left  *BSTNode
	right *BSTNode
}

type IndexBST struct {
	root *BSTNode
}

func NewIndexBST() *IndexBST {
	return &IndexBST{}
}

func (i *IndexBST) get(key string) *valueMetadata {
	current := i.root

	for current != nil {
		if current.key == key {
			return &current.value
		}

		if current.key > key {
			current = current.left
		} else {
			current = current.right
		}
	}

	return nil
}

func (i *IndexBST) insert(key string, valueMeta valueMetadata) error {
	newNode := &BSTNode{key: key, value: valueMeta}

	if i.root == nil {
		i.root = newNode
		return nil
	}

	current := i.root

	for {
		if current.key == key {
			current.value = valueMeta
			break
		} else if current.key > key {
			if current.left == nil {
				current.left = newNode
				break
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = newNode
				break
			}
			current = current.right
		}
	}

	return nil
}

func (i *IndexBST) delete(key string) error {
	return nil
}
