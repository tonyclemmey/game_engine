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

// Create a new game of hangman, add it to the singleton 'theBoys' and return
// to the caller.
func NewHangman() *hangman {
	dict.Ci = uint32(rand.Intn(len(dict.Words)))
	theBoys.Games++
	game := &hangman{Word: dict.NextWord(),
		Played: make([]rune, 0, 256),
		Game:   theBoys.Games,
		Timer:  time.NewTimer(time.Second * 300)}
	theBoys.Episode[theBoys.Games] = game
	return game
}

func Play(w http.ResponseWriter, r *http.Request) {
	// Just for initial development create a new game on each play attempt.
	h := NewHangman()
	// On each play, reset the timer
	h.Timer = time.NewTimer(time.Second * 60)
	// If the timer expires, remove the game from theBoys and free memory
	go func() {
		<-h.Timer.C
		log.Println("expired:", h.Game)
		if _, ok := theBoys.Episode[h.Game]; ok {
			delete(theBoys.Episode, h.Game)
		}
	}()
	// In the real world we will use an anonymous structure to sanitize the
	// data sent back.
	answer := struct {
		Word   string
		Missed []rune
		Game   uint64
	}{h.Word, h.Played, h.Game}
	// Send the result down the wire.
	if bytes, err := json.Marshal(answer); err == nil {
		log.Println(string(bytes))
		w.Write(bytes)
	} else {
		log.Fatalln(err)
	}
}
