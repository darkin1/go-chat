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

		chMessageType := make(chan int)
		chMessage := make(chan []byte)

		go func(connection *websocket.Conn, chMessageType chan int, chMessage chan []byte) {

			for {
				messageType, message, err := connection.ReadMessage()

				if err != nil {
					log.Panicln("ReadMessage broken: ", err)
				}

				log.Printf("Message that recieve: %s", message)

				// _ = conn.WriteMessage(messageType, message)

				chMessageType <- messageType
				chMessage <- message

				defer connection.Close()
			}

		}(conn, chMessageType, chMessage)

		go func(connection *websocket.Conn, chMessageType chan int, chMessage chan []byte) {

			for {
				messageType := <-chMessageType
				message := <-chMessage

				messagex := append([]byte("sent: "), message...)

				err := conn.WriteMessage(messageType, messagex)

				if err != nil {
					log.Panicln("WriteMessage broken: ", err)
				}
			}
		}(conn, chMessageType, chMessage)

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
