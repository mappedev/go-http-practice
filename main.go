package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Todo struct {
	ID   int
	Task string
}

var todos []Todo

func main() {
	// Tasks plural handler
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			allTodos := getAllTasks()
			json.NewEncoder(w).Encode(allTodos)

		case "POST":
			var todoBody Todo
			json.NewDecoder(r.Body).Decode(&todoBody)

			todo := createTask(todoBody)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(todo)
		}
	})

	// Tasks singular handler
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdFromPathUrl(r.URL.Path)

		if err == nil {
			switch r.Method {
			case "GET":
				todo, err := getTask(id)
				if err == nil {
					json.NewEncoder(w).Encode(todo)
					return
				}

			case "PUT":
				var todoBody Todo
				json.NewDecoder(r.Body).Decode(&todoBody)

				todo, err := updateTask(id, todoBody)
				if err == nil {
					json.NewEncoder(w).Encode(todo)
					return
				}

			case "DELETE":
				todo, err := deleteTask(id)
				if err == nil {
					json.NewEncoder(w).Encode(todo)
					return
				}
			}
		}

		notFoundTask(w)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}

// Services
func getAllTasks() []Todo {
	return todos
}

func getTask(id int) (Todo, error) {
	for _, todo := range todos {
		if todo.ID == id {
			return todo, nil
		}
	}

	return Todo{}, fmt.Errorf("error getting task %d", id)
}

func createTask(todo Todo) Todo {
	todos = append(todos, todo)
	return todo
}

func updateTask(id int, todoBody Todo) (Todo, error) {
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Task = todoBody.Task

			return todos[i], nil
		}
	}

	return Todo{}, fmt.Errorf("error updating task %d", id)
}

func deleteTask(id int) (Todo, error) {
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return todo, nil
		}
	}

	return Todo{}, fmt.Errorf("error removing task %d", id)
}

// Utils
func getIdFromPathUrl(path string) (int, error) {
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	return strconv.Atoi(idStr)
}

// Exceptions
func notFoundTask(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Task not found."))
}
