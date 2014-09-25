package hangman

import (
	"encoding/json"
	"github.com/jhcook/game_engine/dictionary"
	"log"
	"math/rand"
	"net/http"
)

var m *men
var dict *dictionary.Dictionary

type hangman struct {
	Word   string
	Played []rune
	Game   uint64
}

type men struct {
	Games   uint64
	Episode map[uint64]*hangman
}

func NewMen() {
	dict = dictionary.NewDictionary("")
	m = &men{Episode: make(map[uint64]*hangman)}
}

func NewHangman() *hangman {
	dict.Ci = uint32(rand.Intn(len(dict.Words)))
	m.Games++
	game := &hangman{Word: dict.NextWord(),
		Played: make([]rune, 16, 256),
		Game:   m.Games}
	m.Episode[m.Games] = game
	return game
}

func Play(w http.ResponseWriter, r *http.Request) {
	h := NewHangman()
	if bytes, err := json.Marshal(h); err == nil {
		log.Println(string(bytes))
		w.Write(bytes)
	} else {
		log.Fatalln(err)
	}
}
