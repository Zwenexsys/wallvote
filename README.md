# WallVote - Simple Online Realtime Collaboration Platform Example

This is created for demonstraration at Google I/O Extended - Yangon 2014 on 5th July.

This application shows how to use the
[router](https://github.com/gorilla/mux) package,
[websocket](https://github.com/gorilla/websocket) package and
[jQuery](http://jquery.com) to implement a simple online realtime collaboration platform, shared Wall for Voting.

## Main Idea

To create a paperless conference wall, accessible by everyone in the room on the same network to share and post their idea for voting.

## Running the example

The example requires a working Go development environment. The [Getting
Started](http://golang.org/doc/install) page describes how to install the
development environment.

Once you have Go up and running, you can download, build and run the example
using the following commands.

    $ go get github.com/gorilla/mux
    $ go get github.com/gorilla/websocket
    $ cd `go list -f '{{.Dir}}' github.com/Zwenexsys/wallvote`
    $ go run *.go

## Base Code

The implemention is the extension of the code given in this [Chat Example](http://github.com/gorilla/chat)
