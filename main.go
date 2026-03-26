// relay/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

func main() {
	// Создаем хост для релея
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/4001"),
		libp2p.EnableRelayService(),
	)
	if err != nil {
		log.Fatalf("Failed to create relay host: %v", err)
	}
	defer h.Close()

	// Включаем сервис релея
	_, err = relay.New(h)
	if err != nil {
		log.Fatalf("Failed to create relay service: %v", err)
	}

	log.Printf("Relay node started with ID: %s", h.ID())
	log.Printf("Relay node addresses:")
	for _, addr := range h.Addrs() {
		log.Printf("  %s/p2p/%s", addr, h.ID())
	}

	// Ждем сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down relay...")
}
