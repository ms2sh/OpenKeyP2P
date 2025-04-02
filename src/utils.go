package openkeyp2p

import (
	"encoding/binary"
	"fmt"
)

func Uint64ToBytesLE(n uint64) []byte {
	bytes := make([]byte, 8) // uint64 = 8 Bytes
	binary.LittleEndian.PutUint64(bytes, n)
	return bytes
}

func BytesToUint64LE(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func ParseVersion(num OpenKeyP2PVesion) string {
	// Beta-Status: die letzten 2 Ziffern
	beta := num % 100

	// Build: die 6 Ziffern davor
	build := (num / 100) % 1000000

	// Release: die 4 Ziffern davor
	release := (num / 100000000) % 10000

	// Version: alles davor
	version := num / 10000000000

	// Formatierte Ausgabe
	versionStr := fmt.Sprintf("%d.%d.%d", version, release, build)

	// Falls es eine Beta-Version ist
	if beta > 0 {
		versionStr = fmt.Sprintf("%s-Beta %d", versionStr, beta)
	}

	return versionStr
}

func CalculateQUICPayloadSize(mtu int, isIPv6 bool) int {
	const udpHeaderSize = 8   // UDP-Header ist immer 8 Bytes
	const quicHeaderSize = 15 // Durchschnittliche QUIC-Header-Größe
	const ipv4HeaderSize = 20 // IPv4-Header-Größe
	const ipv6HeaderSize = 40 // IPv6-Header-Größe

	// Wähle die richtige IP-Header-Größe
	ipHeaderSize := ipv4HeaderSize
	if isIPv6 {
		ipHeaderSize = ipv6HeaderSize
	}

	// Berechnung der nutzbaren QUIC-Payload
	quicPayloadSize := mtu - ipHeaderSize - udpHeaderSize - quicHeaderSize

	// Sicherheitspuffer: Viele Netzwerke begrenzen die effektive MTU weiter (z. B. VPN, Tunnel)
	if quicPayloadSize > 1350 {
		quicPayloadSize = 1350
	}

	return quicPayloadSize
}
