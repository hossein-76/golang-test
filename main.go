package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"log"
)
var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	"localhost", 5432, "hossein", "abcd@1234", "go_db")
type User struct {
	age  int
	email string
	first_name string
	last_name string
}
func mainPage(w http.ResponseWriter, r *http.Request){
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `
	INSERT INTO users (age, email, first_name, last_name)
	VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, user.age, user.email, user.first_name, user.last_name)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(`{"asdf":"asdf"}`)
}

func getPage(w http.ResponseWriter, r *http.Request){
	var email string
	id := 1
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT id, email FROM users WHERE id=$1;`
	row := db.QueryRow(sqlStatement, 1)
	switch err := row.Scan(&id, &email); err {
	case sql.ErrNoRows:
	  fmt.Println("No rows were returned!")
	case nil:
	  fmt.Println(id, email)
	default:
	  panic(err)
	}
	data := fmt.Sprintf(`{"email": %s, "id": %d}`, email, id)
	json.NewEncoder(w).Encode(data)
}
func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", mainPage).Methods("POST")
	myRouter.HandleFunc("/", getPage).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

