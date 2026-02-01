package main

import (
	"To-Do-List/storage"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	store := storage.NewStorage("tasks.json")
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			taskList, err := store.List()
			if err != nil {
				http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(taskList)
		case "POST":
			var newTask storage.Task
			err := json.NewDecoder(r.Body).Decode(&newTask)
			if err != nil {
				http.Error(w, "Invalid task data", http.StatusBadRequest)
				return
			}
			err = store.Add(newTask)
			if err != nil {
				http.Error(w, "Error adding task", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		}

	})
	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
