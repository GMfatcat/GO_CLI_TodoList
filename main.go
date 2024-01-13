/* CLI TODO LIST */
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"cli_todolist/module/app"
)

const (
	todoFile = "./data/todos.json"
)

func handleError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func main() {

	urgent := flag.Bool("urgent", false, "mark as urgent todos, must add before -add flag")
	add := flag.Bool("add", false, "add a new todo")
	complete := flag.Int("complete", 0, "mark a todo as completed")
	del := flag.Int("del", 0, "delete a todo")
	cleanup := flag.Bool("cleanup", false, "clean up the todo list")
	list := flag.Bool("list", false, "list all todos")

	flag.Parse()

	todos := &app.Todos{}

	if err := todos.Load(todoFile); err != nil {
		handleError(err)
	}

	switch {
	case *urgent && *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			handleError(err)
		}
		// Urgent = true
		todos.Add(task, true)
		fmt.Printf("Urgent Task: %s", task)

	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			handleError(err)
		}
		// Urgent = false
		todos.Add(task, false)

	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			handleError(err)
		}
		err = todos.Store(todoFile)
		if err != nil {
			handleError(err)
		}

	case *del > 0:
		err := todos.Delete(*del)
		if err != nil {
			handleError(err)
		}
		err = todos.Store(todoFile)
		if err != nil {
			handleError(err)
		}

	case *cleanup:
		todos.CleanUp()

	case *list:
		todos.Print()
		return

	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}
	// Store after add / del / complete
	if err := todos.Store(todoFile); err != nil {
		handleError(err)
	}

}

func getInput(r io.Reader, args ...string) (string, error) {

	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}

	return text, nil

}
