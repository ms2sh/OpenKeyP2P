package p2p

import (
	"fmt"
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

func _ControlStreamReaderRoutineRootFunction(conn *NodeP2PConnection, wg *sync.WaitGroup) {
	wasinited := false
	for {
		select {
		case <-conn.context.Done(): // Abbruch, wenn der Kontext geschlossen wurde
			return
		default:
			// Es darf nur 1x ein Init Signal gesendet werden
			if !wasinited {
				wasinited = true
				wg.Done()
			}

			// Es wird auf Daten gewartet
			data, err := conn.controlStream.ReadBytes()
			if err != nil {
				fmt.Println("By", err)
				conn.contextCancel(fmt.Errorf("_ControlStreamReaderRoutineRootFunction: %s", err))
				return
			}

			fmt.Println(data)
		}
	}
}

func _TrafficStreamReaderRoutineRootFunction(conn *NodeP2PConnection, wg *sync.WaitGroup) {
	wg.Done()
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
	localEndpointStr := getLocalIPAndHostFromConn(conn.conn)
	remoteEndpointStr := getRemoteIPAndHostFromConn(conn.conn)
	if conn.isIncommingConnection {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Stream Readers started %s -> %s", remoteEndpointStr, localEndpointStr)
	} else {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Stream Readers started %s -> %s", localEndpointStr, remoteEndpointStr)
	}

	// Es wird Signalisiert das die Writer Routinen ausgeführt werden
	wg.Done()

	return nil
}
