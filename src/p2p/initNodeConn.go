package p2p

import (
	"context"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
	"github.com/quic-go/quic-go"
)

func _InitNodeConn(ctx context.Context, cancel context.CancelCauseFunc, isIncommingConnection bool, config *NodeP2PConnectionConfig, conn quic.Connection) (*NodeP2PConnection, error) {
	// Log
	if isIncommingConnection {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "An attempt is made to initialize the connection %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
	} else {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "An attempt is made to initialize the connection %s -> %s", conn.LocalAddr(), conn.RemoteAddr())
	}

	// Die Control Streams werden geöffnet
	controlStream, err := _TryOpenP2PConnectionControlStream(isIncommingConnection, conn, config, nil, nil, ctx)
	if err != nil {
		return nil, err
	}

	// Log
	if isIncommingConnection {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Control Streams opened %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
	} else {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Control Streams opened %s -> %s", conn.LocalAddr(), conn.RemoteAddr())
	}

	// Die Package Traffic Strams werden geöffnet
	_, err = _TryOpenQuicBidirectionalStream(isIncommingConnection, conn, []byte("h"), ctx)
	if err != nil {
		return nil, err
	}

	// Log
	if isIncommingConnection {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Package Traffic Streams opened %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
	} else {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Package Traffic Streams opened %s -> %s", conn.LocalAddr(), conn.RemoteAddr())
	}

	// Die Verbindung wird erzeugt
	nodeConn := &NodeP2PConnection{
		conn:                  conn,
		config:                config,
		context:               ctx,
		contextCancel:         cancel,
		isIncommingConnection: isIncommingConnection,
		controlStream:         controlStream,
	}

	return nodeConn, nil
}
