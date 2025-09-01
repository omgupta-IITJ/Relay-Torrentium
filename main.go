package main

import (
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	// Get PORT from Render env
	port := os.Getenv("PORT")
	if port == "" {
		port = "10010" // local dev fallback
	}

	// IMPORTANT: only use /tcp/, not /wss/, because Render will handle TLS
	addrStr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", port)

	listenAddr, err := ma.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}

	_, err = libp2p.New(
		libp2p.ListenAddrs(listenAddr),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Libp2p host started on:", listenAddr)
	select {}
}
