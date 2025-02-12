package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tsladecek/onetable"
)

func main() {
	pfolderPath := flag.String("folder", "", "Path to folder where data is/will be stored")
	pindex := flag.String("index", "hashtable", "Index to use. Currently supported: [hashtable, bst]")
	help := flag.Bool("help", false, "Print Help")

	flag.Parse()

	printHelp := func() {
		println("OneTable REPL")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *help {
		printHelp()
	}

	if pfolderPath == nil || *pfolderPath == "" || pindex == nil || *pindex == "" {
		printHelp()
	}

	if *pindex != "hashtable" && *pindex != "bst" {
		printHelp()
	}

	index := onetable.NewIndexHashTable()
	t, err := onetable.New(*pfolderPath, index)
	if err != nil {
		panic(err.Error())
	}
	reader := bufio.NewReader(os.Stdin)
	println("---Starting OneTable console---\n")
	println("Available commands:")
	println("get <key>")
	println("insert <key> <value>")
	println("delete <key>\n")

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}

		inputArr := strings.Split(strings.Trim(input, "\n"), " ")

		command := inputArr[0]
		if len(inputArr) < 2 {
			log.Println("Invalid input")
			continue
		}
		key := inputArr[1]

		if command == "get" {
			value, err := t.Get(key)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			} else {
				if value == nil {
					fmt.Printf(">Key '%s' not found\n", key)
					continue
				}
				fmt.Printf(">%s: %s\n", key, string(value))
			}
			continue
		}

		if command == "delete" {
			err := t.Delete(key)
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Printf(">Deleted key: %s\n", key)
			continue
		}

		if command == "insert" {
			if len(inputArr) < 3 {
				fmt.Println("Invalid insert instruction. Expected space separated key and value")
				continue
			}
			value := inputArr[2]
			err := t.Insert(key, []byte(value))

			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf(">Inserted %s: %s\n", key, value)
			continue
		}
		println("Invalid instruction")
	}
}
