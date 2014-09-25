package hangman

import (
	"encoding/json"
	"github.com/jhcook/game_engine/dictionary"
//	"github.com/jhcook/game_engine/util"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var theBoys *men = nil
var dict *dictionary.Dictionary = nil

type hangman struct {
	Word   string
	Played []rune
	Game   uint64
	Timer  *time.Timer
}

type men struct {
	Games   uint64
	Episode map[uint64]*hangman
}

// Instantiates 'theBoys' if not already instantiated as this is intended to
// be singleton.
func NewMen() {
	if theBoys != nil {
		log.Println("NewMen: already instantiated")
		return
	}
	dict = dictionary.NewDictionary("")
	theBoys = &men{Episode: make(map[uint64]*hangman)}
}

func NewHangman() *hangman {
	dict.Ci = uint32(rand.Intn(len(dict.Words)))
	theBoys.Games++
	game := &hangman{Word: dict.NextWord(),
		Played: make([]rune, 16, 256),
		Game:   theBoys.Games,
		Timer:	time.NewTimer(time.Second*300)}
	theBoys.Episode[theBoys.Games] = game
	return game
}

func Play(w http.ResponseWriter, r *http.Request) {
//	h := NewHangman()
	answer := struct{Word string}{"Justin"}
	if bytes, err := json.Marshal(answer); err == nil {
		log.Println(string(bytes))
		w.Write(bytes)
	} else {
		log.Fatalln(err)
	}
}
