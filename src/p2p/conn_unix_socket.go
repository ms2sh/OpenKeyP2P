package p2p

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"
)

type UnixSocketAckConn struct {
	conn           net.Conn
	readChan       chan []byte
	writeChan      chan []byte
	errorChan      chan error
	disconnectChan chan struct{}
	closeChan      chan struct{}
	closed         bool
	mutex          sync.Mutex
}

// NewUnixSocketConn creates a new UnixSocketAckConn
func upgradeConnToUnixSocketAckConn(conn net.Conn) (*UnixSocketAckConn, error) {
	usc := &UnixSocketAckConn{
		conn:           conn,
		readChan:       make(chan []byte, 1024),
		writeChan:      make(chan []byte, 10),
		errorChan:      make(chan error, 1),
		disconnectChan: make(chan struct{}, 1),
		closeChan:      make(chan struct{}),
	}

	go usc.readLoop()
	go usc.writeLoop()

	return usc, nil
}

func (usc *UnixSocketAckConn) calculateChecksum(data []byte) uint64 {
	hash := sha256.Sum256(data)
	return binary.BigEndian.Uint64(hash[:8])
}

func (usc *UnixSocketAckConn) setError(err error) {
	select {
	case usc.errorChan <- err:
	default:
	}
	usc.signalDisconnect()
}

func (usc *UnixSocketAckConn) signalDisconnect() {
	select {
	case usc.disconnectChan <- struct{}{}:
	default:
	}
}

func (usc *UnixSocketAckConn) readLoop() {
	buf := make([]byte, unixSocketMaxFrameSize)
	for {
		select {
		case <-usc.closeChan:
			return
		default:
			n, err := usc.conn.Read(buf)
			if err != nil {
				usc.setError(err)
				close(usc.readChan)
				return
			}

			if n < unixSocketAckSize || buf[0] != 1 {
				usc.setError(errors.New("invalid ACK received"))
				return
			}

			//checksum := binary.BigEndian.Uint64(buf[1:9])

			data := make([]byte, n)
			copy(data, buf[:n])
			usc.readChan <- data
		}
	}
}

func (usc *UnixSocketAckConn) writeLoop() {
	for {
		select {
		case <-usc.closeChan:
			return
		case data := <-usc.writeChan:
			frameSize := len(data)
			if frameSize > unixSocketFramePayloadSize {
				usc.setError(errors.New("data exceeds maximum frame payload size"))
				return
			}

			// Create frame header
			frame := make([]byte, unixSocketFrameHeaderSize+frameSize)
			checksum := usc.calculateChecksum(data)
			binary.BigEndian.PutUint64(frame[0:8], checksum)
			binary.BigEndian.PutUint32(frame[8:12], uint32(frameSize))
			copy(frame[12:], data)

			// Write frame
			if _, err := usc.conn.Write(frame); err != nil {
				usc.setError(err)
				return
			}

			// Wait for ACK
			ack := <-usc.readChan
			if len(ack) < unixSocketAckSize || ack[0] != 1 {
				usc.setError(errors.New("invalid ACK received"))
				return
			}

			ackChecksum := binary.BigEndian.Uint64(ack[1:9])
			if ackChecksum != checksum {
				usc.setError(errors.New("ACK validation failed"))
				return
			}
		}
	}
}

// Write sends data in frames with acknowledgment validation
func (usc *UnixSocketAckConn) Write(data []byte) (int, error) {
	usc.mutex.Lock()
	defer usc.mutex.Unlock()
	if usc.closed {
		return 0, errors.New("connection closed")
	}

	select {
	case err := <-usc.errorChan:
		return 0, err
	default:
	}

	usc.writeChan <- data
	return len(data), nil
}

// Read is not used for ACK but reserved for payloads
func (usc *UnixSocketAckConn) Read(buf []byte) (int, error) {
	usc.mutex.Lock()
	defer usc.mutex.Unlock()
	if usc.closed {
		return 0, errors.New("connection closed")
	}
	select {
	case err := <-usc.errorChan:
		return 0, err
	default:
	}
	n := copy(buf, <-usc.readChan)
	return n, nil
}

// Close closes the connection
func (usc *UnixSocketAckConn) Close() error {
	usc.mutex.Lock()
	defer usc.mutex.Unlock()
	if usc.closed {
		return errors.New("connection already closed")
	}
	close(usc.closeChan)
	usc.signalDisconnect()
	usc.closed = true
	return usc.conn.Close()
}

// IsDisconnected checks if the connection is disconnected
func (usc *UnixSocketAckConn) IsDisconnected() bool {
	select {
	case <-usc.disconnectChan:
		return true
	default:
		return false
	}
}

// LocalAddr returns the local network address
func (usc *UnixSocketAckConn) LocalAddr() net.Addr {
	return usc.conn.LocalAddr()
}

// RemoteAddr returns the remote network address
func (usc *UnixSocketAckConn) RemoteAddr() net.Addr {
	return usc.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated with the connection
func (usc *UnixSocketAckConn) SetDeadline(t time.Time) error {
	return usc.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
func (usc *UnixSocketAckConn) SetReadDeadline(t time.Time) error {
	return usc.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
func (usc *UnixSocketAckConn) SetWriteDeadline(t time.Time) error {
	return usc.conn.SetWriteDeadline(t)
}
