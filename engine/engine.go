package main

import (
	"github.com/jhcook/game_engine/hangman"
	"log"
	"net/http"
	"code.google.com/p/go.net/websocket"
)

func apiHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v: %v from %v", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

func apiHandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return apiHandler(next)
}

func main() {
	// Create the singleton "men" in hangman to keep track of all hangman episodes
	hangman.NewMen()
	//    http.Handle("/list", apiHandlerFunc())
	//    http.Handle("/new", apiHandlerFunc())
	//    http.Handle("/status", apiHandlerFunc())
	// For websockets we need to use HandleFunc to override the http Origin header.
	// http://stackoverflow.com/questions/19708330/serving-a-websocket-in-go
	// Otherwise, we will get a 403 and waste a night trying to debug :)
	http.HandleFunc("/wshangman", func(w http.ResponseWriter, req *http.Request) {
			s := websocket.Server{Handler: websocket.Handler(hangman.Playws)}
			s.ServeHTTP(w, req)
		})
		//}websocket.Handler(hangman.Playws))
	// Provide an HTTP interface
	http.Handle("/hangman", apiHandlerFunc(hangman.Playhttp))
	http.ListenAndServe(":3000", nil)
}
