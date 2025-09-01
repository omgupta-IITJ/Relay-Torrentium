package main

import (
    libp2p "github.com/libp2p/go-libp2p"
    circuit "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
    noise "github.com/libp2p/go-libp2p/p2p/security/noise"
    ws "github.com/libp2p/go-ws-transport"
	"log"
	"fmt"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
    h, err := libp2p.New(
        libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/443/wss"),
        libp2p.Security(noise.ID, noise.New),
        libp2p.Transport(ws.New),
        libp2p.EnableRelay(),
    )
    if err != nil {
        log.Fatal(err)
    }

    _, err = circuit.New(h)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Peer ID:", h.ID().String())
    for _, addr := range h.Addrs() {
        fmt.Println(addr.Encapsulate(ma.StringCast("/p2p/" + h.ID().String())))
    }

    select {}
}
