package main

import (
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

func main() {
	// ctx := context.Background()

	// Create a new libp2p host that acts as a relay
	host, err := libp2p.New(
		libp2p.EnableRelay(), // Enable relay on this host
	)
	if err != nil {
		log.Fatalf("Failed to create relay host: %v", err)
	}

	// Start the relay service
	_, err = relayv2.New(host)
	if err != nil {
		log.Fatalf("Failed to start relay service: %v", err)
	}

	fmt.Println("🚀 Relay server is running!")
	fmt.Println("Peer ID:", host.ID())
	fmt.Println("Multiaddresses:")
	for _, addr := range host.Addrs() {
		fmt.Printf(" - %s/p2p/%s\n", addr, host.ID())
	}

	select {} // keep running
}
