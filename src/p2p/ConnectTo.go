package p2p

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"github.com/quic-go/quic-go"
)

func ConnectToNode(nodeUri string, tlsConfig *tls.Config, config *NodeP2PConnectionConfig) error {
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
	var finalNodeAddress string
	switch IdentifyAddressType(parsedURL.Hostname()) {
	case AddressTypeIPv4Address:
		host, port, err := net.SplitHostPort(parsedURL.Host)
		if err != nil {
			return err
		}
		finalNodeAddress = fmt.Sprintf("%s:%s", host, port)
	case AddressTypeIPv6Address:
		finalNodeAddress = fmt.Sprintf("[%s]:%s", parsedURL.Host, parsedURL.Port())
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
	conn, err := quic.DialAddr(ctx, finalNodeAddress, tlsConfig, nil)
	if err != nil {
		err = fmt.Errorf("ConnectToNode: %w", err)
		cancel(err)
		return err
	}

	// Die Verbindung wird Initialisiert
	nodeConn, err := _InitNodeConn(ctx, cancel, false, config, conn)
	if err != nil {
		return err
	}

	// Die Verbindung wird vorbereitet
	if err := _VarsAddNodeConnection(nodeConn); err != nil {
		return err
	}

	// Falls alles passt, wird eine Verbindung hergestellt (hier Dummy-Rückgabe)
	return nil
}
