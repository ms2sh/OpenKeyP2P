package p2p

import (
	"crypto/tls"
	"fmt"
	"net/url"
)

func ConnectToNode(nodeUri string, tlsConfig *tls.Config, config *NodeP2PConnectionConfig) (*NodeP2PConnection, error) {
	controlUnlockByEnd := false
	defer func() {
		if !controlUnlockByEnd {
			return
		}
		controlLock.Unlock()
	}()

	if !wasSetuped {
		return nil, fmt.Errorf("you must setup p2p node functions, call Setup()")
	}

	parsedURL, err := url.Parse(nodeUri)
	if err != nil {
		return nil, fmt.Errorf("ConnectToNode: " + err.Error())
	}

	// Erlaubt nur "quic" als Protokoll
	if parsedURL.Scheme != "quic" {
		return nil, fmt.Errorf("only quic as protocol allowed")
	}

	// Der Host muss vorhanden sein
	if parsedURL.Hostname() == "" {
		return nil, fmt.Errorf("no host found")
	}

	// Der Port muss vorhanden sein
	if parsedURL.Port() == "" {
		return nil, fmt.Errorf("no port found")
	}

	// Der Pfad muss leer sein
	if parsedURL.Path != "" {
		return nil, fmt.Errorf("path must be empty")
	}

	// Query-Parameter müssen leer sein
	if parsedURL.RawQuery != "" {
		return nil, fmt.Errorf("query parameters are not allowed")
	}

	// Es wird eine Verbindung mit dem Node hergestellt
	var finalNodeAddress string
	switch IdentifyAddressType(parsedURL.Hostname()) {
	case AddressTypeIPv4Address:
		finalNodeAddress = fmt.Sprintf("%s:%s", parsedURL.Host, parsedURL.Port())
	case AddressTypeIPv6Address:
		finalNodeAddress = fmt.Sprintf("[%s]:%s", parsedURL.Host, parsedURL.Port())
	case AddressTypeDomain:
		ipadr, err := GetIpFromDomain(parsedURL.Hostname())
		if err != nil {
			return nil, fmt.Errorf("can't find ip for domain %s", parsedURL.Hostname())
		}

		finalNodeAddress = fmt.Sprintf("%s:%s", ipadr, parsedURL.Port())
	case AddressTypeUnkown:
		return nil, fmt.Errorf("unkown protocol '%s' , only quic supported", parsedURL.Scheme)
	default:
		return nil, fmt.Errorf("unkown protocol, only quic supported", parsedURL.Scheme)
	}

	// Die Eigentliche Verbindung wird aufgebaut
	_, err = openNodeConnection(finalNodeAddress, false, tlsConfig, config)
	if err != nil {
		return nil, fmt.Errorf("ConnectToNode: " + err.Error())
	}

	// Falls alles passt, wird eine Verbindung hergestellt (hier Dummy-Rückgabe)
	return &NodeP2PConnection{}, nil
}
