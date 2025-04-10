package p2p

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"slices"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
	"github.com/quic-go/quic-go"
)

func _InitQUICNodeConn(localhostNetworkInterface *net.Interface, ctx context.Context, cancel context.CancelCauseFunc, isIncommingConnection bool, config NodeP2PConnectionConfig, conn quic.Connection) (*NodeP2PConnection, error) {
	// Der Lokale EP sowie der Remote EP wird abgerufen
	localEndpointStr := getLocalIPFromConn(conn)
	remoteEndpointStr := getRemoteIPAndHostFromConn(conn)

	// LOG
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "An attempt is made to initialize the connection %s -> %s", localEndpointStr, remoteEndpointStr)

	// Die Control Streams werden geöffnet
	controlStream, err := _TryOpenP2PConnectionControlStream(localhostNetworkInterface, isIncommingConnection, conn, config, NodeP2PSocketAddress(localEndpointStr), NodeP2PSocketAddress(remoteEndpointStr), ctx, cancel)
	if err != nil {
		return nil, err
	}

	// Es wird geprüft ob die Version unterstützt wird (LOKAL)
	localAcceptRemoteVersion := slices.Contains(openkeyp2p.SUPPORTED_VERSION, controlStream.GetDestinationVersion())
	if !localAcceptRemoteVersion {
		return nil, fmt.Errorf("the version of peer dosent accepted")
	}

	// Es wird geprüft ob die Gegenseite die Lokale Version Akzeptiert
	remoteAcceptLocalVersion := slices.Contains(controlStream.destPeerHelloPacket.SupportedVersions, openkeyp2p.Version)
	if !remoteAcceptLocalVersion {
		return nil, fmt.Errorf("remote peer dosent accept the local peer version")
	}

	// Die Gemeinsam Unterstützen Funktionen werden ermittelt
	connectionConfig := _DeterminesCommonConfig(controlStream.destPeerHelloPacket.NodeConfigOptions, config)

	// Log
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Control Streams opened %s -> %s", localEndpointStr, remoteEndpointStr)

	// Die Package Traffic Strams werden geöffnet
	trafficStream, err := _TryOpenP2PConnectionTrafficStream(isIncommingConnection, conn, NodeP2PSocketAddress(localEndpointStr), NodeP2PSocketAddress(remoteEndpointStr), ctx, cancel)
	if err != nil {
		return nil, err
	}

	// Log
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Package Traffic Streams opened %s -> %s", localEndpointStr, remoteEndpointStr)

	// Die Verbindung wird erzeugt
	nodeConn := &NodeP2PConnection{
		conn:                    conn,
		config:                  connectionConfig,
		ctx:                     ctx,
		localKeepalivePacketIds: new(sync.Map),
		contextCancel:           cancel,
		isIncommingConnection:   isIncommingConnection,
		controlStream:           controlStream,
		packageTrafficStream:    trafficStream,
		writerControlBuffer:     openkeyp2p.NewThreadSafeContextBuffer(ctx),
		keepaliveTime:           12 * time.Second,
		localSocketAddress:      NodeP2PSocketAddress(localEndpointStr),
		remoteSocketAddress:     NodeP2PSocketAddress(remoteEndpointStr),
	}

	// Die Verbindung wird zurückgegeben
	return nodeConn, nil
}
