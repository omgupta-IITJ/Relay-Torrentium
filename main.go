package main

import (
	"context"
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	identify "github.com/libp2p/go-libp2p/p2p/protocol/identify"
	mplex "github.com/libp2p/go-libp2p/p2p/muxer/mplex"
	noise "github.com/libp2p/go-libp2p/p2p/security/noise"
	ws "github.com/libp2p/go-ws-transport"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()

	// Create a libp2p host
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/10000/ws",   // WebSocket for local/dev
			"/ip4/0.0.0.0/tcp/443/wss",    // Secure WebSocket (Render maps HTTPS→443)
		),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
		libp2p.Transport(ws.New),
		libp2p.EnableRelay(), // enable relay support
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start relay service
	_, err = circuit.New(h)
	if err != nil {
		log.Fatal(err)
	}

	// Start identify service (important for peers to fetch addresses)
	identify.NewIDService(h)

	fmt.Println("✅ Relay node started")
	fmt.Println("Peer ID:", h.ID().Pretty())
	fmt.Println("Listening addresses:")
	for _, addr := range h.Addrs() {
		fmt.Println(addr.Encapsulate(ma.StringCast("/p2p/" + h.ID().Pretty())))
	}

	<-ctx.Done()
}
