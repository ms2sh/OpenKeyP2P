package p2p

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
	"github.com/quic-go/quic-go"
)

func _StreamWriteBytePacket(conn quic.Connection, stream quic.Stream, helloPackage []byte) error {
	// Log
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P_QUIC, "Send Stream hello package %s -> %s", conn.LocalAddr(), conn.RemoteAddr())

	// Das Hallo Paket wird übertragen
	if _, err := stream.Write(openkeyp2p.HELLO_STREAM_PACKAGE); err != nil {
		return err
	}

	return nil
}

func _StreamReadBytePacket(conn quic.Connection, stream quic.Stream) ([]byte, error) {
	// Die Länge des Byte-Slices auf die Länge des zu erwartenden Pakets setzen
	helloStromBytesSlice := make([]byte, len(openkeyp2p.HELLO_STREAM_PACKAGE))

	// Das Hello Stream Paket der gegenseite wird eingelesen
	n, err := io.ReadFull(stream, helloStromBytesSlice)
	if err != nil {
		// Wenn beim Lesen ein Fehler auftritt, behandeln wir diesen Fehler
		return nil, err
	}

	// Es wird geprüft ob die Eingetroffnenen Datenmengen korrekt sind
	if n != len(helloStromBytesSlice) {
		return nil, fmt.Errorf("invalid hello stream package")
	}

	// Es wird geprüft ob es sich um ein Korrektes Hello Stream Package handelt
	if !bytes.Equal(helloStromBytesSlice, openkeyp2p.HELLO_STREAM_PACKAGE) {
		return nil, fmt.Errorf("invalid hello stream package")
	}

	// Log
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P_QUIC, "Recived Stream hello package %s -> %s", conn.LocalAddr(), conn.RemoteAddr())

	return nil, nil
}

func _TryOpenQuicBidirectionalStream(isIncommingConnection bool, conn quic.Connection, helloPackage []byte, ctx context.Context) (*QuicBidirectionalStream, error) {
	// Die Contexts für den Control Stream werden erzeugt
	ctx, cancel := context.WithCancelCause(ctx)

	// Es wird selektiert, ob es sich um eine eingehende oder um eine ausgehende Verbindung handelt
	var inStream quic.Stream
	var outStream quic.Stream
	var streamErr error
	var recivedPacket []byte
	if isIncommingConnection {
		// Es wird ein ausgehe4nder Stream geöffnet
		outStream, streamErr = conn.OpenStreamSync(ctx)
		if streamErr != nil {
			cancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, streamErr
		}

		// Das Hello Stream Package wird an die Gegenseite übertragen
		if err := _StreamWriteBytePacket(conn, outStream, helloPackage); err != nil {
			cancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, err
		}

		// Es wird auf einen eingehenden Stream gewartet
		inStream, streamErr = conn.AcceptStream(ctx)
		if streamErr != nil {
			cancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, streamErr
		}

		// Es wird auf das Eingehende Hello Stream Package gewartet
		var err error
		recivedPacket, err = _StreamReadBytePacket(conn, inStream)
		if err != nil {
			cancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, err
		}
	} else {
		// Es wird auf einen eingehenden Stream gewartet
		inStream, streamErr = conn.AcceptStream(ctx)
		if streamErr != nil {
			cancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, streamErr
		}

		// Es wird auf das HelloPackage der gegenseite gewartet
		var err error
		recivedPacket, err = _StreamReadBytePacket(conn, inStream)
		if err != nil {
			cancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, err
		}

		// Es wird eine ausgehende Verbindung aufgebaut
		outStream, streamErr = conn.OpenStreamSync(ctx)
		if streamErr != nil {
			cancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, streamErr
		}

		// Das Hello Stream Package wird an die Gegenseite übertragen
		if err := _StreamWriteBytePacket(conn, outStream, helloPackage); err != nil {
			cancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, err
		}
	}

	// Das Finale Objekt wird erzeugt
	finalObject := &QuicBidirectionalStream{
		inStream:                inStream,
		outStream:               outStream,
		lock:                    new(sync.Mutex),
		ctx:                     ctx,
		ctxCancle:               cancel,
		_recivedHelloBytePacket: recivedPacket,
		_sendHelloBytePacket:    helloPackage,
	}

	// Log
	if isIncommingConnection {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Bidirectional streams were generated %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
	} else {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Bidirectional streams were generated %s -> %s", conn.LocalAddr(), conn.RemoteAddr())
	}

	// Das Objekt wird zurückgegeben
	return finalObject, nil
}

func (q *QuicBidirectionalStream) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Falls bereits geschlossen, nichts tun
	if q.ctx.Err() != nil {
		return
	}

	// Streams schließen
	q.inStream.Close()
	q.outStream.Close()

	// Kontext beenden
	q.ctxCancle(fmt.Errorf("stream closed"))
}

func (q *QuicBidirectionalStream) WriteBytes(byts []byte) error {
	return nil
}

func (q *QuicBidirectionalStream) ReadBytes() ([]byte, error) {
	return nil, nil
}
