package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB

func initDB() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	log.Println("Successfully connected to the database")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	err := db.QueryRow(
		"INSERT INTO users(name, email) VALUES($1, $2) RETURNING id",
		user.Name, user.Email).Scan(&user.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name, &user.Email)
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	_, err := db.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3", user.Name, user.Email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initDB()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", router)
}
