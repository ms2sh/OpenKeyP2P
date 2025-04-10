package p2p

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/ms2sh/OpenKeyP2P/src/logging"
	"github.com/quic-go/quic-go"
)

func _StreamWriteBytePacket(stream quic.Stream, data []byte, localSocketEp NodeP2PSocketAddress, remoteSocketEp NodeP2PSocketAddress, connCtxCancel context.CancelCauseFunc) error {
	// Der Header, bestehend aus der Datenlänge wird hinzugefügt
	dataLength := len(data)
	dataLengthBytes := openkeyp2p.Uint64ToBytesLE(uint64(dataLength))
	finalDataBlock := append(dataLengthBytes, data...)

	// Der Schreibvorgang wird durchgeführt
	if _, err := stream.Write(finalDataBlock); err != nil {
		var (
			netErr    net.Error
			streamErr *quic.StreamError
			connErr   *quic.ApplicationError
		)

		switch {
		case errors.As(err, &netErr) && netErr.Timeout():
			return fmt.Errorf("network timeout: %w", err)
		case errors.Is(err, net.ErrClosed):
			return fmt.Errorf("connection closed: %w", err)
		case errors.As(err, &streamErr):
			return fmt.Errorf("QUIC stream error (code %d): %w", streamErr.ErrorCode, err)
		case errors.As(err, &connErr):
			return fmt.Errorf("QUIC connection error (code %d): %w", connErr.ErrorCode, err)
		default:
			return fmt.Errorf("write failed: %w", err)
		}
	}

	// LOG
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P_QUIC, "Data Packet writed, %d bytes %s -> %s", len(finalDataBlock), localSocketEp, remoteSocketEp)

	return nil
}

func _StreamReadBytePacket(stream quic.Stream, localSocketEp NodeP2PSocketAddress, remoteSocketEp NodeP2PSocketAddress, connCtxCancel context.CancelCauseFunc) ([]byte, error) {
	// Die Länge des Datensatzes wird ausgelesen
	dataLengthBytes := make([]byte, 8)
	n, err := io.ReadFull(stream, dataLengthBytes)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, fmt.Errorf("stream ended prematurely (expected %d bytes): %w", 8, err)
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf("read timeout after %d bytes: %w", 8, err)
		}
		return nil, fmt.Errorf("failed to read payload (expected %d bytes): %w", 8, err)
	}
	if n < 8 {
		return nil, fmt.Errorf("invalid data")
	}
	dataLength := openkeyp2p.BytesToUint64LE(dataLengthBytes)

	// Der Restliche Datensatz wird ausgelesen
	dataBytes := make([]byte, dataLength)
	_, err = io.ReadFull(stream, dataBytes)
	if err != nil {
		return nil, err
	}

	// Log
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P_QUIC, "Data Packet readed, %d bytes %s -> %s", dataLength+8, localSocketEp, remoteSocketEp)

	return dataBytes, nil
}

