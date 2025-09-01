package main

import (
	"context"
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()

	// Start a libp2p host
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			// WebSocket secure transport on 0.0.0.0:443
			"/ip4/0.0.0.0/tcp/10000/wss",
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Enable relay v2 on this host
	_, err = relayv2.New(h)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Relay node started 🚀")
	fmt.Println("Peer ID:", h.ID())

	// Print out relay’s multiaddresses
	for _, addr := range h.Addrs() {
		fmt.Println(" - ", addr.Encapsulate(ma.StringCast(fmt.Sprintf("/p2p/%s", h.ID()))))
	}

	<-ctx.Done()
}
