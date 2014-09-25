/*
Package hangman is an implementation of the game Hangman many of us played
as a child.

Author: Justin Cook <jhcook@gmail.com>
*/

package hangman

import (
	"encoding/json"
	"github.com/jhcook/game_engine/dictionary"
	"github.com/jhcook/game_engine/util"
	"log"
	"math/rand"
	"net/http"
	"time"
	"unicode/utf8"
)

var theBoys *men = nil
var dict *dictionary.Dictionary = nil

type hangman struct {
	Word   string
	WrdUni []rune
	Right  []rune
	Wrong  []rune
	Game   uint64
	Timer  *time.Timer
}

type Message struct {
	Cmd  string
	Gid  uint64
	Auth string
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
	wrd := dict.NextWord()
	theBoys.Games++
	game := &hangman{Word: wrd,
		WrdUni: util.StringToRuneArray(wrd),
		Right:  make([]rune, len(wrd)),
		Wrong:  make([]rune, 0, 256),
		Game:   theBoys.Games,
		Timer:  time.NewTimer(time.Second * 300)}
	theBoys.Episode[theBoys.Games] = game
	return game
}

func (g *hangman) evalChar(chr string) bool {
	var correct bool = false
	// Convert chr to rune
	letter, _ := utf8.DecodeRuneInString(chr)
	// See if rune is in the word
	for i, v := range g.WrdUni {
		if letter == v {
			g.Right[i] = letter
			correct = true
		}
	}
	if !correct {
		g.Wrong = append(g.Wrong, letter)
	}
	return correct
}

func NewGame() *hangman {
	h := NewHangman()
	h.Timer = time.NewTimer(time.Second * 60)
	// If the timer expires, remove the game from theBoys and free memory
	go func() {
		<-h.Timer.C
		log.Println("expired:", h.Game)
		if _, ok := theBoys.Episode[h.Game]; ok {
			delete(theBoys.Episode, h.Game)
		}
	}()
	return h
}

func Play(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON
	decoder := json.NewDecoder(r.Body)
	var msg Message
	var game *hangman
	if err := decoder.Decode(&msg); err != nil {
		log.Println("hangman.Play.Decode:", err)
		return
	}
	if len(msg.Cmd) > 0 && msg.Gid > 0 && len(msg.Auth) > 0 {
		if game, ok := theBoys.Episode[msg.Gid]; ok {
			log.Println(game)
		} else {
			return
		}
	} else if msg.Cmd == "NEW" {
		game = NewHangman()
	} else {
		return
	}
	// On each play, reset the timer
	game.Timer.Reset(time.Second * 300)
	// Use an anonymous structure to sanitize the data sent back.
	answer := struct {
		Word   string
		Missed []rune
		Game   uint64
	}{game.Word, game.Wrong, game.Game}
	// Send the result down the wire.
	if bytes, err := json.Marshal(answer); err == nil {
		log.Println(string(bytes))
		w.Write(bytes)
	} else {
		log.Fatalln(err)
	}
}
