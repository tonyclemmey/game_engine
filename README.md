Game Server Engine
==============

This is the game server that is implemented at a basic level. It has wrapper
capabilities that will be used to cater to various other plugins. At this time,
it is not a plugin framework. Rather, there is one game that is integrated
in to the code base. 

Language
--------------

The server is implemented in Go using HTTP and JSON. To build, checkout the
code, navigate to the 'engine' and `go install`. If your Go environment is
configured properly, it will work by `$GOPATH/bin/engine`. The process listens
on port 3000 by default, but this can be changed in the code.

Games
--------------

The only game written at this time is Hangman. A Python cli and web clients
are provided. Right now, the game does not end. You may hang your character,
but you can continue on to solve the word. The Macmillan API is used to
provide a hint. A running example is at: 

http://richmond.cookgetsitdone.com/hangman.html

Configuration
--------------

There is no configuration capabilities at this time. 

Protocol
--------------

Yet to be documented.
