package main

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	noise "github.com/libp2p/go-libp2p/p2p/security/noise"
	ws "github.com/libp2p/go-libp2p/p2p/transport/websocket"
	ma "github.com/multiformats/go-multiaddr"
	// priv "privkey/privkey"
)

func loadandgenratepriv() crypto.PrivKey {
	// generating a new RSA private key of 2048 bits
	const pemFile = "ed25519_key.pem"
	if _, err := os.Stat(pemFile); err == nil {
		data, err := os.ReadFile(pemFile)
		if err != nil {
			log.Println("couldnt read the file", err)
		}
		block, _ := pem.Decode(data)
		if block == nil {
			log.Printf("pemfile not decoded")
		}
		privkey, err := crypto.UnmarshalEd25519PrivateKey(block.Bytes)
		if err != nil {
			log.Fatalf("Failed to parse RSA private key: %v", err)
		}

		return privkey
	}
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		log.Printf("error in creating rsa private key: %v", err)
	}
	// encoding private key to PEM format
	// header label of PEM file that will appear in file
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		log.Fatalf("Failed to marshal private key: %v", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "ED25519 PRIVATE KEY",
		Bytes: privBytes,
	}

	privateKeyFile, err := os.Create(pemFile)
	if err != nil {
		log.Println("Error creating private key file:", err)
	}
	defer privateKeyFile.Close()

	pem.Encode(privateKeyFile, privateKeyPEM)
	privateKeyFile.Close()
	return priv
}

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
	log.Println(" Starting libp2p relay...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("unable to load env file", err)
	}

	publicDNS := os.Getenv("RENDER_EXTERNAL_URL")

	port := os.Getenv("PORT")
	if port == "" {
		port = "443"
	}

	priv := loadandgenratepriv()

	// using DER(DISTINGUISHED ENCODING RULES) FOR RSA->crypto
	// Convert to PKCS1 DER bytes

	// Now turn it into libp2p's crypto.PrivKey

	// using a layer of noise protocol for secured connection between peer and relay
	h, err := libp2p.New(
		libp2p.Identity(priv),
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
		log.Printf("Failed to create libp2p host: %v", err)
	}

	// relay as a libp2p host
	_, err = relayv2.New(h)
	if err != nil {
		log.Printf(" Failed to enable relay v2: %v", err)
	}

	// Register the notifiee so we get connection/stream events
	h.Network().Notify(&NotifyBundle{})

	log.Println("Relay started successfully")
	log.Printf(" Peer ID: %s", h.ID())
	fmt.Println("public domain name system is:", publicDNS)

	for _, addr := range h.Addrs() {
		log.Printf(" Listening on: %s/p2p/%s", addr, h.ID())
	}

	select {}

}
