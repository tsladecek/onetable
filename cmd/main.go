package main

import "github.com/tsladecek/onetable"

func main() {
	index := onetable.NewIndexHashTable()
	t, err := onetable.New("/home/tomas/Downloads", index)
	if err != nil {
		panic(err.Error())
	}

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
