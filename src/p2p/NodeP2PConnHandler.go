package p2p

import (
	"context"
	"fmt"
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

func _SyncHandleConnection(conn *NodeP2PConnection) {
	// Die Waitgroup wird erzeugt
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// Die IP Adressen der Verbindung werden abgerufen
	localEndpointStr := getLocalIPFromConn(conn.conn)
	remoteEndpointStr := getRemoteIPAndHostFromConn(conn.conn)

	// Die Leseroutinen werden gestartet
	if err := _StartReaderRoutinesForNodeConn(conn, wg); err != nil {
		return
	}

	// Die Schreibroutinen werden gestartet
	if err := _StartWriterRoutinesForNodeConn(conn, wg); err != nil {
		return
	}

	// Es wird darauf gewartet dass beide Routinen signalisieren das der Vorgang erfolgreich war
	wg.Wait()
	wg.Add(1)

	// Die Keepalive Routine welche prüft wann die Letze erfolgreiche übertragung passiert ist wird gestartet
	if err := _StartKeepaliveRoutinesForNodeConn(conn, wg); err != nil {
		logging.LogError(openkeyp2p.LOG_LEVEL_P2P, "Error by starting 'keepalive' routine. Error(%s) %s -> %s", err, localEndpointStr, remoteEndpointStr)
		return
	}

	// Es wird darauf gewartet das die Ping Routine gestartet wurde
	wg.Wait()

	// Log
	logtxt := "A new connection has been established %s -> %s"
	logtxt = fmt.Sprintf("%s\n   -> Version: %s", logtxt, openkeyp2p.ParseVersion(conn.controlStream.GetDestinationVersion()))
	logtxt = fmt.Sprintf("%s\n   -> CMTU: %d", logtxt, conn.controlStream.GetMTU())
	logtxt = fmt.Sprintf("%s\n   -> ACK-Peer-Packet: %t", logtxt, conn.controlStream.GetACKPeerPacket())
	if conn.config.HasConfigEntryWithValue("auto-routing", "yes") {
		logtxt = logtxt + "\n   -> AutoRouting: Enabled"
	} else {
		logtxt = logtxt + "\n   -> AutoRouting: Disabeld"
	}
	logging.LogInfo(openkeyp2p.LOG_LEVEL_P2P, logtxt, localEndpointStr, remoteEndpointStr)

	// Es wird darauf gewartet dass der Context geschlossen wird
	<-conn.ctx.Done()

	// Ermitteln, warum der Kontext beendet wurde
	switch conn.ctx.Err() {
	case context.Canceled:
		logging.LogInfo(openkeyp2p.LOG_LEVEL_P2P, "The connection has been closed %s -> %s", localEndpointStr, remoteEndpointStr)
	case context.DeadlineExceeded:
		logging.LogInfo(openkeyp2p.LOG_LEVEL_P2P, "The connection has been closed %s -> %s", localEndpointStr, remoteEndpointStr)
	default:
		fmt.Println("Unbekannter Abbruchgrund:", conn.ctx.Err())
	}
}

func _AsyncHandleConnection(conn *NodeP2PConnection, callback func()) {
	go func() {
		_SyncHandleConnection(conn)
		callback()
	}()
}
