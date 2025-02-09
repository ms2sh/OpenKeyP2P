package p2p

import (
	"fmt"
	"net"
	"sync"

	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

// UpgradeConn creates and returns a new instance of Conn, wrapping the provided net.Conn.
func UpgradeConn(conn net.Conn, role Role, supportedVersions []string) (*FWConn, error) {
	var maxFrameSize uint64
	var connType _ConnType
	var dterr error
	switch {
	case isTCP(conn):
		maxFrameSize = (1 * 1024 * 1024) / 32
		connType = TCP
	case isTLS(conn):
		maxFrameSize = (1 * 1024 * 1024) / 32
		connType = TLS
	case isUnixSocket(conn):
		connType = UnixSocket
		maxFrameSize = (1 * 1024 * 1024) / 32
		conn, dterr = upgradeConnToUnixSocketAckConn(conn)
	default:
		dterr = fmt.Errorf("not supported connection type")
	}

	if dterr != nil {
		return nil, dterr
	}

	c := &FWConn{conn: conn, connType: connType, maxFrameSize: maxFrameSize}
	c.cond = sync.NewCond(&c.mu) // Initialize the condition variable

	versionInfo := VersionInfo{}

	if role == Client {
		// Client sends supported versions
		versionInfo.SupportedVersions = supportedVersions
		if err := c.sendVersionInfo(&versionInfo); err != nil {
			return nil, fmt.Errorf("error sending version information: %w", err)
		}

		// Client receives the agreed version from the server
		if err := c.receiveServerResponse(&versionInfo); err != nil {
			return nil, fmt.Errorf("error receiving server response: %w", err)
		}

		logging.LogInfo("Agreed version: %s", versionInfo.AgreedVersion)
	} else if role == Server {
		// Server receives supported versions from the client
		var clientVersionInfo VersionInfo
		if err := c.receiveClientVersionInfo(&clientVersionInfo); err != nil {
			return nil, fmt.Errorf("error receiving client version information: %w", err)
		}

		// Server determines the common version
		agreedVersion, err := determineAgreedVersion(clientVersionInfo.SupportedVersions, supportedVersions)
		if err != nil {
			return nil, fmt.Errorf("error determining agreed version: %w", err)
		}

		// Server sends the agreed version to the client
		response := VersionInfo{
			AgreedVersion: agreedVersion,
		}
		if err := c.sendVersionInfo(&response); err != nil {
			return nil, fmt.Errorf("error sending server response: %w", err)
		}

		versionInfo.AgreedVersion = agreedVersion
		logging.LogInfo("Agreed version with client: %s", agreedVersion)
	} else {
		return nil, fmt.Errorf("unknown role: %v", role)
	}

	return c, nil
}
