package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
		// io.WriteString(w, "hello, world!\n")
	})

	http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
			return
		}
		// defer conn.Close()

		go func(connection *websocket.Conn) {

			for {
				_, message, err := connection.ReadMessage()

				if err != nil {
					log.Panicln("ReadMessage broken: ", err)
				}

				log.Printf("Message that recieve: %s", message)

				// _ = conn.WriteMessage(messageType, message)

				defer connection.Close()
			}

		}(conn)

		// for {
		// 	mt, message, err := conn.ReadMessage()
		// 	if err != nil {
		// 		log.Println("read:", err)
		// 		break
		// 	}
		// 	log.Printf("recv: %s", message)
		//
		// 	// test := []byte("test")
		//
		// 	err = conn.WriteMessage(mt, message)
		// 	if err != nil {
		// 		log.Println("write:", err)
		// 		break
		// 	}
		//
		// }

	})

	http.ListenAndServe("localhost:3000", nil)
}
