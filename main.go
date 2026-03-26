// relay.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

const identifyProtocol = "/ipfs/id/1.0.0"
const publicIPProtocol = "/publicIPProtocol/1.0.0"

type RelayNode struct {
	host host.Host
}

func main() {
	// Создаем relay хост
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/4001",
			"/ip4/0.0.0.0/udp/4001/quic-v1",
		),
		libp2p.EnableRelayService(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	// Включаем relay сервис
	_, err = relay.New(h)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("=== RELAY NODE ===")
	log.Printf("ID: %s", h.ID())
	log.Printf("Addresses:")
	for _, addr := range h.Addrs() {
		log.Printf("  %s/p2p/%s", addr, h.ID())
	}
	log.Println("=================")

	// Добавить в relay.go
	h.SetStreamHandler(publicIPProtocol, func(stream network.Stream) {
		defer stream.Close()

		// Получаем удаленный адрес (это и есть публичный адрес клиента)
		remoteAddr := stream.Conn().RemoteMultiaddr()
		publicAddr := fmt.Sprintf("%s/p2p/%s", remoteAddr, stream.Conn().RemotePeer())

		// Отправляем клиенту его публичный адрес
		stream.Write([]byte(publicAddr))

		log.Printf("Gave public address to %s: %s", stream.Conn().RemotePeer(), publicAddr)
	})

	// Ждем сигнала
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
