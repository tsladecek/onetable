# Onetable

Simple persistent key-value store

## How does it work?
To use the `OneTable` you need to choose an index.
The index is an in memory data structure used for fast lookups
of of value metadata (offset and length).

- The values are stored in an append only file, which does not make
much sense without the index
- The index data is stored also in an append only csv file in 
format `{key: string},{offset: int},{length: int}`

This allows for fast lookups and inserts without loading the
entire file content to memory

```go
index := onetable.NewIndexHashTable()
t, err := onetable.New("/path/to/folder/where/data/will/be/stored", index)
if err != nil {
    panic(err.Error())
}

_ = t.Insert("key", []byte("val0"))
v, _ := t.Get("key") // val = "val0"

_ = t.Insert("key", []byte("val1"))
v, _ = t.Get("key") // v = "val1"

_ = t.Delete("key")
v, _ = t.Get("key") // v = nil
```
---

You can also run a interactive session:

```shell
go run cmd/main.go --folder "/path/to/data"
```
