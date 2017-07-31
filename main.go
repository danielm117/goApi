package main

import (
	"time"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"regexp"
)

type Game struct {
		Board string `json:"board"`
		Player string `json:"player"`
}

func NextMoveHandler(w http.ResponseWriter, r *http.Request)  {
	var game Game
	validGame := true
	w.Header().Set("Content-Type", "text/plain")
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	//check board size
	if len(game.Board) != 64 {
		validGame = false
		w.Write([]byte("invalid board size\n"))
	}

	//check game pieces
	pattern_w := regexp.MustCompile("w")
	matches_w := len(pattern_w.FindAllStringIndex(game.Board, -1))
	if matches_w > 12 {
		validGame = false
		w.Write([]byte("invalid number of w on board\n"))
	//	log.Println("invalid number of w on board\n")
	}
	pattern_b := regexp.MustCompile("b")
	matches_b := len(pattern_b.FindAllStringIndex(game.Board, -1))
	if matches_b > 12 {
		validGame = false
		w.Write([]byte("invalid number of b on board\n"))
	//	log.Println("invalid number of b on board\n")
	}
	//check pieces distribution
	j:= 7 //end of row
	for i:=0;i<len(game.Board);i++ {
		fmt.Print(string(game.Board[i]))
		fmt.Print(i)
		if i != j {
			if (string(game.Board[i]) != "0")&& (string(game.Board[i+1]) != "0"){
				//check right piece
				w.Write([]byte("invalid distribution of pieces on board\n"))
				validGame = false
				break
			}else if (len(game.Board)>i+8)&&(string(game.Board[i]) != "0")&& (string(game.Board[i+8]) != "0"){
				//check button piece
				w.Write([]byte("invalid distribution of pieces on board\n"))
				validGame = false
				break
			}
		} else {
			j+=8
			//check button piece
			if (j<=len(game.Board)) && (string(game.Board[i]) != "0")&& (string(game.Board[j]) != "0")  {
				w.Write([]byte("invalid distribution of pieces on board\n"))
				validGame = false
				break
			}
		}
	}

	//check valid player
	if game.Player != "w" && game.Player != "b" {
		validGame = false
		w.Write([]byte("invalid player"))
	//	log.Println("invalid player")
	}

	if validGame {
		w.Write([]byte("valid\n"))
	//	log.Println("valid\n")
	}
//	w.WriteHeader(http.StatusOK)
//	w.Write()

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"<h1>Hola Mundo</h1>")
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api/game", NextMoveHandler).Methods("POST")

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
