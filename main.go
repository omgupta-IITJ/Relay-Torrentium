package main

import (
	"fmt"
	"log"
	"os"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	noise "github.com/libp2p/go-libp2p/p2p/security/noise"
	ws "github.com/libp2p/go-libp2p/p2p/transport/websocket"
	ma "github.com/multiformats/go-multiaddr"
)

type NotifyBundle struct{}

// if relay start listening on multiaddr
func (nb *NotifyBundle) Listen(n network.Network, a ma.Multiaddr) {
	log.Printf("[notifiee] Listen: %s\n", a)
}

// if relay stops listening on multiaddr
func (nb *NotifyBundle) ListenClose(n network.Network, a ma.Multiaddr) {
	log.Printf("[notifiee] ListenClose: %s\n", a)
}

// if someone connects on this multiaddr
func (nb *NotifyBundle) Connected(n network.Network, c network.Conn) {
	log.Printf("[notifiee] Connected: %s <-> %s  peer=%s\n",
		c.LocalMultiaddr(), c.RemoteMultiaddr(), c.RemotePeer().String())
}

func (nb *NotifyBundle) Disconnected(n network.Network, c network.Conn) {
	log.Printf("[notifiee] Disconnected: %s <-> %s  peer=%s\n",
		c.LocalMultiaddr(), c.RemoteMultiaddr(), c.RemotePeer().String())
}

// OpenedStream is called when a stream is opened on a connection
func (nb *NotifyBundle) OpenedStream(net network.Network, stream network.Stream) {
	log.Printf("[notifiee] OpenedStream: from=%s to=%s protocol=%s\n",
		stream.Conn().LocalPeer().String(), stream.Conn().RemotePeer().String(), stream.Protocol())
}

// ClosedStream is called when a stream is closed
func (nb *NotifyBundle) ClosedStream(net network.Network, stream network.Stream) {
	log.Printf("[notifiee] ClosedStream: from=%s to=%s protocol=%s\n",
		stream.Conn().LocalPeer().String(), stream.Conn().RemotePeer().String(), stream.Protocol())
}

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

	// Register the notifiee so we get connection/stream events
	h.Network().Notify(&NotifyBundle{})

	log.Println("✅ Relay started successfully")
	log.Printf("🆔 Peer ID: %s", h.ID())

	for _, addr := range h.Addrs() {
		log.Printf("📡 Listening on: %s/p2p/%s", addr, h.ID())
	}

	select {}

}
