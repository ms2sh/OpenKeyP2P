package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ms2sh/OpenKeyP2P/src/p2p"
)

func main() {
	// Erstelle einen TLS-Listener auf Port 4433
	listener, err := net.Listen("tcp", ":4433")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server läuft auf Port 4433...")

	// Akzeptiere eingehende Verbindungen
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("Failed to accept connection: %v", err)
	}

	upgrconn := p2p.UpgradeConn(conn)

	data, err := upgrconn.Read()
	if err != nil {
		panic(err)
	}

	fmt.Println(len(data))
}
