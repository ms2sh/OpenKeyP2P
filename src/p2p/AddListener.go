package p2p

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
	"github.com/quic-go/quic-go"
)

func _HandleSession(session quic.Connection, listenerConfig *NodeP2PListenerConfig) {
	// Context erzeugen
	ctx, cancel := context.WithCancelCause(context.Background())

	// Ermittelt das passende Netzwerkinterface anhand der IP Adresse
	ip, _, err := net.SplitHostPort(session.LocalAddr().String())
	if err != nil {
		err = fmt.Errorf("_HandleSession: %w", err)
		cancel(err)
		return
	}
	localhostNetworkInterface, err := getInterfaceByIP(ip)
	if err != nil {
		err = fmt.Errorf("_HandleSession: %w", err)
		cancel(err)
		return
	}

	// Verbindung wird Initalisieren
	conn, err := _InitQUICNodeConn(localhostNetworkInterface, ctx, cancel, true, listenerConfig.GetConnectionConfig(), session)
	if err != nil {
		ert := fmt.Errorf("fehler beim Initalisieren einer Verbindung: %v", err)
		cancel(ert)
		return
	}

	// Verbindung wird Global zwischengespeichert
	if err := _VarsAddNodeConnection(conn); err != nil {
		cancel(err)
		return
	}

	// Der Handler wird gestartet
	_SyncHandleConnection(conn)

	// Die Verbindung wird Global gelöscht
	_VarsDeleteNodeConnection(conn)
}

func _StartListenerGoroutine(listeneraddr openkeyp2p.LocalListenerAddress, listener *NodeP2Listener, config *NodeP2PListenerConfig) {
	logging.LogInfo(openkeyp2p.LOG_LEVEL_P2P, "Accepts incoming connections on %s", listeneraddr)
	go func() {
		for {
			// Neue QUIC-Verbindung akzeptieren
			session, err := listener.listener.Accept(context.Background())
			if err != nil {
				logging.LogError(openkeyp2p.LOG_LEVEL_P2P, "Error by accepting connection %s", err, listeneraddr)
				continue
			}

			// Die Lokale sowie die Remote IP werden abgerufen
			remoteEndpointStr := getRemoteIPAndHostFromConn(session)

			// LOG
			logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Incoming connection accepted %s -> %s", remoteEndpointStr, listeneraddr)

			// Falls NIST ECC genutzt wird, Verbindung weiterverarbeiten
			go _HandleSession(session, config)
		}
	}()
}

func AddListener(localIp string, localPort uint32, tlsConfig *tls.Config, config *NodeP2PListenerConfig) error {
	// Prüft ob die Gloablen Variablen Initalisiert wurden
	if !_VarsWasSetuped() {
		return fmt.Errorf("you must setup p2p node functions, call Setup()")
	}

	// Die Lokale IP wird geprüft
	var finalAddress string
	localIpType := IdentifyAddressType(localIp)
	if localIpType == AddressTypeIPv4Address {
		finalAddress = fmt.Sprintf("%s:%d", localIp, localPort)
	} else if localIpType == AddressTypeIPv6Address {
		finalAddress = fmt.Sprintf("[%s]:%d", localIp, localPort)
	} else {
		return fmt.Errorf("invalid local ip address type")
	}

	// LOG
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "A new listener is started on %s", finalAddress)

	// UDP-Listener erstellen
	addr, err := net.ResolveUDPAddr("udp", finalAddress)
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	// QUIC-Listener starten
	listener, err := quic.Listen(udpConn, tlsConfig, nil)
	if err != nil {
		return err
	}

	// Das Rückgabe Objekt wird erstellt
	resolve := &NodeP2Listener{
		config:   config,
		listener: listener,
		lock:     new(sync.Mutex),
	}

	// Die Goroutine für den Listener wird gestaret
	_StartListenerGoroutine(openkeyp2p.LocalListenerAddress(finalAddress), resolve, config)

	// Das Objket wird zurückgegeben
	return nil
}
