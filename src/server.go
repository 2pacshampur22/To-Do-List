package main

import (
	"To-Do-List/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	store, err := storage.NewPostgresStorage("")
	if err != nil {
		fmt.Printf("Error connecting to Postgres: %s\n", err)
		return
	}
	fmt.Println("Connected to Postgres database")
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
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
		case "DELETE":
			deleteTaskId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/tasks/"))
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}
			err = store.Delete(deleteTaskId)
			if err != nil {
				http.Error(w, "Error deleting task", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		case "PUT":
			updateTaskId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/tasks/"))
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}
			err = store.Done(updateTaskId)
			if err != nil {
				http.Error(w, "Error updating task", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		}

	})
	fmt.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", enableCors(http.DefaultServeMux))
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
