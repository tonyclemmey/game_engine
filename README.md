Game Server Engine
==============

This is the game server that is implemented at a basic level. It has wrapper
capabilities that will be used to cater to various other plugins. At this time,
it is not a plugin framework. Rather, there is one game that is integrated
in to the code base. 

Language
--------------

The server is implemented in Go using HTTP and JSON. To build, checkout the
code, navigate to the repo's root and `go install`. If your Go environment is
configured properly, it will work by `$GOPATH/bin/game_engine`. The process listens
on port 3000 by default, but this can be changed in the code.

Games
--------------

The only game written at this time is Hangman. A Python cli and web clients
are provided. After the man is hanged, the word will be returned. The Cambridge
Dictionary API is used to provide a hint. After you build, it will be available
locally -- with configured hostname -- at:

http://hangman.example.com/hangman.html

In the example above, websocket is listening on port 8080 and is proxied by
mod_proxy_ws in Apache.

Configuration
--------------

There are no configuration capabilities at this time. 

Protocol
--------------

The protocol is RESTful or Websocket over HTTP and uses JSON structures to
both send and receive data in a syncronous style.

The flow will look like the following:

Received for a new game 
`{"Cmd":"NEW","Hint":"a hint to the word","Curr":[0,0,0,0,0,0,0,0,0,0,0],"Missed":[],"Game":36,"Cred":"xxxxxxxxxxxxxxxxxxxxxxx"}`

A guess from the client
`{"Cmd":"P1T","Play":"r","Gid":36,"Auth":"xxxxxxxxxxxxxxxxxxxxxxx"}`

Response from the server
`{"Cmd":"P1T","Curr":[0,0,0,0,114,0,0,0,0,0,0],"Missed":[],"Game":36}`
