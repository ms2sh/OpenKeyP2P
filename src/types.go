package openkeyp2p

import (
	"crypto/ed25519"
)

// Gibt den Prefix einer Node Adresse an
type OpenKeyP2PPrefix string

// Gibt die Version eines Nodes an
type OpenKeyP2PVesion uint64

// Stellt den Verwendeten Addresstypen dat
type OpenKeyP2PKeyType uint8

// Stellt den ED25519 PublicKey dar
type OpenKeyP2PPublicKey ed25519.PublicKey

// Stellt ein Base32 Alphabet dar
type Base32Alphabet string

// Stellt einen Verfügbaren Hash Algroytmus dar
type HashAlgorithm uint16

// Stellt einen Hash dar
type HashSlice []byte

// LogLevel
type LogLevel uint8

type LocalListenerAddress string
