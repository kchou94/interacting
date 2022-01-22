package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bernylinville/interacting/todo"
)

// Default file name
var todoFileName = ".todo.json"

func main() {
	// Parsing command line flags
	task := flag.String("task", "", "Task to be included in the Todo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Define an items list
	l := &todo.List{}

	// Check if the user define the ENV VAR for a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	// Use the Get method to read todo items from file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the number of arguments provided
	switch {
	// For no extra arguments, print the list
	case *list:
		// List current todo items
		fmt.Print(l)
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *task != "":
		l.Add(*task)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}
