package p2p

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/quic-go/quic-go"
)

func openNodeConnection(nodeUri string, isIncommingConnection bool, tlsConfig *tls.Config, config *NodeP2PConnectionConfig) (*NodeP2PConnection, error) {
	// Jeder Client bekommt seinen eigenen Kontext
	ctx, cancel := context.WithCancelCause(context.Background())

	// Die Quic Verbindung wird aufgebaut
	conn, err := quic.DialAddr(ctx, nodeUri, tlsConfig, nil)
	if err != nil {
		err = fmt.Errorf("openNodeConnection: %w", err)
		cancel(err)
		return nil, err
	}

	// Die Verbindung wird erzeugt
	nodeConn := &NodeP2PConnection{
		conn:                  conn,
		config:                config,
		context:               ctx,
		contextCancel:         cancel,
		isIncommingConnection: isIncommingConnection,
		controlStream:         nil,
	}

	// Der Control Stream wird vorbereitet
	controlStream, err := initControlStreams(conn, ctx, cancel)
	if err != nil {
		err = fmt.Errorf("openNodeConnection: %w", err)
		cancel(err)
		return nil, err
	}

	// Wenn es sich um eine eintreffende Verbindung handelt, wird der gegennseite das Server Propmt geschickt
	// Sollte es sich um eine ausgehende Verbindung, wird die
	if isIncommingConnection {
		if err = controlStream.IncommingSideInit(); err != nil {
			err = fmt.Errorf("openNodeConnection: %w", err)
			cancel(err)
			return nil, err
		}
	} else {
		if err = controlStream.OutgoingSideInit(); err != nil {
			err = fmt.Errorf("openNodeConnection: %w", err)
			cancel(err)
			return nil, err
		}
	}

	// Der Routing Stream wird vorbereitet
	routingStream, err := initRoutingStreams(conn, ctx, cancel)
	if err != nil {
		err = fmt.Errorf("openNodeConnection: %w", err)
		cancel(err)
		return nil, err
	}

	// Der Package Traffic Stream wird vorbereitet
	packageTrafficStream, err := initPackageTrafficStreams(conn, ctx, cancel)
	if err != nil {
		err = fmt.Errorf("openNodeConnection: %w", err)
		cancel(err)
		return nil, err
	}

	// Die Control Routine wird gestartet
	if err = startGoroutineControl(controlStream, ctx, cancel); err != nil {
		err = fmt.Errorf("openNodeConnection: %w", err)
		cancel(err)
		return nil, err
	}

	// Die Routing Routine wird gestartet
	if err = startGoroutineRouting(routingStream, ctx, cancel); err != nil {
		err = fmt.Errorf("openNodeConnection: %w", err)
		cancel(err)
		return nil, err
	}

	// Die Traffic Routine wird gestartet
	if err = startGoroutinePackageTraffic(packageTrafficStream, ctx, cancel); err != nil {
		err = fmt.Errorf("openNodeConnection: %w", err)
		cancel(err)
		return nil, err
	}

	// Die Unreliable Datagrame Sitzungs Funktionen werden erstellt
	if err = startGoroutineUnreliableDatagrammHandle(conn, ctx, cancel); err != nil {
		err = fmt.Errorf("openNodeConnection: %w", err)
		cancel(err)
		return nil, err
	}

	// Wenn es sich um eine eingehende Verbindung handelt, wird Signalisiert dass
	// die Verbindung erfolgreich Initalisiert wurde, wenn es sich um eine Ausgehende,
	// Verbindung handelt, dann wird auf das Signal gewartet.
	if isIncommingConnection {
		if err = controlStream.SignalIncommingInitalizationComplete(); err != nil {
			err = fmt.Errorf("openNodeConnection: %w", err)
			cancel(err)
			return nil, err
		}
		if err = controlStream.WaitInitalizationComplete(); err != nil {
			err = fmt.Errorf("openNodeConnection: %w", err)
			cancel(err)
			return nil, err
		}
	} else {
		if err = controlStream.WaitInitalizationComplete(); err != nil {
			err = fmt.Errorf("openNodeConnection: %w", err)
			cancel(err)
			return nil, err
		}
		if err = controlStream.SignalIncommingInitalizationComplete(); err != nil {
			err = fmt.Errorf("openNodeConnection: %w", err)
			cancel(err)
			return nil, err
		}
	}

	return nodeConn, nil
}
