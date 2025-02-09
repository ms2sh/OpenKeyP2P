package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/ms2sh/OpenKeyP2P/src/logging"
	"github.com/ms2sh/OpenKeyP2P/src/p2p"
)

// Generate1MBString erzeugt einen String mit einer Größe von 1 MB.
func Generate1MBString() string {
	// Der Buchstabe 'a' hat eine Größe von 1 Byte in UTF-8.
	return strings.Repeat("a", 5000)
}

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:4433")
	if err != nil {
		fmt.Println(err)
		return
	}
	upgradedConn, err := p2p.UpgradeConn(conn, p2p.Client, []string{"1.0"})
	if err != nil {
		logging.LogError("Fehler beim Upgrade der Verbindung: %v", err)
		return
	}

	err = upgradedConn.Write([]byte(Generate1MBString()))
	if err != nil {
		panic(err)
	}

}
