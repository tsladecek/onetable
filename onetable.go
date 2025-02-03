package main

import (
	"github.com/tsladecek/onetable/inmemory"
	"github.com/tsladecek/onetable/persistent"
)

func main() {
	im := inmemory.New()

	index := persistent.NewIndexHashTable()
	t, err := persistent.New("/home/tomas/Downloads", index)
	if err != nil {
		panic(err.Error())
	}
	im.Insert("key", []byte("val"))

	_ = t.Insert("key", []byte("val0"))
	v, _ := t.Get("key")
	println(string(v))

	_ = t.Insert("key", []byte("val1"))
	v, _ = t.Get("key")
	println(string(v))

	_ = t.Delete("key")
	v, _ = t.Get("key")
	println(string(v))
}
