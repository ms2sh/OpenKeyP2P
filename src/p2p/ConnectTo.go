package p2p

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"github.com/quic-go/quic-go"
)

func ConnectTo(nodeUri string, tlsConfig *tls.Config, config NodeP2PConnectionConfig) error {
	controlUnlockByEnd := false
	defer func() {
		if !controlUnlockByEnd {
			return
		}
		controlLock.Unlock()
	}()

	if !_VarsWasSetuped() {
		return fmt.Errorf("you must setup p2p node functions, call Setup()")
	}

	parsedURL, err := url.Parse(nodeUri)
	if err != nil {
		return fmt.Errorf("ConnectToNode: " + err.Error())
	}

	// Erlaubt nur "quic" als Protokoll
	if parsedURL.Scheme != "quic" {
		return fmt.Errorf("only quic as protocol allowed")
	}

	// Der Host muss vorhanden sein
	if parsedURL.Hostname() == "" {
		return fmt.Errorf("no host found")
	}

	// Der Port muss vorhanden sein
	if parsedURL.Port() == "" {
		return fmt.Errorf("no port found")
	}

	// Der Pfad muss leer sein
	if parsedURL.Path != "" {
		return fmt.Errorf("path must be empty")
	}

	// Query-Parameter müssen leer sein
	if parsedURL.RawQuery != "" {
		return fmt.Errorf("query parameters are not allowed")
	}

	// Es wird eine Verbindung mit dem Node hergestellt
	useAsProxy := false
	var finalNodeAddress string
	iapt := IdentifyAddressType(parsedURL.Hostname())
	switch iapt {
	case AddressTypeIPv4Address:
		host, port, err := net.SplitHostPort(parsedURL.Host)
		if err != nil {
			return err
		}
		finalNodeAddress = fmt.Sprintf("%s:%s", host, port)
	case AddressTypeIPv6Address:
		finalNodeAddress = fmt.Sprintf("[%s]:%s", parsedURL.Host, parsedURL.Port())
	case AddressTypeOnionV3:
	case AddressTypeDomain:
		ipadr, err := GetIpFromDomain(parsedURL.Hostname())
		if err != nil {
			return fmt.Errorf("can't find ip for domain %s", parsedURL.Hostname())
		}

		finalNodeAddress = fmt.Sprintf("%s:%s", ipadr, parsedURL.Port())
	case AddressTypeUnkown:
		return fmt.Errorf("unkown protocol '%s' , only quic supported", parsedURL.Scheme)
	default:
		return fmt.Errorf("unkown protocol, only quic supported: %s", parsedURL.Scheme)
	}

	// Jeder Client bekommt seinen eigenen Kontext
	ctx, cancel := context.WithCancelCause(context.Background())

	// Die Quic Verbindung wird aufgebaut
	var nodeConn *NodeP2PConnection
	var conn quic.Connection
	if !useAsProxy {
		// Die Quic Verbindung wird aufgebaut
		var err error
		conn, err = quic.DialAddr(ctx, finalNodeAddress, tlsConfig, nil)
		if err != nil {
			err = fmt.Errorf("ConnectToNode: %w", err)
			cancel(err)
			return err
		}

		// Es wird das passende Interface für die Lokale IP-Adresse ermittelt
		ip, _, err := net.SplitHostPort(conn.LocalAddr().String())
		if err != nil {
			err = fmt.Errorf("ConnectToNode: %w", err)
			cancel(err)
			return err
		}
		if ip == "::" || ip == "0.0.0.0" {
			ip = getLocalIPFromConn(conn)
		}
		localhostNetworkInterface, err := getInterfaceByIP(ip)
		if err != nil {
			err = fmt.Errorf("ConnectToNode: %w", err)
			cancel(err)
			return err
		}

		// Die Verbindung wird Initialisiert
		nodeConn, err = _InitQUICNodeConn(localhostNetworkInterface, ctx, cancel, false, config, conn)
		if err != nil {
			cancel(err)
			return err
		}
	} else {
		cancel(nil)
		return fmt.Errorf("not supported parameter")
	}

	// Die Verbindung wird vorbereitet
	if err := _VarsAddNodeConnection(nodeConn); err != nil {
		cancel(err)
		return err
	}

	// Die Handler Routine wird gestartet
	_AsyncHandleConnection(nodeConn, func() {
		// Die Verbindung wurde getrennt, sie wird aus dem Verbindungsspeicher entfernt
		_VarsDeleteNodeConnection(nodeConn)
	})

	// Falls alles passt, wird eine Verbindung hergestellt (hier Dummy-Rückgabe)
	return nil
}
