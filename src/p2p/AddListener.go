package p2p

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
	"github.com/quic-go/quic-go"
)

func _HandleSession(session quic.Connection, config *NodeP2PListenerConfig) {
	ctx, cancel := context.WithCancelCause(context.Background())

	_, err := _InitNodeConn(ctx, cancel, true, &NodeP2PConnectionConfig{AllowAutoRouting: config.AllowAutoRouting, AllowTrafficForwarding: config.AllowTrafficForwarding}, session)
	if err != nil {
		ert := fmt.Errorf("fehler beim Initalisieren einer Verbindung: %v", err)
		fmt.Println(ert)
		cancel(ert)
	}
}

func _StartListenerGoroutine(listeneraddr openkeyp2p.LocalListenerAddress, listener *NodeP2Listener, config *NodeP2PListenerConfig) {
	logging.LogInfo(openkeyp2p.LOG_LEVEL_P2P, "Accepts incoming connections on %s", listeneraddr)
	go func() {
		for {
			// Neue QUIC-Verbindung akzeptieren
			session, err := listener.listener.Accept(context.Background())
			if err != nil {
				log.Printf("Fehler beim Akzeptieren einer Verbindung: %v", err)
				continue
			}
			go _HandleSession(session, config)
		}
	}()
}

func AddListener(localIp string, localPort uint32, tlsConfig *tls.Config, config *NodeP2PListenerConfig) error {
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

	//Log
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
