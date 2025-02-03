package main

import (
	"github.com/tsladecek/onetable/inmemory"
	"github.com/tsladecek/onetable/persistent"
)

type OneTable interface {
	Get(key string) ([]byte, error)
	Insert(key string, value []byte) error
	Delete(key string) error
}

func main() {
	im := inmemory.New()
	t, err := persistent.New("/home/tomas/Downloads")
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
}
