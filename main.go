package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	DB  *sql.DB
	err error
)

type Note struct {
	Id    string `json:"Id"`
	Title string `json:"Title"`
	Body  string `json:"Body"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func connection() *sql.DB {
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/notesapi-go")

	if err != nil {
		panic(err.Error())
	}

	return DB
}

var allNotes []Note

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hallo. Home Page is Here!")
}

func updateNotes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := DB.Prepare("UPDATE notes SET title = ?, body = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyValue := make(map[string]string)
	json.Unmarshal(body, &keyValue)
	newTitle := keyValue["title"]
	newBody := keyValue["body"]

	_, err = stmt.Exec(newTitle, newBody, params["id"])
	if err != nil {
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := Response{
		Code:    http.StatusOK,
		Message: "Data Update Successfully!",
	}
	resp, _ := json.Marshal(response)
	w.Write(resp)
}

func getAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var allNotes []Note

	result, err := DB.Query("SELECT * FROM notes")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var mynote Note
		err := result.Scan(&mynote.Id, &mynote.Title, &mynote.Body)
		if err != nil {
			panic(err.Error())
		}
		allNotes = append(allNotes, mynote)
	}

	json.NewEncoder(w).Encode(allNotes)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := DB.Query("SELECT * FROM notes WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var myNote Note

	for result.Next() {
		err := result.Scan(&myNote.Id, &myNote.Title, &myNote.Body)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(myNote)
}

func createNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var itemNotes Note
	_ = json.NewDecoder(r.Body).Decode(&itemNotes)

	insertNotes(DB, itemNotes)

	w.WriteHeader(http.StatusOK)
	response := Response{
		Code:    http.StatusOK,
		Message: "Success Created",
	}

	resp, _ := json.Marshal(response)
	w.Write(resp)
}

func deleteNotes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := DB.Prepare("DELETE FROM notes WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := Response{
		Code:    http.StatusOK,
		Message: "Data Delete Successfully!",
	}
	resp, _ := json.Marshal(response)
	w.Write(resp)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/api/v1/notes", createNotes).Methods("POST")
	router.HandleFunc("/api/v1/notes", getAllNotes).Methods("GET")
	router.HandleFunc("/api/v1/notes/{id}", getNote).Methods("GET")
	router.HandleFunc("/api/v1/notes/{id}", deleteNotes).Methods("DELETE")
	router.HandleFunc("/api/v1/notes/{id}", updateNotes).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func insertNotes(db *sql.DB, note Note) error {
	fmt.Println(db)
	sql := "INSERT INTO notes(title, body) values (?, ?)"
	result, err := db.Exec(sql, note.Title, note.Body)

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(result)

	return nil
}

func main() {
	connection()
	// defer DB.Close()

	err = DB.Ping()

	if err != nil {
		fmt.Println(err)
	}

	handleRequests()
}
