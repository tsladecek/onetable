package onetable

import "testing"

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
