package main

import (
	"time"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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

		output := ""
		valid_move := ""
		end_row := 7
		start_row := 0
		player := game.Player
		opponent := "w"
		if player == "w" {
			opponent = "b"
		}
		empty_place := "0"
		board := game.Board

		if player == "b" {
			for i := start_row; i < len(board); i++ {
				if string(board[i]) == player {
					if i != start_row {
						//check corner left
						if i+7 < len(board) && string(board[i+7]) == opponent {
							if i+14 < len(board) && (i+15)%8 != 0 && string(board[i+14]) == empty_place {
								output = strconv.Itoa(i+1) + "-" + strconv.Itoa(i+15)
								break
							}
						} else if i+7 < len(board) && string(board[i+7]) == empty_place && valid_move == "" {
							valid_move = strconv.Itoa(i+1) + "-" + strconv.Itoa(i+8)
						}
					} else {
						start_row += 8
					}
					if i != end_row {
						//check corner right
						if i+9 < len(board) && string(board[i+9]) == opponent {
							if i+18 < len(board) && (i+18)%8 != 0 && (string(board[i+18]) == empty_place) {
								output = strconv.Itoa(i+1) + "-" + strconv.Itoa(i+19)
								break
							}
						} else if i+9 < len(board) && string(board[i+9]) == empty_place && valid_move == "" {
							valid_move = strconv.Itoa(i+1) + "-" + strconv.Itoa(i+10)
						}
					} else {
						end_row += 8
					}
				}
			}
		} else {
			for i := start_row; i < len(board); i++ {
				if string(board[i]) == player {
					if i != end_row {
						//check corner right
						if i-7 >= 0 && string(board[i-7]) == opponent {
							if i-14 >= 0 && (i-14)%8 != 0 && (string(board[i-14]) == empty_place) {
								output = strconv.Itoa(i+1) + "-" + strconv.Itoa(i-13)
								break
							}
						} else if i-7 >= 0 && string(board[i-7]) == empty_place && valid_move == "" {
							valid_move = strconv.Itoa(i+1) + "-" + strconv.Itoa(i-6)
						}
					} else {
						end_row += 8
					}
					if i != start_row {
						//check corner left
						if i-9 >= 0 && string(board[i-9]) == opponent {
							if i-18 >= 0 && (i-17)%8 != 0 && string(board[i-18]) == empty_place {
								output = strconv.Itoa(i+1) + "-" + strconv.Itoa(i-17)
								break
							}
						} else if i-9 >= 0 && string(board[i-9]) == empty_place && valid_move == "" {
							valid_move = strconv.Itoa(i+1) + "-" + strconv.Itoa(i-8)
						}
					} else {
						start_row += 8
					}
				}
			}
		}
		if output != "" {
		//	fmt.Println(output)
			w.Write([]byte(output))
		} else if valid_move != "" {
		//	fmt.Println(valid_move)
			w.Write([]byte(valid_move))
		} else {
			w.Write([]byte("No se encontro movimiento posible"))
		}
	}

	//	w.WriteHeader(http.StatusOK)

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
