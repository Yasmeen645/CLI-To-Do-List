package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	
	"time"
)

// Task represents a to-do item
type Task struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Done     bool      `json:"done"`
	Deadline time.Time `json:"deadline,omitempty"`
}

// loadTasks reads tasks from tasks.txt file
func loadTasks() ([]Task, error) {
	file, err := os.ReadFile("tasks.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(file, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// saveTasks writes tasks to tasks.txt file
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.txt", data, 0644)
}

// addTask creates a new task and adds it to the list
func addTask(tasks []Task, title string, deadline string) ([]Task, int) {
	var newID int
	if len(tasks) == 0 {
		newID = 1
	} else {
		maxID := tasks[0].ID
		for _, task := range tasks[1:] {
			if task.ID > maxID {
				maxID = task.ID
			}
		}
		newID = maxID + 1
	}

	var dl time.Time
	if deadline != "" {
		parsed, err := time.Parse("2006-01-02", deadline)
		if err == nil {
			dl = parsed
		}
	}

	newTask := Task{
		ID:       newID,
		Title:    title,
		Done:     false,
		Deadline: dl,
	}

	tasks = append(tasks, newTask)
	return tasks, newID
}

// deleteTask removes a task by ID
func deleteTask(tasks []Task, id int) ([]Task, bool) {
	for i, task := range tasks {
		if task.ID == id {
			return append(tasks[:i], tasks[i+1:]...), true
		}
	}
	return tasks, false
}

// markDone sets a task as done by ID
func markDone(tasks []Task, id int) ([]Task, bool) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			return tasks, true
		}
	}
	return tasks, false
}

// clearTasks removes all tasks
func clearTasks() []Task {
	return []Task{}
}

// printUsage shows available commands
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  add \"task name\" [deadline YYYY-MM-DD] - Add a new task with optional deadline")
	fmt.Println("  list                                  - List all tasks")
	fmt.Println("  delete <id>                           - Delete a task by ID")
	fmt.Println("  done <id>                             - Mark a task as done by ID")
	fmt.Println("  clear                                 - Delete all tasks")
}

// Colors
var (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

func main() {
	// Load existing tasks
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	// Check command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Task title is required")
			printUsage()
			os.Exit(1)
		}
		title := os.Args[2]
		var deadline string
		if len(os.Args) > 3 {
			deadline = os.Args[3]
		}
		var newID int
		tasks, newID = addTask(tasks, title, deadline)
		fmt.Printf("%sAdded task #%d:%s %s\n", green, newID, reset, title)

	case "list":
		if len(tasks) == 0 {
			fmt.Println(yellow + "No tasks found" + reset)
			break
		}
		fmt.Println("Tasks:")
		for _, task := range tasks {
			status := red + "Not Done" + reset
			if task.Done {
				status = green + "Done" + reset
			}
			dl := ""
			if !task.Deadline.IsZero() {
				dl = " (Deadline: " + task.Deadline.Format("2006-01-02") + ")"
			}
			fmt.Printf("#%d: %s [%s]%s\n", task.ID, task.Title, status, dl)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Task ID is required")
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: ID must be a number")
			os.Exit(1)
		}
		var found bool
		tasks, found = deleteTask(tasks, id)
		if !found {
			fmt.Printf("Error: Task #%d not found\n", id)
			os.Exit(1)
		}
		fmt.Printf("%sDeleted task #%d%s\n", red, id, reset)

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Error: Task ID is required")
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: ID must be a number")
			os.Exit(1)
		}
		var found bool
		tasks, found = markDone(tasks, id)
		if !found {
			fmt.Printf("Error: Task #%d not found\n", id)
			os.Exit(1)
		}
		fmt.Printf("%sMarked task #%d as done%s\n", green, id, reset)

	case "clear":
		tasks = clearTasks()
		fmt.Println(yellow + "All tasks cleared!" + reset)

	default:
		printUsage()
		os.Exit(1)
	}

	// Save tasks if modified
	if command == "add" || command == "delete" || command == "done" || command == "clear" {
		if err := saveTasks(tasks); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
			os.Exit(1)
		}
	}
}
