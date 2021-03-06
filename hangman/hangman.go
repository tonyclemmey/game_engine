/*
Package hangman is an implementation of the game Hangman many of us played
as a child.

Author: Justin Cook <jhcook@gmail.com>
*/

package hangman

import (
	"golang.org/x/net/websocket"
	"encoding/json"
	"github.com/jhcook/game_engine/dictionary"
	"github.com/jhcook/game_engine/util"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"regexp"
	"unicode"
	"unicode/utf8"
)

var theBoys *men = nil
var dict *dictionary.Dictionary = nil

// This is used to represent each game
type hangman struct {
	Word   string      // Word represented as an ascii string
	Defo   string      // Definition of Word for hint
	WrdUni []rune      // Word represented as unicode points
	Right  []rune      // Correct tries
	Wrong  []rune      // Incorrect tries
	Game   uint64      // Game ID
	Cmd    string      // Last destructive command used
	TwoP   bool        // Is this a two-player game
	P1cred string      // Credentials for player one
	P2cred string      // Credentials for player two
	Timer  *time.Timer // Timer to remove the game from memory
}

// This is used as the structure for JSON decoding
type Message struct {
	Cmd  string
	Gid  uint64
	Play string
	Auth string
}

// This is a structure that holds references to each game instance
type men struct {
	Games   uint64
	Episode map[uint64]*hangman
}

// Instantiates 'theBoys' if not already instantiated as this is intended to
// be singleton.
func NewMen() {
	if theBoys != nil {
		log.Println(util.GetFuncName(), ": already instantiated")
		return
	}
	dict = dictionary.NewDictionary("")
	theBoys = &men{Episode: make(map[uint64]*hangman)}
}

// This function returns a rune array with non-Latin characters filled in
func createRuneArray(wrd string) []rune {
    stuff := make([]rune, len(wrd))
    for ind, runeCh := range wrd {
        if !unicode.In(runeCh, unicode.Latin) {
            stuff[ind] = runeCh
        }
    }
    /* This is ugly, but it is necessary to check the end of the word since
       multiple codes are provided when a word may be a noun, verb, etc. */
    var i int
    for i = len(stuff); i>=0; i-- {
        if stuff[i-1] != 0 {
            continue
        } else {
            break
        }
    }
    return stuff[:i]
}

// Create a new game of hangman, add it to the singleton 'theBoys' and return
// to the caller.
func NewHangman(np int) *hangman {
	var p2 string
	var err error
	var dictEntry dictionary.DictEntry

GETWORD:
	dict.Ci = uint32(rand.Intn(len(dict.Words)))
	dictEntry.Word = dict.NextWord()
	if err = dictEntry.GetDefinition(); err != nil {
		log.Printf("%s.GetDefinition: %v\n", util.GetFuncName(), err)
		err = nil
	}
	for dictEntry.Definition == "" {
		goto GETWORD
	}
	theBoys.Games++
	if np == 2 {
		p2 = util.Rand_str(64)
	} else {
		p2 = ""
	}

	// Remove the word from the hint
	re, _ := regexp.Compile("(?i)" + dictEntry.Word + `[\w']*`)
	res2 := re.ReplaceAllString(dictEntry.Definition, "-----")

	game := &hangman{
		Word:   dictEntry.Word,
		Defo:   res2,
		WrdUni: util.StringToRuneArray(dictEntry.Word),
//		Right:  make([]rune, len(dictEntry.Word)),
        Right:  createRuneArray(dictEntry.Word),
		Wrong:  make([]rune, 0, 6),
		Game:   theBoys.Games,
		Cmd:    "NEW",
		P1cred: util.Rand_str(64),
		P2cred: p2,
		Timer:  time.NewTimer(time.Second * 60)}
	go func() {
		<-game.Timer.C
		log.Printf("%v expired: %d\n", util.GetFuncName(), game.Game)
		if _, ok := theBoys.Episode[game.Game]; ok {
			delete(theBoys.Episode, game.Game)
		}
	}()
	theBoys.Episode[theBoys.Games] = game
	return game
}

