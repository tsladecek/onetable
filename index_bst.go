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

func (index *IndexBST) get(key string) *valueMetadata {
	current := index.root

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

func (index *IndexBST) insert(key string, valueMeta valueMetadata) error {
	newNode := &BSTNode{key: key, value: valueMeta}

	if index.root == nil {
		index.root = newNode
		return nil
	}

	current := index.root

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

func (index *IndexBST) delete(key string) error {
	var parent *BSTNode
	current := index.root

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
		index.root = nil
		return nil
	}

	var replacement *BSTNode
	if current.left != nil && current.right != nil {
		// Two children
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

func inorder(buffer *[]*BSTNode, node *BSTNode) {
	if node == nil {
		return
	}

	if node.left != nil {
		inorder(buffer, node.left)
	}

	*buffer = append(*buffer, node)

	if node.right != nil {
		inorder(buffer, node.right)
	}
}

func inorderBetween(buffer *[]*item, node *BSTNode, fromKey string, toKey string) {
	if node == nil {
		return
	}

	if node.left != nil {
		inorderBetween(buffer, node.left, fromKey, toKey)
	}

	if node.key >= fromKey && node.key <= toKey {
		*buffer = append(*buffer, &item{key: node.key, value: node.value})
	}

	if node.right != nil {
		inorderBetween(buffer, node.right, fromKey, toKey)
	}
}

func (index *IndexBST) between(fromKey string, toKey string) ([]*item, error) {
	res := &[]*item{}
	inorderBetween(res, index.root, fromKey, toKey)
	return *res, nil
}
