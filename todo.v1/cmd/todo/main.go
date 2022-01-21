package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bernylinville/interacting/todo"
)

// Hardcoding the file name
const todoFileName = ".todo.json"

func main() {
	// Parsing command line flags
	task := flag.String("task", "", "Task to be included in the Todo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")

	flag.Parse()

	// Define an items list
	l := &todo.List{}

	// Use the Get method to read todo items from file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the number of arguments provided
	switch {
	// For no extra arguments, print the list
	case len(os.Args) == 1:
		// List current todo items
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	// Concatenate all provided arguments with a space and
	// add to the list as an item
	default:
		// Concatenate all arguments with a space
		item := strings.Join(os.Args[1:], " ")
		l.Add(item)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
