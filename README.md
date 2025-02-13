# Onetable

Simple persistent key-value store

## How does it work?
To use the `OneTable` you need to choose an index.
The index is an in memory data structure used for fast lookups
of value metadata (offset and length).

- The values are stored in an append only file, which does not make
much sense without the index
- The index data is stored also in an append only csv file in 
format `{key: string},{offset: int},{length: int}`

This allows for fast lookups and inserts without loading the
entire file content to memory

```go
index := onetable.NewIndexHashTable()
// or 
index := onetable.NewIndexBST()
t, err := onetable.New("/path/to/folder/where/data/will/be/stored", index)
if err != nil {
    panic(err.Error())
}

err = t.Insert("a", []byte("val a"))
err = t.Insert("b", []byte("val b"))
err = t.Insert("c", []byte("val c"))
v, found := t.Get("a") // val a

// get sorted values in range
items, err := t.between("a", "b") // []{Key: string, Value: []byte}

// delete key
err = t.delete("c")

```
---

You can also run a repl session:

```shell
go run cmd/repl/main.go --folder "/path/to/data"
```
