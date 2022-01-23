package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bernylinville/interacting/todo"
)

// Default file name
var todoFileName = ".todo.json"

func main() {
	// Parsing command line flags
	add := flag.Bool("add", false, "Add task to Todo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("del", 0, "Delete task from Todo list")
	inspect := flag.Bool("inspect", false, "Enable verbose output")
	hideCompleted := flag.Bool("hidecompleted", false, "Hide completed tasks")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintf(flag.CommandLine.Output(), "AddNewTaskFromArguments, AddNewTaskFromSTDIN\n")
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
	case *add:
		// When any arguments (excluding flags) are provided, they will be
		// used as the new task
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, task := range t {
			l.Add(task)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *del > 0:
		// Delete the given item
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *inspect:
		for _, t := range *l {
			fmt.Printf("Task: %s\nDone: %t\nCreatedAt: %v\nCompletedAt: %v\n\n", t.Task, t.Done, t.CreatedAt, t.CompletedAt)
		}
	case *hideCompleted:
		for k, t := range *l {
			prefix := "  "
			if t.Done {
				continue
			}

			// Adjust the item number k to print numbers starting from 1 instead of 0
			fmt.Printf("%s%d: %s\n", prefix, k+1, t.Task)
		}
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// getTask function decides where to get the description for a new task from: arguments or STDIN
func getTask(r io.Reader, args ...string) ([]string, error) {
	tasks := []string{}

	if len(args) > 0 {
		tasks = append(tasks, strings.Join(args, " "))
		return tasks, nil
	}

	s := bufio.NewScanner(r)

	// Use for loop to read multiple lines from STDIN
	for s.Scan() {
		if err := s.Err(); err != nil {
			tasks = append(tasks, "")
			return tasks, err
		}

		if len(s.Text()) == 0 {
			tasks = append(tasks, "")
			return tasks, fmt.Errorf("Task cannot be blank")
		}

		tasks = append(tasks, s.Text())
	}

	return tasks, nil
}
