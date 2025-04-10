package p2p

import (
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

func _ControlStreamWriterRoutineRootFunction(conn *NodeP2PConnection, wg *sync.WaitGroup) {
	wasinited := false
	for {
		select {
		case <-conn.ctx.Done(): // Abbruch, wenn der Kontext geschlossen wurde
			return
		default:
			// Es darf nur 1x ein Init Signal gesendet werden
			if !wasinited {
				wasinited = true
				wg.Done()
			}

			// Es wird geprüft ob neue Daten verfügbar sind, wenn ja werden diese gesendet
			data, err := conn.writerControlBuffer.Get()
			if err != nil {
				conn.contextCancel(err)
				return
			}

			// Die Daten werden geschrieben
			if err := conn.controlStream.WriteBytes(data.([]byte)); err != nil {
				conn.contextCancel(err)
				return
			}
		}
	}
}

func _TrafficStreamWriterRoutineRootFunction(conn *NodeP2PConnection, wg *sync.WaitGroup) {
	wg.Done()
}

func _StartWriterRoutinesForNodeConn(conn *NodeP2PConnection, wg *sync.WaitGroup) error {
	// Es wird eine Waiting Group erzeugt
	wgt := new(sync.WaitGroup)
	wgt.Add(2)

	// Der Reader für den Controlstream wird gestartet
	go _ControlStreamWriterRoutineRootFunction(conn, wgt)

	// Der Reader für den Trafficstream wird gestartet
	go _TrafficStreamWriterRoutineRootFunction(conn, wgt)

	// Es wird darauf gewartet dass beide Routinen ausgeführt werden
	wgt.Wait()

	// Log
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Stream Writers started %s -> %s", conn.localSocketAddress, conn.remoteSocketAddress)

	// Es wird Signalisiert das die Writer Routinen ausgeführt werden
	wg.Done()

	return nil
}
