package p2p

import (
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

func _StartWriterRoutinesForNodeConn(conn *NodeP2PConnection, wg *sync.WaitGroup) error {
	// Log
	localEndpointStr := getLocalIPAndHostFromConn(conn.conn)
	remoteEndpointStr := getRemoteIPAndHostFromConn(conn.conn)
	if conn.isIncommingConnection {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Writer started %s -> %s", remoteEndpointStr, localEndpointStr)
	} else {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Writer started %s -> %s", localEndpointStr, remoteEndpointStr)
	}

	// Es wird Signalisiert das die Writer Routinen ausgeführt werden
	wg.Done()

	return nil
}
