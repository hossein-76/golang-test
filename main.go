package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	// "strings"
)
var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	"localhost", 5432, "hossein", "abcd@1234", "go_db")
type User struct {
	Age  int `"json":"age"`
	Email string `"json":"email"`
	First_name string `"json":"first_name"`
	Last_name string `"json":"last_name"`
}
func mainPage(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.Body)
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var user User
	err = json.Unmarshal(b, &user)
	fmt.Println(user)
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
	_, err = db.Exec(sqlStatement, user.Age, user.Email, user.First_name, user.Last_name)
	if err != nil {
		panic(err)
	}
	data := fmt.Sprintf(`{"email": %s, "age": %d}`, user.Email, user.Age)

	json.NewEncoder(w).Encode(data)
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

