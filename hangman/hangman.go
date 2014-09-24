package hangman

import (
    "log"
    "net/http"
    "encoding/json"
//	"github.com/jhcook/game_engine/dictionary"
)

type Hangman struct {
    Word string
    Played []rune
    Game uint64
}

func Play(w http.ResponseWriter, r *http.Request) {
    game := Hangman{"fuck", []rune("Justin"), 1234567}
    if bytes, err := json.Marshal(game); err == nil {
        log.Println(string(bytes))
        w.Write(bytes)
    } else {
        log.Fatal(err)
    }
}
