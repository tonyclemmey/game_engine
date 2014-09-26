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
configured properly, it will work by `$GOROOT/bin/engine`. The process listens
on port 3000 by default, but this can be changed in the code.

Configuration
--------------

There is no configuration capabilities at this time. 

Protocol
--------------

Yet to be documented.
