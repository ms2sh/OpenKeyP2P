package p2p

import (
	"context"
	"fmt"
	"io"
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
	"github.com/quic-go/quic-go"
)

func _StreamWriteBytePacket(conn quic.Connection, stream quic.Stream, data []byte) error {
	// Der Header, bestehend aus der Datenlänge wird hinzugefügt
	dataLength := len(data)
	dataLengthBytes := openkeyp2p.Uint64ToBytesLE(uint64(dataLength))
	finalDataBlock := append(dataLengthBytes, data...)

	// Der Schreibvorgang wird durchgeführt
	if _, err := stream.Write(finalDataBlock); err != nil {
		return err
	}

	// LOG
	localEndpointStr := getLocalIPAndHostFromConn(conn)
	remoteEndpointStr := getRemoteIPAndHostFromConn(conn)
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P_QUIC, "Data Packet writed, %d bytes %s -> %s", len(finalDataBlock), localEndpointStr, remoteEndpointStr)

	return nil
}

func _StreamReadBytePacket(conn quic.Connection, stream quic.Stream) ([]byte, error) {
	// Die Länge des Datensatzes wird ausgelesen
	dataLengthBytes := make([]byte, 8)
	n, err := io.ReadFull(stream, dataLengthBytes)
	if err != nil {
		// Wenn beim Lesen ein Fehler auftritt, behandeln wir diesen Fehler
		return nil, err
	}
	if n < 8 {
		return nil, fmt.Errorf("invalid data")
	}
	dataLength := openkeyp2p.BytesToUint64LE(dataLengthBytes)

	// Der Restliche Datensatz wird ausgelesen
	dataBytes := make([]byte, dataLength)
	_, err = io.ReadFull(stream, dataBytes)
	if err != nil {
		// Wenn beim Lesen ein Fehler auftritt, behandeln wir diesen Fehler
		return nil, err
	}

	// Log
	localEndpointStr := getLocalIPAndHostFromConn(conn)
	remoteEndpointStr := getRemoteIPAndHostFromConn(conn)
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P_QUIC, "Data Packet readed, %d bytes %s -> %s", dataLength+8, localEndpointStr, remoteEndpointStr)

	return dataBytes, nil
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
		quicConn:                conn,
		_recivedHelloBytePacket: recivedPacket,
		_sendHelloBytePacket:    helloPackage,
	}

	// Log
	localEndpointStr := getLocalIPAndHostFromConn(conn)
	remoteEndpointStr := getRemoteIPAndHostFromConn(conn)
	if isIncommingConnection {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Bidirectional streams were generated %s -> %s", localEndpointStr, remoteEndpointStr)
	} else {
		logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Bidirectional streams were generated %s -> %s", localEndpointStr, remoteEndpointStr)
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
	return _StreamWriteBytePacket(q.quicConn, q.outStream, byts)
}

func (q *QuicBidirectionalStream) ReadBytes() ([]byte, error) {
	return _StreamReadBytePacket(q.quicConn, q.inStream)
}
