package p2p

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
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

// determineAgreedVersion ermittelt eine gemeinsame Version zwischen Client und Server.
func determineAgreedVersion(clientVersions, serverVersions []string) (string, error) {
	versionSet := make(map[string]struct{})
	for _, v := range serverVersions {
		versionSet[v] = struct{}{}
	}

	for _, v := range clientVersions {
		if _, exists := versionSet[v]; exists {
			return v, nil // Erste gemeinsame Version gefunden
		}
	}

	return "", fmt.Errorf("keine gemeinsame Version gefunden")
}

// IsUnixSocket prüft, ob die Verbindung ein Unix-Socket ist.
func isUnixSocket(conn net.Conn) bool {
	if _, ok := conn.(*net.UnixConn); ok {
		return true
	}
	return false
}

// IsTLS prüft, ob die Verbindung TLS ist.
func isTLS(conn net.Conn) bool {
	if _, ok := conn.(*tls.Conn); ok {
		return true
	}
	return false
}

// IsTCP prüft, ob die Verbindung TCP ist.
func isTCP(conn net.Conn) bool {
	if _, ok := conn.(*net.TCPConn); ok {
		return true
	}
	return false
}
