package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cli_todolist/module/app"
)

const todoFile = "./data/todos.json"

// Set longer interval in real application -- or the terminal will be messed up
const checkInterval = 60 * time.Minute

// const checkInterval = 5 * time.Second

var mu sync.Mutex

func main() {
	fmt.Printf("\n[CLI_TODO_APP] Start Monitoring...\n")
	go monitorTasks()
	// Quit signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// If Sys Interrupt or ctrl-C to force terminate
	select {
	case <-c:
		fmt.Printf("\n[CLI_TODO_APP] Stop Monitoring...\n")

	}
}

func monitorTasks() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkAndPrintUrgentTasks()
		}
	}
}

func checkAndPrintUrgentTasks() {
	mu.Lock()
	defer mu.Unlock()

	todos := &app.Todos{}
	if err := todos.Load(todoFile); err != nil {
		handleError(err)
	}

	for _, task := range *todos {
		if !task.Done && task.Urgent {
			currentTime := time.Now().Format(time.RFC822)
			fmt.Printf("\n[%v] Reminder -> Urgent Task: %s\n", currentTime, task.Task)
		}
	}
}

func handleError(err error) {
	fmt.Println("Error:", err)
}
