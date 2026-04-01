// server.go
package main

import (
	"fmt"
	"net"
)

var clients = make(map[string]*net.UDPAddr)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":9999")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	fmt.Println("Server started on :9999")

	buf := make([]byte, 1024)

	for {
		n, clientAddr, _ := conn.ReadFromUDP(buf)
		id := string(buf[:n])

		fmt.Println("Client connected:", id, clientAddr.String())
		clients[id] = clientAddr

		// Если есть оба клиента — отправляем им адреса друг друга
		if len(clients) >= 2 {
			var a, b *net.UDPAddr
			for _, v := range clients {
				if a == nil {
					a = v
				} else {
					b = v
				}
			}

			conn.WriteToUDP([]byte(b.String()), a)
			conn.WriteToUDP([]byte(a.String()), b)

			fmt.Println("Exchanged addresses")
		}
	}
}
