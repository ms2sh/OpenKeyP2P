package p2p

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/quic-go/quic-go"
)

func openNodeConnection(nodeUri string, isIncommingConnection bool, tlsConfig *tls.Config, config *NodeP2PConnectionConfig) (*NodeP2PConnection, error) {
	// Jeder Client bekommt seinen eigenen Kontext
	ctx, cancel := context.WithCancelCause(context.Background())

	// Die Quic Verbindung wird aufgebaut
	conn, err := quic.DialAddr(ctx, nodeUri, tlsConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("openNodeConnection: " + err.Error())
	}

	// Das Basis Paket wird erstellt
	nodeConn := &NodeP2PConnection{
		conn:                  conn,
		config:                config,
		context:               ctx,
		contextCancel:         cancel,
		isIncommingConnection: isIncommingConnection,
		controlStream:         nil,
	}

	// Die Stream Acceptor werden gestartet
	if err := initControlStreams(nodeConn); err != nil {
		err = errors.New(fmt.Sprintf("openNodeConnection: %s", err.Error()))
		cancel(err)
		return nil, err
	}

	// Die Routing Request Streams werden geöffnet
	if err := initRoutingStreams(nodeConn); err != nil {
		err = errors.New(fmt.Sprintf("openNodeConnection: %s", err.Error()))
		cancel(err)
		return nil, err
	}

	// Die Package Traffic Streams werden geöffnet
	if err := initPackageTrafficStreams(nodeConn); err != nil {
		err = errors.New(fmt.Sprintf("openNodeConnection: %s", err.Error()))
		cancel(err)
		return nil, err
	}

	// Die Unreliable Datagrame Sitzungs Funktionen werden erstellt
	initUnreliableDatagramsHandle(nodeConn)

	return nodeConn, nil
}
