package onetable

import "testing"

func TestIndexHashTableBetween(t *testing.T) {
	index := NewIndexHashTable()

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
