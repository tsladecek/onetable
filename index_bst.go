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

		if key < current.key {
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
		} else if key < current.key {
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
	var parent *BSTNode
	current := i.root

	for current != nil {
		if current.key == key {
			break
		}

		parent = current

		if key < current.key {
			current = current.left
		} else {
			current = current.right
		}
	}

	if current == nil {
		return nil
	}

	if parent == nil {
		i.root = nil
		return nil
	}

	var replacement *BSTNode
	if current.left == nil && current.right == nil {
		// 1. Leaf Node
		replacement = nil
	} else if current.left != nil && current.right != nil {
		// 3. Two children
		// find smallest child of right node an use this as a replacement
		ptemp := current
		ctemp := current.right

		for ctemp.left != nil {
			ptemp = ctemp
			ctemp = ctemp.left
		}

		ptemp.left = nil
		ctemp.left = current.left
		ctemp.right = ptemp
		replacement = ctemp
	} else {
		// 2. One child
		replacement = current.left
		if replacement == nil {
			replacement = current.right
		}
	}

	if current.key < parent.key {
		parent.left = replacement
	} else {
		parent.right = replacement
	}
	return nil
}
