package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type ToDo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

var db *sql.DB
var mutex = &sync.Mutex{}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"task" TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v\n", err)
	}
}

func addTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newTask ToDo
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	result, err := db.Exec("INSERT INTO todos (task) VALUES (?)", newTask.Task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting task into database: %v", err), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	newTask.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func viewTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	rows, err := db.Query("SELECT id, task FROM todos")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []ToDo
	for rows.Next() {
		var task ToDo
		err := rows.Scan(&task.ID, &task.Task)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows: %v", err), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask ToDo
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	_, err = db.Exec("UPDATE todos SET task = ? WHERE id = ?", updatedTask.Task, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating task: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Task with ID %d updated successfully\n", id)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	_, err = db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting task: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Task with ID %d deleted successfully\n", id)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/add", addTask)
	http.HandleFunc("/view", viewTasks)
	http.HandleFunc("/update", updateTask)
	http.HandleFunc("/delete", deleteTask)
	http.Handle("/", http.FileServer(http.Dir("."))) 

	fmt.Println("Server is running on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

