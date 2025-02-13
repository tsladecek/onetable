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
					indexHashTable.insert(key, valueMetadata{offset: typeOffset(i), length: 128})
				}
			}
		},
	)

	b.Run("BST insert",
		func(b *testing.B) {
			for b.Loop() {
				for i, key := range keys {
					indexBST.insert(key, valueMetadata{offset: typeOffset(i), length: 128})
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
				indexHashTable.get(keys[idx])
			}
		}
	})

	b.Run("BST get", func(b *testing.B) {
		for b.Loop() {
			for _, idx := range indexesToGet {
				indexBST.get(keys[idx])
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
				indexHashTable.between(keys[leftIdx[i]], keys[rightIdx[i]])
			}
		}
	})

	b.Run("BST between", func(b *testing.B) {
		for b.Loop() {
			for i := 0; i < n; i++ {
				indexBST.between(keys[leftIdx[i]], keys[rightIdx[i]])
			}
		}
	})

	b.Run("Hashtable delete and insert", func(b *testing.B) {
		for b.Loop() {
			for _, key := range keys {
				indexHashTable.delete(key)
				indexHashTable.insert(key, valueMetadata{})
			}
		}
	})

	b.Run("BST delete and insert", func(b *testing.B) {
		for b.Loop() {
			for _, key := range keys {
				indexBST.delete(key)
				indexBST.insert(key, valueMetadata{})
			}
		}
	})
}
