package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

var allNotes []Note

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hallo. Home Page is Here!")
}

func deleteNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	_deleteItemId(params["id"])

	json.NewEncoder(w).Encode(allNotes)
}

func _deleteItemId(id string) {
	for index, item := range allNotes {
		if item.Id == id {
			allNotes = append(allNotes[:index], allNotes[index+1:]...)
			break
		}
	}
}

func updateNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var itemNotes Note
	_ = json.NewDecoder(r.Body).Decode(&itemNotes)

	params := mux.Vars(r)

	_deleteItemId(params["id"])
	allNotes = append(allNotes, itemNotes)

	json.NewEncoder(w).Encode(allNotes)
}

func getAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(allNotes)
}

func createNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var itemNotes Note
	_ = json.NewDecoder(r.Body).Decode(&itemNotes)

	allNotes = append(allNotes, itemNotes)

	w.WriteHeader(http.StatusOK)
	response := Response{
		Code: http.StatusOK,
		Message: "Success Created",
	}

	resp,_ := json.Marshal(response)
	w.Write(resp)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/api/v1/notes", getAllNotes).Methods("GET")
	router.HandleFunc("/api/v1/notes", createNotes).Methods("POST")
	router.HandleFunc("/api/v1/notes/{id}", deleteNotes).Methods("DELETE")
	router.HandleFunc("/api/v1/notes/{id}", updateNotes).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func main() {
	allNotes = append(allNotes, Note{
		Id:    "123",
		Title: "Golang",
		Body:  "Golang adalah bahasa dari Google",
	})
	handleRequests()
}
