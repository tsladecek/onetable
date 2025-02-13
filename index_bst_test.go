package onetable

import (
	"sort"
	"testing"
)

func TestBSTInsert(t *testing.T) {
	bst := NewIndexBST()

	err := bst.insert("d", valueMetadata{})
	if err != nil {
		t.Fatal("Expected nil error after inserting root")
	}

	if bst.root.key != "d" {
		t.Fatal("root key not correct")
	}

	err = bst.insert("a", valueMetadata{})
	if err != nil {
		t.Fatal("Expected nil error after inserting 'a'")
	}

	if bst.root.left.key != "a" {
		t.Fatal("root.left key not correct")
	}

	err = bst.insert("f", valueMetadata{})
	if err != nil {
		t.Fatal("Expected nil error after inserting 'f'")
	}

	if bst.root.right.key != "f" {
		t.Fatal("root.right key not correct")
	}

	if bst.root.right.value.offset != 0 {
		t.Fatal("root.right.value.offset not 0")
	}

	offset := typeOffset(123)
	err = bst.insert("f", valueMetadata{offset: offset})

	if err != nil {
		t.Fatal("Expected nil error after inserting 'f'")
	}

	if bst.root.right.key != "f" {
		t.Fatal("root.right key not correct")
	}

	if bst.root.right.value.offset != offset {
		t.Fatalf("root.right.value.offset not %d", offset)
	}
}

func TestBSTGet(t *testing.T) {
	bst := NewIndexBST()

	v := bst.get("key")

	if v != nil {
		t.Fatal("Expecting nil value in empty tree")
	}

	for i, k := range []string{"d", "a", "b", "f"} {
		err := bst.insert(k, valueMetadata{offset: typeOffset(i)})

		if err != nil {
			t.Fatal(err.Error())
		}
	}

	for i, k := range []string{"d", "a", "b", "f"} {
		v = bst.get(k)
		if v == nil {
			t.Fatal("Expected non nil value")
		}

		if v.offset != typeOffset(i) {
			t.Fatalf("Wrong offset %d for key %s. Expected %d", v.offset, k, i)
		}
	}
}

func TestBSTDeleteLeaf(t *testing.T) {
	bst := NewIndexBST()
	err := bst.delete("key")

	if err != nil {
		t.Fatal("Expecting no error when deleting in empty tree")
	}

	for i, k := range []string{"d", "b", "a", "c1", "c0", "c2", "f", "e", "g"} {
		err := bst.insert(k, valueMetadata{offset: typeOffset(i)})

		if err != nil {
			t.Fatal(err.Error())
		}
	}

	// try removing leaf
	err = bst.delete("g")

	if err != nil {
		t.Fatal("Expecting no error when deleting node g")
	}

	if bst.root.right.right != nil {
		t.Fatal("Root.right.right node is not nil")
	}

	if bst.root.right.left.key != "e" {
		t.Fatal("Root.right.left node is not e")
	}

	v := bst.get("g")
	if v != nil {
		t.Fatal("Node g found in tree even after deletion")
	}
}

func TestBSTDeleteOneChild(t *testing.T) {
	bst := NewIndexBST()

	for i, k := range []string{"d", "b", "a", "c1", "c0", "c2", "f", "e"} {
		err := bst.insert(k, valueMetadata{offset: typeOffset(i)})

		if err != nil {
			t.Fatal(err.Error())
		}
	}

	// try removing node with only one children
	err := bst.delete("f")

	if err != nil {
		t.Fatal("Expecting no error when node f")
	}

	if bst.root.right.key != "e" {
		t.Fatal("Root.right node.key is not e")
	}

	v := bst.get("f")
	if v != nil {
		t.Fatal("Node f found in tree even after deletion")
	}
}

func TestBSTDeleteTwoChildren(t *testing.T) {
	bst := NewIndexBST()
	for i, k := range []string{"d", "b", "a", "c1", "c0", "c2", "f"} {
		err := bst.insert(k, valueMetadata{offset: typeOffset(i)})

		if err != nil {
			t.Fatal(err.Error())
		}
	}

	// try removing node with only one children
	err := bst.delete("b")

	if err != nil {
		t.Fatal("Expecting no error when node f")
	}

	if bst.root.left.key != "c0" {
		t.Fatal("Root.left node.key is not c0")
	}

	if bst.root.left.left == nil || bst.root.left.left.key != "a" {
		t.Fatal("Root.left.left node.key is not a")
	}

	if bst.root.left.right.key != "c1" {
		t.Fatal("Root.left.left node.key is not a")
	}

	v := bst.get("b")
	if v != nil {
		t.Fatal("Node b found in tree even after deletion")
	}
}

func TestBSTInorder(t *testing.T) {
	bst := NewIndexBST()
	items := []string{"d", "b", "a", "c", "f", "e", "g"}
	for i, k := range items {
		err := bst.insert(k, valueMetadata{offset: typeOffset(i)})

		if err != nil {
			t.Fatal(err.Error())
		}
	}

	res := &[]*BSTNode{}
	inorder(res, bst.root)

	sort.Strings(items)
	for i, v := range *res {
		if items[i] != v.key {
			t.Fatalf("Expected %s, Received %s", items[i], v.key)
		}
	}
}

func TestBSTBetween(t *testing.T) {
	index := NewIndexBST()

	if index == nil {
		t.Fatal("Index is nil after initialization")
	}

	for _, key := range []string{"a", "b", "c1", "c0", "c2", "c", "d", "e"} {
		err := index.insert(key, valueMetadata{})
		if err != nil {
			t.Fatalf("Failed to insert key %s", key)
		}
	}

	between, err := index.between("c", "d")

	if err != nil {
		t.Fatal("Valid index.between call failed")
	}

	expected := []string{"c", "c0", "c1", "c2", "d"}
	for i, item := range between {
		if item.key != expected[i] {
			t.Fatalf("Keys do not match. Expected %s, Got %s", expected[i], item.key)
		}
	}
}
