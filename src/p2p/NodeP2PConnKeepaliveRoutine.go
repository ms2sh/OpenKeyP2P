package p2p

import (
	"fmt"
	"sync"
	"time"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

func _StartKeepaliveRoutinesForNodeConn(conn *NodeP2PConnection, wg *sync.WaitGroup) error {
	// Die Parameter werden gepfüft
	if conn == nil || conn.conn == nil {
		return fmt.Errorf("connection must not be nil")
	}
	if conn.contextCancel == nil {
		return fmt.Errorf("context cancel must not be nil")
	}
	if conn.keepaliveTime <= 0 {
		return fmt.Errorf("keepalive time must be positive")
	}

	// Wird als Routine ausgeführt
	go func(conn *NodeP2PConnection) {
		ticker := time.NewTicker(conn.keepaliveTime) // Der Time wartet bis neue Daten gesendet werden
		wasChangesTickerTime := false                // Gibt an das die Zeit des Tickers verändert wurde
		currentKeepaliveInterval := conn.keepaliveTime

		// Wird ausgeführt wenn die Funktion zuende ist
		defer func() {
			ticker.Stop()
			logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Keepalive routine stopped %s -> %s", conn.localSocketAddress, conn.remoteSocketAddress)
		}()

		// Log
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Keepalive routine started %s -> %s", conn.localSocketAddress, conn.remoteSocketAddress)

		// Es wird Signalisiert das die Goroutine ausgeführt wird, damit kann die nächste Routine gestartet werden.
		wg.Done()

		// Die Schleife wird solange ausgeführt bis keine
		for {
			select {
			case <-ticker.C:
				if err := _SendKeepaliveSignal(conn, currentKeepaliveInterval); err != nil {
					// LOG
					logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Error by sending keepalive packet :: %s %s -> %s", err, conn.localSocketAddress, conn.remoteSocketAddress)

					// Es ist ein fehler aufgetreten, sofern die Verbindung nicht geschlossen wurde,
					// wird nun das Schreiben blockiert bis ein Keepalive empfangen wurden
					_SignalWritingLockThenNoKeepaliveResponseRecived(conn)

					// Die Timezeit wird verringer
					ticker.Stop()
					currentKeepaliveInterval = conn.keepaliveTime / 2
					ticker = time.NewTicker(currentKeepaliveInterval)
					wasChangesTickerTime = true

					// Der Vorgang wird neugestartet
					continue
				}

				// Sollte die Tickerzeit verädnert wurden sein, wird sie auf den Standrdwert zurückgesetzt
				if wasChangesTickerTime {
					ticker.Stop()
					currentKeepaliveInterval = conn.keepaliveTime
					ticker = time.NewTicker(currentKeepaliveInterval)
					wasChangesTickerTime = false
				}
			case <-conn.ctx.Done(): // Falls die Verbindung geschlossen wird
				// Der Fehler wird an den Context übergeben
				conn.contextCancel(fmt.Errorf("connection closed"))
				return
			}
		}
	}(conn)

	return nil
}
