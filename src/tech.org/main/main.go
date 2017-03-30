package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"fmt"
)

const (
	ADDR string = ":8080"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
}

func serveWs(hub *Hoster, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	if err != nil {
		log.Println(err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}


func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))

	for {

		if ws == nil {
			break
		}
		// Read
		_, msg, err := ws.ReadMessage()

		if err != nil {
			log.Fatal(err)
		} else {

			var answer string

			switch string(msg) {
				case "hello":
				case "lol":
					answer = "kek"
					fallthrough
				default:
					answer = "Your message: " + string(msg)
			}

			ws.WriteMessage(websocket.TextMessage, []byte(answer))
		}
		fmt.Printf("%s\n", msg)
	}
}

func main() {
	hoster := newHoster()
	go hoster.run()
	http.HandleFunc("/", homeHandler) // set http request handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {  // websocket handler
		serveWs(hoster, w, r)
	})

	if err := http.ListenAndServe(ADDR, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}