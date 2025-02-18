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

func initControlStreams(conn quic.Connection, context context.Context, contextCancel context.CancelCauseFunc) (*_NodeP2PConnectionControlStream, error) {
	// Der Chan nimmt den Eingehenden Stream entgegen
	inStreamChan := make(chan *_ChanStreamErrorResult, 1)

	// Startet die Routine welche den Control Stream der Gegenseite Akzeptiert
	go _streamAcceptOppositeSide(conn, context, contextCancel, inStreamChan)

	// Es wird ein Control Stream mit der gegenseite aufgebaut
	outStream, err := conn.OpenStreamSync(context)
	if err != nil {
		return nil, fmt.Errorf("initControlStreams: Fehler beim Öffnen des Control Streams: %w", err)
	}

	// Es wird geprüpft ob der Eingehende Stream geöffnet wurde
	inStreamResult := <-inStreamChan
	if inStreamResult == nil {
		return nil, fmt.Errorf("initControlStreams: connection error")
	}
	if inStreamResult.err != nil {
		return nil, fmt.Errorf("initControlStreams: " + inStreamResult.err.Error())
	}

	// Die Streams werden abgespeichert
	cstream := &_NodeP2PConnectionControlStream{
		outControlStream: outStream,
		inControlStream:  inStreamResult.stream,
	}

	return cstream, nil
}

func initRoutingStreams(conn quic.Connection, context context.Context, cancel context.CancelCauseFunc) (*_NodeP2PConnectionRoutingStream, error) {
	// Der Chan nimmt den Eingehenden Stream entgegen
	inStreamChan := make(chan *_ChanStreamErrorResult, 1)

	// Startet die Routine welche den Control Stream der Gegenseite Akzeptiert
	go _streamAcceptOppositeSide(conn, context, cancel, inStreamChan)

	// Es wird ein Control Stream mit der gegenseite aufgebaut
	outStream, err := conn.OpenStreamSync(context)
	if err != nil {
		return nil, fmt.Errorf("initControlStreams: Fehler beim Öffnen des Control Streams: %w", err)
	}

	// Es wird geprüpft ob der Eingehende Stream geöffnet wurde
	inStreamResult := <-inStreamChan
	if inStreamResult == nil {
		return nil, fmt.Errorf("initControlStreams: connection error")
	}
	if inStreamResult.err != nil {
		return nil, fmt.Errorf("initControlStreams: " + inStreamResult.err.Error())
	}

	// Die Streams werden abgespeichert
	rstream := &_NodeP2PConnectionRoutingStream{
		inRoutingStream:  outStream,
		outRoutingStream: inStreamResult.stream,
	}

	return rstream, nil
}

func initPackageTrafficStreams(conn quic.Connection, context context.Context, contextCancel context.CancelCauseFunc) (*_NodeP2PConnectionPackageTrafficStream, error) {
	// Der Chan nimmt den Eingehenden Stream entgegen
	inStreamChan := make(chan *_ChanStreamErrorResult, 1)

	// Startet die Routine welche den Control Stream der Gegenseite Akzeptiert
	go _streamAcceptOppositeSide(conn, context, contextCancel, inStreamChan)

	// Es wird ein Control Stream mit der gegenseite aufgebaut
	outStream, err := conn.OpenStreamSync(context)
	if err != nil {
		return nil, fmt.Errorf("initControlStreams: Fehler beim Öffnen des Control Streams: %w", err)
	}

	// Es wird geprüpft ob der Eingehende Stream geöffnet wurde
	inStreamResult := <-inStreamChan
	if inStreamResult == nil {
		return nil, fmt.Errorf("initControlStreams: connection error")
	}
	if inStreamResult.err != nil {
		return nil, fmt.Errorf("initControlStreams: " + inStreamResult.err.Error())
	}

	// Die Streams werden abgespeichert
	ptstream := &_NodeP2PConnectionPackageTrafficStream{
		inPackageTrafficStream:  outStream,
		outPackageTrafficStream: inStreamResult.stream,
	}

	return ptstream, nil
}

func startGoroutineControl(routingStream *_NodeP2PConnectionControlStream, context context.Context, contextCancel context.CancelCauseFunc) error {
	return nil
}

func startGoroutineRouting(routingStream *_NodeP2PConnectionRoutingStream, context context.Context, contextCancel context.CancelCauseFunc) error {
	return nil
}

func startGoroutinePackageTraffic(routingStream *_NodeP2PConnectionPackageTrafficStream, context context.Context, contextCancel context.CancelCauseFunc) error {
	return nil
}

func startGoroutineUnreliableDatagrammHandle(conn quic.Connection, context context.Context, contextCancel context.CancelCauseFunc) error {
	for {
		msg, err := conn.ReceiveDatagram(context)
		if err != nil {
			log.Println("Fehler beim Empfangen von Datagrammen:", err)
			return nil
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
