package main

import (
	"time"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"strconv"
	"fmt"
)

type Note struct {
		Title string `json:"title"`
		Description string `json:"description"`
		CreateAt time.Time `json:"create_at"`
}
var noteStore = make(map[string]Note)

var id int

func GetNoteHandler(w http.ResponseWriter, r *http.Request)  {
	var notes []Note
	for _, v := range noteStore {
			notes = append(notes, v)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(notes)
	if err != nil {
			panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func PostNoteHandler(w http.ResponseWriter, r *http.Request)  {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
			panic(err)
	}
	note.CreateAt = time.Now()
	id++
	k := strconv.Itoa(id)
	noteStore[k] = note

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"<h1>Hola Mundo</h1>")
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")

	server := &http.Server{
		Addr: 		":8080",
		Handler: 	r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1<<20,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
