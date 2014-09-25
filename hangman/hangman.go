package hangman

import (
	"encoding/json"
	"github.com/jhcook/game_engine/dictionary"
	"log"
	"net/http"
	"math/rand"
)

type hangman struct {
	Word   string
	Played []rune
	Game   uint64
}

func NewHangman() *hangman {
	d := ""
	dict := dictionary.NewDictionary(d)
	dict.Ci = uint32(rand.Intn(len(dict.Words)))
	game := &hangman{Word: dict.NextWord(),
					 Played: make([]rune, 16, 256)}
	log.Println(dict.Word)
	return game
}

/*
func (h *hangman) {
	;
}
*/

func Play(w http.ResponseWriter, r *http.Request) {
	h := NewHangman()
	log.Println(h)
	if bytes, err := json.Marshal(h); err == nil {
		log.Println(string(bytes))
		w.Write(bytes)
	} else {
		log.Fatal(err)
	}
}
