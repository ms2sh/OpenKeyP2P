package p2p

import (
	"crypto/tls"
	"fmt"
	"net"
)

// Conn wraps a net.Conn with a Mutex and a Condition to synchronize read/write operations.
type TlsConn struct {
	conn *Conn
}

// UpgradeConn creates and returns a new instance of Conn, wrapping the provided net.Conn.
func UpgradeTlsConn(conn net.Conn) (*TlsConn, error) {
	// Es wird versucht die Verbindug zu einer TLS Verbindung zu Upgraden
	tlsConn, isTlsConn := conn.(*tls.Conn)
	if !isTlsConn {
		return nil, fmt.Errorf("is no tls connection")
	}

	// TLS-Verbindungsstatus erhalten (Handschlag ist bereits durchgeführt)
	state := tlsConn.ConnectionState()

	// Prüfen, ob es Client-Zertifikate gibt
	if len(state.PeerCertificates) > 0 {
		fmt.Println("Client-Zertifikat:")
		fmt.Println("  Subject:", state.PeerCertificates[0].Subject)
		fmt.Println("  Issuer:", state.PeerCertificates[0].Issuer)
	}

	// Zugriff auf die verwendeten CA-Zertifikate
	fmt.Println("Verwendete CA-Zertifikate:")
	for _, cert := range state.VerifiedChains {
		for _, caCert := range cert {
			fmt.Println("  CA Issuer:", caCert.Issuer)
		}
	}

	return &TlsConn{conn: UpgradeConn(tlsConn)}, nil
}

// Write is responsible for writing data to the connection in fragments.
// It ensures that only one operation can write or read at a time.
func (c *TlsConn) Write(b []byte) error {
	return c.conn.Write(b)
}

// Read reads data from the connection in fragments.
// It ensures that only one operation can write or read at a time.
func (c *TlsConn) Read() ([]byte, error) {
	return c.conn.Read()
}

// Close closes the connection, releasing any resources.
func (c *TlsConn) Close() error {
	return c.conn.Close()
}
