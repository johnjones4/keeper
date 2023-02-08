package main

import (
	"os"
)

const usage = "Usage: todo [path]"

func main() {
	if len(os.Args) != 2 {
		panic(usage)
	}

	list := &todoList{path: os.Args[1]}
	err := list.open()
	if err != nil {
		panic(err)
	}

	rt := newRuntime(list)
	err = rt.run()
	if err != nil {
		panic(err)
	}
}
