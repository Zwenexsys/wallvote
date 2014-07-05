// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/gorilla/websocket"
    "github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
    "encoding/json"
    "fmt"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

    // Name of the connection
    name string
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
        log.Printf("READ: %s", message)

        // Parse and handle message
        ccmd := Command{}
        if cerr := json.Unmarshal(message, &ccmd); cerr != nil { // Upon JSON Error send error 
            scmd := Command{Cmd: "message", Message: "Invalid request."}
            jscmd, _ := json.Marshal(scmd)
            h.self <- struct{message []byte; src *connection}{jscmd, c}
        }else{
            switch ccmd.Cmd {
                case "add_card":
                    cdata := ccmd.Data.(map[string]interface{})
                    ncard := wall.AddCard(cdata["Text"].(string), cdata["Name"].(string))
                    scmd := Command{Cmd: "card_added", Message: "You added new card.", Data: ncard}
                    scmd1 := Command{Cmd: "card_added", Message: fmt.Sprintf("%s added new card.", ncard.Name), Data: ncard}
                    jscmd, _ := json.Marshal(scmd)
                    jscmd1, _ := json.Marshal(scmd1)
                    h.others <- struct{message []byte; src *connection}{jscmd1, c}
                    h.self <- struct{message []byte; src *connection}{jscmd, c}
                case "plus_card":
                    cdata := ccmd.Data.(map[string]interface{})
                    if wall.IsVoted(cdata["Id"].(string), cdata["Name"].(string)){
                        if wall.UnplusCard(cdata["Id"].(string), cdata["Name"].(string)) {
                            card := wall.Cards[cdata["Id"].(string)]
                            scmd := Command{Cmd: "card_unplused", Message: "Card unplused.", Data: card}
                            jscmd, _ := json.Marshal(scmd)
                            h.broadcast <- struct{message []byte; src *connection}{jscmd, c}
                        }else{
                            break
                            log.Println("Cannot plus card.")
                        }
                    }else{
                        if wall.PlusCard(cdata["Id"].(string), cdata["Name"].(string)) {
                            card := wall.Cards[cdata["Id"].(string)]
                            scmd := Command{Cmd: "card_plused", Message: "Card plused.", Data: card}
                            jscmd, _ := json.Marshal(scmd)
                            h.broadcast <- struct{message []byte; src *connection}{jscmd, c}
                        }else{
                            break
                            log.Println("Cannot plus card.")
                        }
                    }
                case "move_card":
                    cdata := ccmd.Data.(map[string]interface{})
                    if card, status := wall.MoveCard(cdata["Id"].(string), cdata["Name"].(string), cdata["X"].(string), cdata["Y"].(string)); status == true{
                        scmd := Command{Cmd: "card_moved", Message: "Card moved.", Data: card}
                        jscmd, _ := json.Marshal(scmd)
                        h.others <- struct{message []byte; src *connection}{jscmd, c}
                    }else{
                        log.Println("Cannot move card.")
                        break
                    }
                default:
                    scmd := Command{Cmd: "message", Message: "Server not understood!"}
                    jscmd, _ := json.Marshal(scmd)
                    h.self <- struct{message []byte; src *connection}{jscmd, c}
            }
            //h.broadcast <- struct{message []byte; src *connection}{jscmd, c}
        }
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serverWs handles webocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    log.Printf("New connection : %s", vars["name"])

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
    c := &connection{send: make(chan []byte, 256), ws: ws, name: vars["name"]}
	h.register <- c
	go c.writePump()
	c.readPump()
}
