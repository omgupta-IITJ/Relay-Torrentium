package main

import(
	"github.com/gorilla/websocket"
	"fmt"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

var peers = make(map[string]*websocket.Conn)

func handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var peerID string
	// First message from peer is its ID
	err = conn.ReadJSON(&peerID)
	if err != nil {
		log.Println("PeerID error:", err)
		return
	}
	peers[peerID] = conn
	fmt.Println("Peer connected:", peerID)

	for {
		var msg map[string]string
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// relay message to target peer
		target := msg["to"]
		if targetConn, ok := peers[target]; ok {
			targetConn.WriteJSON(msg)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConn)
	log.Println("Relay server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}