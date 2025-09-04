package main

import (
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	noise "github.com/libp2p/go-libp2p/p2p/security/noise"
	ws "github.com/libp2p/go-libp2p/p2p/transport/websocket"
)

func main() {
	log.Println("🚀 Starting libp2p node setup...")
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/443/ws"),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Transport(ws.New),
	)

	if err != nil {
		log.Fatalf(" Failed to create libp2p host: %v", err)
	}

	_, err = relayv2.New(h)
	if err != nil {
		log.Fatalf("failed to enbale relay: %v", err)
	}

	log.Println("libp2p host created successfully.")
	log.Printf("🆔 Peer ID: %s\n", h.ID().String())

	log.Println("📡 Listening on addresses:")

	fmt.Println("Peer ID:", h.ID().String())
	for _, addr := range h.Addrs() {
		log.Printf(" - %s\n", addr)
	}

	log.Println("⏳ Node is running... Press Ctrl+C to exit.")

	select {} // or either use for{} so dont exit the program
}
