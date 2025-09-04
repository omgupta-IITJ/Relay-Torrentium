package main

import (
	"fmt"
	"log"
	"os"

	libp2p "github.com/libp2p/go-libp2p"
	relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	noise "github.com/libp2p/go-libp2p/p2p/security/noise"
	ws "github.com/libp2p/go-libp2p/p2p/transport/websocket"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	log.Println("🚀 Starting libp2p relay...")

	// Render provides external URL like: https://your-app.onrender.com
	publicDNS := os.Getenv("RENDER_EXTERNAL_URL")
	// publicDNS = strings.TrimPrefix(publicDNS, "https://")
	// publicDNS = strings.TrimSuffix(publicDNS, "/")

	port := os.Getenv("PORT")
	if port == "" {
		port = "443"
	}

	// Listen inside container
	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", port)
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(listenAddr),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Transport(ws.New),
		libp2p.AddrsFactory(func(addrs []ma.Multiaddr) []ma.Multiaddr {
			maddr, _ := ma.NewMultiaddr(
				fmt.Sprintf("/dns4/%s/tcp/%s/wss", publicDNS, port),
			)
			return []ma.Multiaddr{maddr}
		}),
	)
	if err != nil {
		log.Fatalf("❌ Failed to create libp2p host: %v", err)
	}

	// ✅ Enable relay v2 service
	_, err = relayv2.New(h)
	if err != nil {
		log.Fatalf("❌ Failed to enable relay v2: %v", err)
	}

	log.Println("✅ Relay started successfully")
	log.Printf("🆔 Peer ID: %s", h.ID())

	for _, addr := range h.Addrs() {
		log.Printf("📡 Listening on: %s/p2p/%s", addr, h.ID())
	}

	select {}

}
