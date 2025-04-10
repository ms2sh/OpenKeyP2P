package p2p

import (
	"bytes"
	"fmt"
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

func _ControlStreamReaderProcess(conn *NodeP2PConnection, wg *sync.WaitGroup, data []byte) error {
	// Die Paketgröße wird geprüft
	if len(data) < 2 {
		fmt.Println("Invalid data recived")
		return nil
	}

	// Es wird versucht zu ermitteln um was für ein Pakettypen es sich handelt
	switch {
	case bytes.Equal(data[:2], Keepalive[:]) || bytes.Equal(data[:2], KeepaliveReply[:]):
		// Es wird geprüft ob der Restliche Datensatz 32 bytes Groß ist
		if x := len(data[2:]); x != 32 {
			fmt.Println("invalid keepalive id", x)
			return nil
		}

		// Das Keepalive Paket wird beantwortet
		isreplay := bytes.Equal(data[:2], KeepaliveReply[:])
		if err := _EnterKeepaliveResponse(conn, NodeP2PKeepaliveProcessId(data[2:]), isreplay); err != nil {
			conn.contextCancel(err)
		}
	case bytes.Equal(data[:2], RoutingChannelCrawler[:]):
	case bytes.Equal(data[:2], UpdatePOWDiff[:]):
	case bytes.Equal(data[:2], UpdateAutoRoutingQuickSearchTable[:]):
	case bytes.Equal(data[:2], PeerDiscovery[:]):
	default:
		fmt.Println("unkown packet type")
		return nil
	}

	return nil
}

func _PackageTrafficStreamReaderProcess(conn *NodeP2PConnection, wg *sync.WaitGroup, data []byte) error {
	// Die Paketgröße wird geprüft
	if len(data) < 2 {
		fmt.Println("Invalid data recived")
		return nil
	}

	// Es wird versucht zu ermitteln um was für ein Pakettypen es sich handelt
	switch {
	case bytes.Equal(data[:2], Datagramm[:]):
	case bytes.Equal(data[:2], RoutingChannelDatagramm[:]):
	default:
		fmt.Println("unkown packet type")
		return nil
	}

	return nil
}

func _ControlStreamReaderRoutineRootFunction(conn *NodeP2PConnection, wg *sync.WaitGroup) {
	// Gib an ob die Initalisierung durchgeführt wurde
	wasinited := false

	// Wird ausgeführt wenn die Funktion am ende ist
	defer func() {
		// LOG
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Control reader routine stopped %s -> %s", conn.localSocketAddress, conn.remoteSocketAddress)
	}()

	// Die Schleife wird verwendet um Eintreffende Daten zu lesen
	for {
		select {
		case <-conn.ctx.Done(): // Abbruch, wenn der Kontext geschlossen wurde
			return
		default: // Es wird geprüft ob Daten vorhanden sind
			// Es darf nur 1x ein Init Signal gesendet werden
			if !wasinited {
				// Gibt an das der Init vorgang erfolgreich war
				wasinited = true

				// Log
				logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Control reader routine started %s -> %s", conn.localSocketAddress, conn.remoteSocketAddress)

				// Es wird Signalisiert das der Reader vollständig ausgeführt wird
				wg.Done()
			}

			// Es wird auf Daten gewartet
			data, err := conn.controlStream.ReadBytes()
			if err != nil {
				// LOG
				logging.LogError(openkeyp2p.LOG_LEVEL_P2P, "Error by reading data on 'ControlStream' {%s} %s -> %s", err, conn.localSocketAddress, conn.remoteSocketAddress)

				// Der Fehler wird an den Context übergeben
				conn.contextCancel(fmt.Errorf("_ControlStreamReaderRoutineRootFunction: %s", err))

				// Nächster Lesevorgang
				continue
			}

			// Das Paket wird verarbeitet
			if err := _ControlStreamReaderProcess(conn, wg, data); err != nil {
				// LOG
				logging.LogError(openkeyp2p.LOG_LEVEL_P2P, "Error by process data on 'ControlStream' {%s} %s -> %s", err, conn.localSocketAddress, conn.remoteSocketAddress)

				// Der Fehler wird an den Context übergeben
				conn.contextCancel(fmt.Errorf("_ControlStreamReaderRoutineRootFunction: %s", err))

				// Nächster Lesevorgang
				continue
			}
		}
	}
}

func _TrafficStreamReaderRoutineRootFunction(conn *NodeP2PConnection, wg *sync.WaitGroup) {
	// Gib an ob die Initalisierung durchgeführt wurde
	wasinited := false

	// Wird ausgeführt wenn die Funktion am ende ist
	defer func() {
		// LOG
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Package traffic reader routine stopped %s -> %s", conn.localSocketAddress, conn.remoteSocketAddress)
	}()

	// Die Schleife wird verwendet um Eintreffende Daten zu lesen
	for {
		select {
		case <-conn.ctx.Done(): // Abbruch, wenn der Kontext geschlossen wurde
			return
		default: // Es wird geprüft ob Daten vorhanden sind
			// Es darf nur 1x ein Init Signal gesendet werden
			if !wasinited {
				// Gibt an das der Init vorgang erfolgreich war
				wasinited = true

				// Log
				logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Package traffic reader routine started %s -> %s", conn.localSocketAddress, conn.remoteSocketAddress)

				// Es wird Signalisiert das der Reader vollständig ausgeführt wird
				wg.Done()
			}

			// Es wird auf Daten gewartet
			data, err := conn.packageTrafficStream.ReadBytes()
			if err != nil {
				// LOG
				logging.LogError(openkeyp2p.LOG_LEVEL_P2P, "Error by reading data on Package traffic stream: {%s} %s -> %s", err, conn.localSocketAddress, conn.remoteSocketAddress)

				// Der Fehler wird an den Context übergeben
				conn.contextCancel(fmt.Errorf("_ControlStreamReaderRoutineRootFunction: %s", err))

				// Nächster Lesevorgang
				continue
			}

			// Das Paket wird verarbeitet
			if err := _PackageTrafficStreamReaderProcess(conn, wg, data); err != nil {
				// LOG
				logging.LogError(openkeyp2p.LOG_LEVEL_P2P, "Error by reading data on Package traffic stream: {%s} %s -> %s", err, conn.localSocketAddress, conn.remoteSocketAddress)

				// Der Fehler wird an den Context übergeben
				conn.contextCancel(fmt.Errorf("_ControlStreamReaderRoutineRootFunction: %s", err))

				// Nächster Lesevorgang
				continue
			}
		}
	}
}

func _StartReaderRoutinesForNodeConn(conn *NodeP2PConnection, wg *sync.WaitGroup) error {
	// Es wird eine Waiting Group erzeugt
	wgt := new(sync.WaitGroup)
	wgt.Add(2)

	// Der Reader für den Controlstream wird gestartet
	go _ControlStreamReaderRoutineRootFunction(conn, wgt)

	// Der Reader für den Trafficstream wird gestartet
	go _TrafficStreamReaderRoutineRootFunction(conn, wgt)

	// Es wird darauf gewartet dass beide Routinen ausgeführt werden
	wgt.Wait()

	// Log
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Stream Readers started %s -> %s", conn.localSocketAddress, conn.remoteSocketAddress)

	// Es wird Signalisiert das die Writer Routinen ausgeführt werden
	wg.Done()

	return nil
}