func (g *hangman) evalChar(chr string) bool {
	var correct bool = false
	// Convert chr to rune
	letter, _ := utf8.DecodeRuneInString(strings.ToLower(chr))
	// See if rune is in the word
	for i, v := range g.WrdUni {
		if letter == unicode.ToLower(v) {
			g.Right[i] = v
			correct = true
		}
	}
	if !correct {
		g.Wrong = append(g.Wrong, letter)
	}
	return correct
}

func (g *hangman) checkAuth(gid uint64, cred string) *hangman {
	if game, ok := theBoys.Episode[gid]; ok {
		if game.P1cred == cred {
			// On each play, reset the timer
			game.Timer.Reset(time.Second * 300)
			return game
		}
	}
	return nil
}

func (msg *Message) play() interface{} {
	var answer interface{} = nil
	var game *hangman
	var cmd string = ""
	var ok bool = false

	switch msg.Cmd {
	case "NEW":
		game = NewHangman(1)
		answer = struct {
			Cmd    string
			Hint   string
			Curr   []rune
			Missed []rune
			Game   uint64
			Cred   string
		}{game.Cmd, game.Defo, game.Right, game.Wrong, game.Game, game.P1cred}
		log.Println(game)

	case "P1T", "P2T":
		if _, ok = theBoys.Episode[msg.Gid]; ok {
			if game = game.checkAuth(msg.Gid, msg.Auth); game != nil {
				cmd = "P1T"
                // msg.Play length should be one char
                if len(msg.Play) == 1 {
				    game.evalChar(msg.Play)
                } else {
                    answer = struct {
                        Error string
                    }{"incorrect number of characters attempted"}
                }
			} else {
				answer = struct {
					Error string
				}{"unauthorized"}
			}
		} else {
			answer = struct {
				Error string
			}{"unknown game id"}
		}

    /* In the case where the user tries to guess the word, the "FIN" command
       is sent along with the string version of the answer. */
    case "FIN":
        if _, ok = theBoys.Episode[msg.Gid]; ok {
            if game = game.checkAuth(msg.Gid, msg.Auth); game != nil {
                cmd = "FIN"
                for i, w := 0, 0; i < len(msg.Play); i += w {
                    runeValue, width := utf8.DecodeRuneInString(msg.Play[i:])
                    if !game.evalChar(string(runeValue)) {
                        answer = nil
                    }
                    w = width
                }
                game.evalChar(msg.Play)
            } else {
                answer = struct {
                    Error string
                }{"unauthorized"}
            }
        } else {
            answer = struct {
                Error string
            }{"unknown game id"}
        }

	case "STATUS":
		if game = game.checkAuth(msg.Gid, msg.Auth); game != nil {
			cmd = "STATUS"
		} else {
			answer = struct {
				Error string
			}{"unauthorized"}
		}

	default:
		answer = struct {
			Error string
		}{"unknown command"}
	}

    var right []rune
	if answer == nil {
		// Use an anonymous structure to sanitize the data sent back.
        if len(game.Wrong) >5 {
            cmd = "FIN"
            right = game.WrdUni
        } else if cmd == "FIN" {
            cmd = "FIN"
            right = game.WrdUni
        } else {
            right = game.Right
        }
		answer = struct {
		    Cmd    string
			Curr   []rune
	    	Missed []rune
			Game   uint64
    	}{cmd, right, game.Wrong, game.Game}
	}
	return answer
}

func Playhttp(w http.ResponseWriter, r *http.Request) {
	var msg Message
	var answer interface{} = nil
	// Decode the JSON
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		bytes, _ := json.Marshal(err)
		log.Printf("%s.Decode: %v\n", util.GetFuncName(), err)
		w.Write(bytes)
		return
	}
	answer = msg.play()
	// Send the result down the wire.
	if bytes, err := json.Marshal(answer); err == nil {
		log.Println(string(bytes))
		w.Write(bytes)
	} else {
		log.Println(err)
	}
}

func Playws(ws *websocket.Conn) {
	var answer interface{} = nil
	var err error = nil
	var msg Message
	for {
		if err = websocket.JSON.Receive(ws, &msg); err != nil {
			log.Printf("%s.Receive: %v\n", util.GetFuncName(), err)
			break
		}
		answer = msg.play()
		if err := websocket.JSON.Send(ws, answer); err != nil {
			log.Printf("%s.Send: %v\n", util.GetFuncName(), err)
			break
		}
	}
}
