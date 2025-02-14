package p2p

import (
	"context"
	"fmt"
	"log"

	"github.com/quic-go/quic-go"
)

func _streamAcceptOppositeSide(conn quic.Connection, context context.Context, cancel context.CancelCauseFunc, inStream chan *_ChanStreamErrorResult) {
	// Warte auf einen Stream (blockiert, bis ein Stream geöffnet wird oder die Verbindung geschlossen wird)
	stream, err := conn.AcceptStream(context)
	if err != nil {
		cancel(fmt.Errorf("_streamAcceptOppositeSide: %s", err.Error()))
		close(inStream)
		return
	}

	inStream <- &_ChanStreamErrorResult{err: nil, stream: stream}
}

func initControlStreams(nodeConn *NodeP2PConnection) error {
	// Der Chan nimmt den Eingehenden Stream entgegen
	inStreamChan := make(chan *_ChanStreamErrorResult, 1)

	// Startet die Routine welche den Control Stream der Gegenseite Akzeptiert
	go _streamAcceptOppositeSide(nodeConn.conn, nodeConn.context, nodeConn.contextCancel, inStreamChan)

	// Es wird ein Control Stream mit der gegenseite aufgebaut
	outStream, err := nodeConn.conn.OpenStreamSync(nodeConn.context)
	if err != nil {
		return fmt.Errorf("initControlStreams: Fehler beim Öffnen des Control Streams: %w", err)
	}

	// Es wird geprüpft ob der Eingehende Stream geöffnet wurde
	inStreamResult := <-inStreamChan
	if inStreamResult == nil {
		return fmt.Errorf("initControlStreams: connection error")
	}
	if inStreamResult.err != nil {
		return fmt.Errorf("initControlStreams: " + inStreamResult.err.Error())
	}

	// Die Streams werden abgespeichert
	cstream := &NodeP2PConnectionControlStream{
		outControlStream: outStream,
		inControlStream:  inStreamResult.stream,
	}

	// Das Stream Objekt wird im Verbindungsobjekt zwischengespeichert
	nodeConn.controlStream = cstream

	return nil
}

func initRoutingStreams(nodeConn *NodeP2PConnection) error {
	// Der Chan nimmt den Eingehenden Stream entgegen
	inStreamChan := make(chan *_ChanStreamErrorResult, 1)

	// Startet die Routine welche den Control Stream der Gegenseite Akzeptiert
	go _streamAcceptOppositeSide(nodeConn.conn, nodeConn.context, nodeConn.contextCancel, inStreamChan)

	// Es wird ein Control Stream mit der gegenseite aufgebaut
	outStream, err := nodeConn.conn.OpenStreamSync(nodeConn.context)
	if err != nil {
		return fmt.Errorf("initControlStreams: Fehler beim Öffnen des Control Streams: %w", err)
	}

	// Es wird geprüpft ob der Eingehende Stream geöffnet wurde
	inStreamResult := <-inStreamChan
	if inStreamResult == nil {
		return fmt.Errorf("initControlStreams: connection error")
	}
	if inStreamResult.err != nil {
		return fmt.Errorf("initControlStreams: " + inStreamResult.err.Error())
	}

	// Die Streams werden abgespeichert
	cstream := &NodeP2PConnectionControlStream{
		outControlStream: outStream,
		inControlStream:  inStreamResult.stream,
	}

	// Das Stream Objekt wird im Verbindungsobjekt zwischengespeichert
	nodeConn.controlStream = cstream

	return nil
}

func initPackageTrafficStreams(nodeConn *NodeP2PConnection) error {
	// Der Chan nimmt den Eingehenden Stream entgegen
	inStreamChan := make(chan *_ChanStreamErrorResult, 1)

	// Startet die Routine welche den Control Stream der Gegenseite Akzeptiert
	go _streamAcceptOppositeSide(nodeConn.conn, nodeConn.context, nodeConn.contextCancel, inStreamChan)

	// Es wird ein Control Stream mit der gegenseite aufgebaut
	outStream, err := nodeConn.conn.OpenStreamSync(nodeConn.context)
	if err != nil {
		return fmt.Errorf("initControlStreams: Fehler beim Öffnen des Control Streams: %w", err)
	}

	// Es wird geprüpft ob der Eingehende Stream geöffnet wurde
	inStreamResult := <-inStreamChan
	if inStreamResult == nil {
		return fmt.Errorf("initControlStreams: connection error")
	}
	if inStreamResult.err != nil {
		return fmt.Errorf("initControlStreams: " + inStreamResult.err.Error())
	}

	// Die Streams werden abgespeichert
	cstream := &NodeP2PConnectionControlStream{
		outControlStream: outStream,
		inControlStream:  inStreamResult.stream,
	}

	// Das Stream Objekt wird im Verbindungsobjekt zwischengespeichert
	nodeConn.controlStream = cstream

	return nil
}

func initUnreliableDatagramsHandle(nodeConn *NodeP2PConnection) {
	for {
		msg, err := nodeConn.conn.ReceiveDatagram(nodeConn.context)
		if err != nil {
			log.Println("Fehler beim Empfangen von Datagrammen:", err)
			return
		}

		// Erste Byte als Typkennung auswerten
		if len(msg) < 1 {
			log.Println("Ungültiges Datagramm empfangen")
			continue
		}

		msgType := msg[0]
		payload := msg[1:]

		switch msgType {
		case 0x01:
			fmt.Println("[Server] Steuerungsnachricht empfangen:", string(payload))
		case 0x02:
			fmt.Println("[Server] Chat-Nachricht empfangen:", string(payload))
		case 0x03:
			fmt.Println("[Server] Spielfigur-Koordinaten empfangen:", string(payload))
		default:
			fmt.Println("[Server] Unbekannter Nachrichtentyp:", msgType, "Daten:", string(payload))
		}
	}
}
