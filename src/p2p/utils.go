package p2p

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
)

// isTemporaryError überprüft, ob ein Fehler vorübergehend ist und einen erneuten Versuch rechtfertigt.
func isTemporaryError(err error) bool {
	// Beispielhafte Implementierung: Prüft, ob der Fehler temporär ist
	// Dies kann erweitert werden, um spezifische Fehlertypen zu erkennen
	var tempErr interface{ Temporary() bool }
	if errors.As(err, &tempErr) {
		return tempErr.Temporary()
	}
	return false
}

// randomUint32 generates a random uint32 number using cryptographically secure random bytes.
func randomUint32() (uint32, error) {
	b := make([]byte, 4)
	// Read 4 cryptographically secure random bytes
	if _, err := rand.Read(b); err != nil {
		return 0, err
	}
	// Convert the 4 bytes into a uint32 (Big-Endian)
	return binary.BigEndian.Uint32(b), nil
}
