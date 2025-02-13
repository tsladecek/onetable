package onetable

import (
	"crypto/rand"
	mrand "math/rand"
	"sort"
	"testing"
)

func BenchmarkIndexInsert(b *testing.B) {
	n := 1000
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = rand.Text()
	}

	indexHashTable := NewIndexHashTable()
	indexBST := NewIndexBST()

	b.Run("Hashtable insert",
		func(b *testing.B) {
			for b.Loop() {
				for i, key := range keys {
					err := indexHashTable.insert(key, valueMetadata{offset: typeOffset(i), length: 128})
					if err != nil {
						b.Fatal(err.Error())
					}
				}
			}
		},
	)

	b.Run("BST insert",
		func(b *testing.B) {
			for b.Loop() {
				for i, key := range keys {
					err := indexBST.insert(key, valueMetadata{offset: typeOffset(i), length: 128})
					if err != nil {
						b.Fatal(err.Error())
					}
				}
			}
		},
	)

	indexesToGet := make([]int, n)
	for i := 0; i < n; i++ {
		indexesToGet[i] = mrand.Intn(n)
	}

	b.Run("Hashtable get", func(b *testing.B) {
		for b.Loop() {
			for _, idx := range indexesToGet {
				_, found := indexHashTable.get(keys[idx])
				if !found {
					b.Fatal("Node not found")
				}
			}
		}
	})

	b.Run("BST get", func(b *testing.B) {
		for b.Loop() {
			for _, idx := range indexesToGet {
				_, found := indexBST.get(keys[idx])
				if !found {
					b.Fatal("Node not found")
				}
			}
		}
	})

	sort.Strings(keys)
	leftIdx := make([]int, n)
	rightIdx := make([]int, n)

	for i := 0; i < n; i++ {
		leftIdx[i] = mrand.Intn(n - indexesToGet[i])
		rightIdx[i] = leftIdx[i] + indexesToGet[i]
	}

	b.Run("Hashtable between", func(b *testing.B) {
		for b.Loop() {
			for i := 0; i < n; i++ {
				items, err := indexHashTable.between(keys[leftIdx[i]], keys[rightIdx[i]])
				if err != nil {
					b.Fatal(err.Error())
				}

				if len(items) != rightIdx[i]-leftIdx[i]+1 {
					b.Fatal("Items length does not equal range")
				}
			}
		}
	})

	b.Run("BST between", func(b *testing.B) {
		for b.Loop() {
			for i := 0; i < n; i++ {
				items, err := indexBST.between(keys[leftIdx[i]], keys[rightIdx[i]])

				if err != nil {
					b.Fatal(err.Error())
				}
				if len(items) != rightIdx[i]-leftIdx[i]+1 {
					b.Fatal("Items length does not equal range")
				}
			}
		}
	})

	b.Run("Hashtable delete and insert", func(b *testing.B) {
		for b.Loop() {
			for _, key := range keys {
				err := indexHashTable.delete(key)
				if err != nil {
					b.Fatal(err.Error())
				}
				err = indexHashTable.insert(key, valueMetadata{})
				if err != nil {
					b.Fatal(err.Error())
				}
			}
		}
	})

	b.Run("BST delete and insert", func(b *testing.B) {
		for b.Loop() {
			for _, key := range keys {
				err := indexBST.delete(key)
				if err != nil {
					b.Fatal(err.Error())
				}
				err = indexBST.insert(key, valueMetadata{})
				if err != nil {
					b.Fatal(err.Error())
				}
			}
		}
	})
}
