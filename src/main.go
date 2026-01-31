package main

import (
	"To-Do-List/storage"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var store storage.Service = storage.NewStorage("tasks.json")

	var input string
	fmt.Println("Please provide an input argument: add, list, done, delete")
	fmt.Println("If you want to add a task, provide name and description after add command using \"\" for both.")
	if len(os.Args) < 2 {
		fmt.Println("Please provide an input argument.")
		return
	}
	input = os.Args[1]

	switch input {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Please provide task name and description.")
			return
		}
		name := os.Args[2]
		description := os.Args[3]
		err := store.Add(storage.Task{
			Name:        name,
			Description: description,
			IsDone:      false,
		})
		if err != nil {
			fmt.Printf("Error adding task: %v\n", err)
			return
		}
		fmt.Println("Task added successfully")
	case "list":
		tasks, err := store.List()
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return
		}

		for _, task := range tasks {
			fmt.Printf("Task: %+v\n", task)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Please provide task ID to mark as done.")
			return
		}
		id := os.Args[2]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid task ID:", id)
			return
		}

		err = store.Done(idInt)
		if err != nil {
			fmt.Printf("Error marking task as done: %v\n", err)
			return
		}
		fmt.Printf("Marked task with ID %d as done\n", idInt)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide task ID to delete.")
			return
		}
		id := os.Args[2]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid task ID:", id)
			return
		}

		err = store.Delete(idInt)
		if err != nil {
			fmt.Printf("Error deleting task: %v\n", err)
			return
		}
		fmt.Printf("Deleted task with id: %v", id)

	default:
		fmt.Println("Unknown command, please use add, list, done, delete or update")
	}
}