func _TryOpenQuicBidirectionalStream(isIncommingConnection bool, conn quic.Connection, helloPackage []byte, localSocketEp NodeP2PSocketAddress, remoteSocketEp NodeP2PSocketAddress, connCtx context.Context, connCtxCancel context.CancelCauseFunc) (*QuicBidirectionalStream, error) {
	// Es wird selektiert, ob es sich um eine eingehende oder um eine ausgehende Verbindung handelt
	var inStream quic.Stream
	var outStream quic.Stream
	var streamErr error
	var recivedPacket []byte
	if isIncommingConnection {
		// Es wird ein ausgehe4nder Stream geöffnet
		outStream, streamErr = conn.OpenStreamSync(connCtx)
		if streamErr != nil {
			errmsg := fmt.Errorf("failed to send hello: %w", streamErr)
			connCtxCancel(errmsg)
			return nil, streamErr
		}

		// Das Hello Stream Package wird an die Gegenseite übertragen
		if err := _StreamWriteBytePacket(outStream, helloPackage, localSocketEp, remoteSocketEp, connCtxCancel); err != nil {
			errmsg := fmt.Errorf("failed to send hello: %w", streamErr)
			connCtxCancel(errmsg)
			return nil, err
		}

		// Es wird auf einen eingehenden Stream gewartet
		inStream, streamErr = conn.AcceptStream(connCtx)
		if streamErr != nil {
			errmsg := fmt.Errorf("failed to send hello: %w", streamErr)
			connCtxCancel(errmsg)
			return nil, streamErr
		}

		// Es wird auf das Eingehende Hello Stream Package gewartet
		var err error
		recivedPacket, err = _StreamReadBytePacket(inStream, localSocketEp, remoteSocketEp, connCtxCancel)
		if err != nil {
			connCtxCancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, err
		}
	} else {
		// Es wird auf einen eingehenden Stream gewartet
		inStream, streamErr = conn.AcceptStream(connCtx)
		if streamErr != nil {
			connCtxCancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, streamErr
		}

		// Es wird auf das HelloPackage der gegenseite gewartet
		var err error
		recivedPacket, err = _StreamReadBytePacket(inStream, localSocketEp, remoteSocketEp, connCtxCancel)
		if err != nil {
			connCtxCancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, err
		}

		// Es wird eine ausgehende Verbindung aufgebaut
		outStream, streamErr = conn.OpenStreamSync(connCtx)
		if streamErr != nil {
			connCtxCancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, streamErr
		}

		// Das Hello Stream Package wird an die Gegenseite übertragen
		if err := _StreamWriteBytePacket(outStream, helloPackage, localSocketEp, remoteSocketEp, connCtxCancel); err != nil {
			connCtxCancel(fmt.Errorf("failed to send hello: %w", streamErr))
			return nil, err
		}
	}

	// Das Finale Objekt wird erzeugt
	finalObject := &QuicBidirectionalStream{
		inStream:                inStream,
		outStream:               outStream,
		lock:                    new(sync.Mutex),
		ctx:                     connCtx,
		ctxCancle:               connCtxCancel,
		quicConn:                conn,
		readMutex:               new(sync.Mutex),
		writeMutex:              new(sync.Mutex),
		_recivedHelloBytePacket: recivedPacket,
		_sendHelloBytePacket:    helloPackage,
		_localSocketEp:          localSocketEp,
		_remoteSocketEp:         remoteSocketEp,
	}

	// Log
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P, "Bidirectional streams were generated %s -> %s", localSocketEp, remoteSocketEp)

	// Das Objekt wird zurückgegeben
	return finalObject, nil
}

func (q *QuicBidirectionalStream) Close() {
	// Falls bereits geschlossen, nichts tun
	if q.ctx.Err() != nil {
		return
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	// Streams schließen
	q.inStream.Close()
	q.outStream.Close()

	// Kontext beenden
	q.ctxCancle(fmt.Errorf("stream closed"))

	// LOG
	logging.LogDebug(openkeyp2p.LOG_LEVEL_P2P_QUIC, "Data Packet writed, %d bytes %s -> %s", q._localSocketEp, q._remoteSocketEp)
}

func (q *QuicBidirectionalStream) WriteBytes(byts []byte) error {
	if err := q.ctx.Err(); err != nil {
		return fmt.Errorf("context closed")
	}

	q.writeMutex.Lock()
	defer q.writeMutex.Unlock()

	err := _StreamWriteBytePacket(q.outStream, byts, q._localSocketEp, q._remoteSocketEp, q.ctxCancle)
	if err != nil {
		return err
	}

	return nil
}

func (q *QuicBidirectionalStream) ReadBytes() ([]byte, error) {
	if err := q.ctx.Err(); err != nil {
		return nil, fmt.Errorf("context closed")
	}

	q.readMutex.Lock()
	defer q.readMutex.Unlock()

	data, err := _StreamReadBytePacket(q.inStream, q._localSocketEp, q._remoteSocketEp, q.ctxCancle)
	if err != nil {
		return nil, err
	}

	return data, nil
}
