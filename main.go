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
var (
	player1 string = "b"
	player2 string = "w"
	empty_place string = "0"
)

func CheckBoardHandler(w http.ResponseWriter, r *http.Request)  {
	var game Game
	w.Header().Set("Content-Type", "text/plain")
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
	}

	validBoardFormat := ValidateBoardFormat(game.Board)
	validPiecesDistribution := ValidatePiecesDistribution(game.Board)
	validPlayer := ValidatePlayer(game.Player)
	if !validBoardFormat{
		w.Write([]byte("invalid board format"))
	}else if !validPlayer {
		w.Write([]byte("invalid player"))
	} else if !validPiecesDistribution {
		w.Write([]byte("invalid distribution of pieces on board"))
	}else {
		w.Write([]byte("Perfect board"))
	}
}

func NextMoveHandler(w http.ResponseWriter, r *http.Request)  {
	var game Game
	w.Header().Set("Content-Type", "text/plain")
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
	}

	validBoardFormat := ValidateBoardFormat(game.Board)
	validPiecesDistribution := ValidatePiecesDistribution(game.Board)
	validPlayer := ValidatePlayer(game.Player)

	if validBoardFormat && validPiecesDistribution && validPlayer {
		var output string
		if game.Player == player1 {
			output = SearchMovePlayer1(game)
		} else {
			output = SearchMovePlayer2(game)
		}
		w.Write([]byte(output))
	}else {
		w.Write([]byte("Check inputs with others services "))
	}
}


func ValidateBoardFormat(board string) bool {

	patternPlayer1 := regexp.MustCompile(player1)
	matchesPlayer1 := len(patternPlayer1.FindAllStringIndex(board, -1))
	patternEmptyPlace := regexp.MustCompile(empty_place)
	matchesEmptyPlace := len(patternEmptyPlace.FindAllStringIndex(board, -1))
	patternPlayer2 := regexp.MustCompile(player2)
	matchesPlayer2 := len(patternPlayer2.FindAllStringIndex(board, -1))

	//check pieces number
	if matchesPlayer1 > 12 || matchesPlayer2 > 12 || (matchesPlayer1+matchesPlayer2+matchesEmptyPlace)!=64 {
		return false
	} else {
		return true
	}
}

func ValidatePiecesDistribution(board string) bool {
	is_valid := true
	end_row:= 7

	for i:=0;i<len(board);i++ {
		if i != end_row {
			if (string(board[i]) != empty_place)&& (string(board[i+1]) != empty_place){
				//check right piece
				is_valid = false
				break
			}else if (len(board)>i+8)&&(string(board[i]) != empty_place)&& (string(board[i+8]) != empty_place){
				//check button piece
				is_valid = false
				break
			}
		} else {
			end_row+=8
			//check button piece
			if (end_row<=len(board)) && (string(board[i]) != empty_place)&& (string(board[end_row]) != empty_place)  {
				is_valid = false
				break
			}
		}
	}
	return is_valid
}

func ValidatePlayer(player string) bool {
	//check valid player
	if player != player1 && player != player2 {
		return false
	}else {
		return true
	}
}

func SearchMovePlayer1(game Game) string {
	//the player 1 move down on the board
	player := game.Player
	board := game.Board
	opponent := player2

	output := ""
	valid_move := ""
	end_row := 7
	start_row := 0

	for i := start_row; i < len(board); i++ {
		if string(board[i]) == player {
			if i != start_row {
				//check corner left
				if i+7 < len(board) && string(board[i+7]) == opponent {
					if i+14 < len(board) && (i+15)%8 != 0 && string(board[i+14]) == empty_place {
						output = strconv.Itoa(i+1) + "-" + strconv.Itoa(i+15)
						break
					}
				} else if i+7 < len(board) && i%8 != 0 &&  string(board[i+7]) == empty_place && valid_move == "" {
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
				} else if i+9 < len(board) && (i+1)%8 != 0 && string(board[i+9]) == empty_place && valid_move == "" {
					valid_move = strconv.Itoa(i+1) + "-" + strconv.Itoa(i+10)
				}
			} else {
				end_row += 8
			}
		}
	}
	if output != "" {
		return output
	} else if valid_move != "" {
		return valid_move
	} else {
		return "No found movements"
	}
}

func SearchMovePlayer2(game Game) string {
	//the player 2 move to up on board
	player := game.Player
	board := game.Board
	opponent := player1

	output := ""
	valid_move := ""
	end_row := 7
	start_row := 0

	for i := start_row; i < len(board); i++ {
		if string(board[i]) == player {
			if i != end_row {
				//check corner right
				if i-7 >= 0 && string(board[i-7]) == opponent {
					if i-14 >= 0 && (i-14)%8 != 0 && (string(board[i-14]) == empty_place) {
						output = strconv.Itoa(i+1) + "-" + strconv.Itoa(i-13)
						break
					}
				} else if i-7 >= 0 && (i-7)%8 != 0 && string(board[i-7]) == empty_place && valid_move == "" {
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
				} else if i-9 >= 0 && i%8 != 0 && string(board[i-9]) == empty_place && valid_move == "" {
					valid_move = strconv.Itoa(i+1) + "-" + strconv.Itoa(i-8)
				}
			} else {
				start_row += 8
			}
		}
	}
	if output != "" {
		return output
	} else if valid_move != "" {
		return valid_move
	} else {
		return "No found movements"
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"<h1>Hola Mundo</h1>")
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api/game", NextMoveHandler).Methods("POST")
	r.HandleFunc("/api/board", CheckBoardHandler).Methods("POST")

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
