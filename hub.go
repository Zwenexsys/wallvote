// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import(
    "fmt"
    "encoding/json"
    "log"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
    // Registered connections.
    connections map[*connection]bool

    // Inbound messages from the connections.
    broadcast chan struct{message []byte; src *connection}
    others chan struct{message []byte; src *connection}

    // Inbound message from connection to self
    self chan struct{message []byte; src *connection}

    // Register requests from the connections.
    register chan *connection

    // Unregister requests from connections.
    unregister chan *connection
}

func (h *hub) connectionsExcept(conn connection) (map[*connection]bool) {
    nconns := h.connections
    delete(nconns, &conn)
    return nconns
}

var h = hub{
    broadcast:   make(chan struct{message []byte; src *connection}),
    others:   make(chan struct{message []byte; src *connection}),
    self:   make(chan struct{message []byte; src *connection}),
    register:    make(chan *connection),
    unregister:  make(chan *connection),
    connections: make(map[*connection]bool),
}

func (h *hub) run() {
    for {
        select {
        case c := <-h.register:
            h.connections[c] = true

            // Setup Command to Self
            scmd := Command{Cmd: "setup", Message: "You joined the Wall.", Data: wall}
            scmd_json, err1 := json.Marshal(&scmd)
            if err1 != nil {
                log.Println(err1)
            }else{
                log.Println(string(scmd_json))
            }

            // Loop self
            c.send <- []byte(scmd_json)

            // Broad cast to all
            bcmd := Command{Cmd: "message", Message: fmt.Sprintf("%s has joined the hub.", c.name)}
            bmessage, err2 := json.Marshal(&bcmd)
            if err2 != nil {
                log.Println(err2)
            }else{
                log.Println(string(bmessage))
            }
            for cn := range h.connections {
                if cn.name != c.name {
                    select {
                    case cn.send <- []byte(bmessage):
                    default:
                        close(cn.send)
                        delete(h.connections, cn)
                    }
                }
            }
        case c := <-h.unregister:
            if _, ok := h.connections[c]; ok {
                delete(h.connections, c)
                close(c.send)

                // Inform the rest
                bmessage := []byte(fmt.Sprintf("%s has left the hub.", c.name))
                for cn := range h.connections {
                    select {
                    case cn.send <- bmessage:
                    default:
                        close(cn.send)
                        delete(h.connections, cn)
                    }
                }
            }
        case b := <-h.self:
            log.Printf("Self: %s", b.src.name);
            b.src.send <- b.message
        case b := <-h.broadcast:
            for c := range h.connections {
                select {
                case c.send <- b.message:
                default:
                    close(c.send)
                    delete(h.connections, c)
                }
            }
        case b := <-h.others:
            for c := range h.connections {
                if c.name != b.src.name {
                    select {
                    case c.send <- b.message:
                    default:
                        close(c.send)
                        delete(h.connections, c)
                    }
                }
            }
        }
    }
}
