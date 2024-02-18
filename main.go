package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocql/gocql"
)

type Todo struct {
	ID          gocql.UUID `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}

var session *gocql.Session

func main() {

	log.SetOutput(os.Stdout)

	log.Println("sarting...")
	var err error
	session, err = connectToScyllaDB()
	if err != nil {
		panic("Failed to connect to ScyllaDB: " + err.Error())
	}
	defer session.Close()
	insertSampleData()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/todos", createTodoHandler)
	http.HandleFunc("/todos/read", readTodosHandler)
	http.HandleFunc("/todos/update", updateTodoHandler)
	http.HandleFunc("/todos/delete", deleteTodoHandler)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

// Function to connect to ScyllaDB
func connectToScyllaDB() (*gocql.Session, error) {
	var cluster = gocql.NewCluster("node-0.gce-asia-south-1.1dee5732be0ae28fa763.clusters.scylla.cloud", "node-1.gce-asia-south-1.1dee5732be0ae28fa763.clusters.scylla.cloud", "node-2.gce-asia-south-1.1dee5732be0ae28fa763.clusters.scylla.cloud")
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "scylla", Password: "84bROzPv7kolEey"}
	cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy("GCE_ASIA_SOUTH_1")

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Function to insert sample data into the table
func insertSampleData() {
	todo1 := Todo{
		ID:          gocql.TimeUUID(),
		UserID:      "user1",
		Title:       "Complete the integration of Scylladb",
		Description: "after writing basic golang code connect scylladb",
		Status:      "pending",
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	todo2 := Todo{
		ID:          gocql.TimeUUID(),
		UserID:      "user1",
		Title:       "complete the implementation of golang in project",
		Description: "write the skelleton code to start porject",
		Status:      "completed",
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	if err := session.Query(
		"INSERT INTO todo.todos (id, user_id, title, description, status, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)",
		todo1.ID, todo1.UserID, todo1.Title, todo1.Description, todo1.Status, todo1.Created, todo1.Updated,
	).Exec(); err != nil {
		fmt.Println("Failed to insert sample data:", err.Error())
		return
	}

	if err := session.Query(
		"INSERT INTO todo.todos (id, user_id, title, description, status, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)",
		todo2.ID, todo2.UserID, todo2.Title, todo2.Description, todo2.Status, todo2.Created, todo2.Updated,
	).Exec(); err != nil {
		fmt.Println("Failed to insert sample data:", err.Error())
		return
	}

	fmt.Println("Sample data inserted successfully")
}

// HTTP handler to create a new TODO item
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	todo.ID = gocql.TimeUUID()
	todo.Created = time.Now()
	todo.Updated = time.Now()

	if err := session.Query(
		"INSERT INTO todo.todos (id, user_id, title, description, status, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)",
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated,
	).Exec(); err != nil {
		http.Error(w, "Failed to create TODO item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// HTTP handler to read TODO items with optional filtering by status
func readTodosHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	status := queryParams.Get("status")
	log.Print("value of status:", status)
	query := "SELECT id, user_id, title, description, status, created, updated FROM todo.todos"

	if status != "" {
		if status == "completed" {
			query = "SELECT * FROM todo.todos WHERE status = 'completed' ALLOW FILTERING;"
			log.Print("done completed")
		} else if status == "pending" {
			query = "SELECT * FROM todo.todos WHERE status = 'pending' ALLOW FILTERING;"
			log.Print("done pending")
		} else {
			http.Error(w, "Invalid status value", http.StatusBadRequest)
			return
		}
	} else {
		query += " ALLOW FILTERING"
	}

	iter := session.Query(query).Iter()

	var todos []Todo
	for {
		var todo Todo
		if !iter.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
			break
		}
		todos = append(todos, todo)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		http.Error(w, "Failed to read TODO items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
}

// HTTP handler to update a TODO item
func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	todo.Updated = time.Now()
	if err := session.Query(
		"UPDATE todo.todos SET title = ?, description = ?, status = ?, updated = ? WHERE id = ?",
		todo.Title, todo.Description, todo.Status, todo.Updated, todo.ID,
	).Exec(); err != nil {
		http.Error(w, "Failed to update TODO item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler to delete a TODO item
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todoID gocql.UUID
	err := json.NewDecoder(r.Body).Decode(&todoID)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := session.Query("DELETE FROM todo.todos WHERE id = ?", todoID).Exec(); err != nil {
		http.Error(w, "Failed to delete TODO item: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
