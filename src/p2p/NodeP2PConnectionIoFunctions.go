package p2p

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

func _SendKeepaliveSignal(o *NodeP2PConnection, duration time.Duration) error {
	// Es wird ein Zufälliger 256 Bit Wert erzeugt
	bitvalue, err := _GenerateRandom256BitValue()
	if err != nil {
		return err
	}

	// Es wird ein neuer Context erzeugt
	ctx, cancel := context.WithTimeoutCause(o.ctx, duration, fmt.Errorf("keepalive timed out"))
	defer cancel()

	// Der Vorgang wird zwischengespeichert
	o.localKeepalivePacketIds.Store(hex.EncodeToString(bitvalue[:]), &_NodeP2pKeepaliveProcess{Ctx: ctx, Cancel: cancel, LMutex: new(sync.Mutex), Finish: false})
	defer o.localKeepalivePacketIds.Delete(hex.EncodeToString(bitvalue[:]))

	// Das Finale Datenpaket wird gebaut
	finalDataPacket := append([]byte(Keepalive[:]), bitvalue...)

	// LOG
	localEndpointStr := getLocalIPFromConn(o.conn)
	remoteEndpointStr := getRemoteIPAndHostFromConn(o.conn)
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P_QUIC, "Try to write Keepalive Packet, %s bytes %s -> %s", hex.EncodeToString(bitvalue), localEndpointStr, remoteEndpointStr)

	// Das Paket wird gesendet
	o.writerControlBuffer.Prepend(finalDataPacket)

	// Warten, bis der Kontext abgeschlossen oder abgelaufen ist
	<-ctx.Done()

	// Prüfen, ob der Timeout die Ursache war
	if errors.Is(context.Cause(ctx), context.DeadlineExceeded) {
		return ErrTimeout
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

func _EnterKeepaliveResponse(o *NodeP2PConnection, kasid NodeP2PKeepaliveProcessId, isResponse bool) error {
	// LOG
	localEndpointStr := getLocalIPFromConn(o.conn)
	remoteEndpointStr := getRemoteIPAndHostFromConn(o.conn)
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P_QUIC, "Enter Keepalive Packet, %s bytes %s -> %s", hex.EncodeToString(kasid[:]), localEndpointStr, remoteEndpointStr)

	// Es wird geprüft ob es einen Lokalen vorgang gibt
	ctxFinal, foundit := o.localKeepalivePacketIds.Load(hex.EncodeToString(kasid[:]))
	if foundit {
		ctxFinalRc := ctxFinal.(*_NodeP2pKeepaliveProcess)
		ctxFinalRc.Cancel()
		return nil
	}

	// Sollte es sich um ein Response sein, wird an dieser Stelle ein Fehler ausgelöst,
	// es handelt sich um ein nicht Bekanntes Keepalive Paket welches als Response angegeben wurde
	if isResponse {
		return fmt.Errorf("unkown response id")
	}

	// Das Paket wird an den Absender zurückgesendet
	// Das Finale Datenpaket wird gebaut
	finalDataPacket := append([]byte(Keepalive[:]), kasid[:]...)

	// Das Paket wird gesendet
	o.writerControlBuffer.Prepend(finalDataPacket)

	// Es ist kein Fehler aufgetreten
	return nil
}
