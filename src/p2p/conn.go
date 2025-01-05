package p2p

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

const (
	FragmentSize    = 1024
	maxWriteRetries = 3
	retryBaseDelay  = 1 * time.Second
	maxFrames       = 1048576
)

// Conn wraps a net.Conn with a Mutex and a Condition to synchronize read/write operations.
type Conn struct {
	conn    net.Conn   // Underlying network connection
	mu      sync.Mutex // Mutex to protect concurrent access
	cond    *sync.Cond // Condition variable for synchronization
	reading bool       // Flag to check if a read operation is active
	writing bool       // Flag to check if a write operation is active
}

// UpgradeConn creates and returns a new instance of Conn, wrapping the provided net.Conn.
func UpgradeConn(conn net.Conn) *Conn {
	c := &Conn{conn: conn}
	c.cond = sync.NewCond(&c.mu) // Initialize the Condition variable
	return c
}

// Write sends the provided byte slice over the connection by splitting it into frames.
// Enhanced error handling: avoids panic, retries on transient errors, and provides detailed error information.
func (c *Conn) Write(b []byte) error {
	// Generate Process ID
	procid, err := randomUint32()
	if err != nil {
		// Log the error and return it
		logging.LogError("Process ID generation failed: %v", err)
		return fmt.Errorf("process ID generation failed: %w", err)
	}

	logging.LogDebug("Generated Process ID: %d", procid)

	// Split the payload into frames
	frames := splitIntoFrames(b, FragmentSize, procid)
	logging.LogInfo("Split data into %d frames", len(frames))

	// Transmit individual frames
	for _, frame := range frames {
		var writeErr error
		for attempt := 1; attempt <= maxWriteRetries; attempt++ {
			logging.LogDebug("Writing frame %d (Attempt %d)", frame.Header.FrameNo, attempt)
			writeErr = c.writeFrame(frame)
			if writeErr == nil {
				logging.LogInfo("Successfully wrote frame %d", frame.Header.FrameNo)
				break // Successfully written, proceed to next frame
			}

			logging.LogError("Error writing frame %d: %v", frame.Header.FrameNo, writeErr)

			// Check if the error is temporary (e.g., network issues)
			if isTemporaryError(writeErr) && attempt < maxWriteRetries {
				waitDuration := time.Duration(attempt) * retryBaseDelay
				logging.LogInfo("Temporary error, retrying after %v", waitDuration)
				time.Sleep(waitDuration)
				continue
			}

			break
		}

		if writeErr != nil {
			// Return the error with context
			return writeErr
		}
	}

	logging.LogInfo("All frames written successfully")
	return nil
}

// writeFrame sends a single frame over the connection.
func (c *Conn) writeFrame(frame *_Frame) error {
	bytedFrame := frameToBytes(frame)
	total := len(bytedFrame)
	written := 0

	logging.LogDebug("Starting to write frame %d (%d bytes)", frame.Header.FrameNo, total)

	for written < total {
		n, err := c.conn.Write(bytedFrame[written:])
		if err != nil {
			logging.LogError("Write error on frame %d: %v", frame.Header.FrameNo, err)
			return fmt.Errorf("write error: %w", err)
		}
		if n == 0 {
			logging.LogError("Connection might be broken: wrote 0 bytes for frame %d", frame.Header.FrameNo)
			return fmt.Errorf("connection might be broken: wrote 0 bytes")
		}
		written += n
		logging.LogDebug("Wrote %d/%d bytes for frame %d", written, total, frame.Header.FrameNo)
	}
	logging.LogDebug("Finished writing frame %d", frame.Header.FrameNo)
	return nil
}

// readFrame liest ein einzelnes Frame von der Verbindung.
func (c *Conn) readFrame() (*_Frame, error) {
	headerSize := 21
	headerBuf := make([]byte, headerSize)

	logging.LogDebug("Reading frame header")

	// Read the 21 header bytes
	n, err := io.ReadFull(c.conn, headerBuf)
	if err != nil {
		if err == io.EOF {
			logging.LogError("Connection closed while reading header: %v", err)
			return nil, fmt.Errorf("connection closed while reading header: %w", err)
		}
		logging.LogError("Failed to read frame header (%d bytes read): %v", n, err)
		return nil, fmt.Errorf("failed to read frame header (%d bytes read): %w", n, err)
	}

	logging.LogDebug("Read frame header successfully")

	// Parse Header
	h := &_Header{}
	h.DataLength = binary.BigEndian.Uint32(headerBuf[0:4])
	h.Checksum = binary.BigEndian.Uint32(headerBuf[4:8])
	h.ProcessId = binary.BigEndian.Uint32(headerBuf[8:12])
	h.FrameNo = binary.BigEndian.Uint64(headerBuf[12:20])
	h.LastFrame = headerBuf[20] == 1

	logging.LogDebug("Parsed header: %+v", h)

	// Read the body if present (DataLength can be 0)
	var bodyBuf []byte
	if h.DataLength > 0 {
		bodyBuf = make([]byte, h.DataLength)

		logging.LogDebug("Reading frame body (%d bytes)", h.DataLength)

		n, err := io.ReadFull(c.conn, bodyBuf)
		if err != nil {
			if err == io.EOF {
				logging.LogError("Connection closed while reading body: %v", err)
				return nil, fmt.Errorf("connection closed while reading body: %w", err)
			}
			logging.LogError("Failed to read frame body (%d bytes read): %v", n, err)
			return nil, fmt.Errorf("failed to read frame body (%d bytes read): %w", n, err)
		}

		logging.LogDebug("Read frame body successfully")

		// Verify checksum
		csum := crc32.ChecksumIEEE(bodyBuf)
		if csum != h.Checksum {
			logging.LogError("Invalid checksum: expected %d, got %d", h.Checksum, csum)
			return nil, fmt.Errorf("invalid checksum: expected %d, got %d", h.Checksum, csum)
		}

		logging.LogDebug("Checksum verified successfully")
	}

	// Construct the frame
	frame := &_Frame{
		Header: h,
		Body:   bodyBuf,
	}

	logging.LogInfo("Read frame %d successfully", h.FrameNo)

	return frame, nil
}

// Read liest Daten von der Verbindung in Fragmenten.
func (c *Conn) Read() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	logging.LogInfo("Starting to read data")

	var (
		framesRead uint64
		dataBuffer = bytes.Buffer{}
	)

	for {
		if framesRead >= maxFrames {
			logging.LogError("Maximum number of frames (%d) exceeded", maxFrames)
			return nil, fmt.Errorf("maximum number of frames (%d) exceeded", maxFrames)
		}

		logging.LogDebug("Reading frame number %d", framesRead)

		frame, err := c.readFrame()
		if err != nil {
			// Prüfe, ob der Fehler ein EOF-Fehler ist
			if errors.Is(err, io.EOF) {
				if framesRead == 0 {
					logging.LogInfo("No frames read, returning EOF")
					return nil, io.EOF
				}
				logging.LogInfo("EOF reached after reading %d frames", framesRead)
				return dataBuffer.Bytes(), nil // Rückgabe der bisher gelesenen Daten
			}
			logging.LogError("Error reading frame: %v", err)
			return nil, err
		}

		// Überprüfe, ob die Frame-Nummer wie erwartet ist
		if frame.Header.FrameNo != framesRead {
			logging.LogError("Invalid frame number: expected %d, got %d", framesRead, frame.Header.FrameNo)
			return nil, fmt.Errorf("invalid frame number: expected %d, got %d", framesRead, frame.Header.FrameNo)
		}

		framesRead++
		if _, err := dataBuffer.Write(frame.GetBodyBytes()); err != nil {
			logging.LogError("Error writing frame data to buffer: %v", err)
			return nil, err
		}

		logging.LogDebug("Frame %d written to buffer", frame.Header.FrameNo)

		// Prüfe, ob dies das letzte Frame ist
		if frame.Header.LastFrame {
			logging.LogInfo("Last frame %d received", frame.Header.FrameNo)
			break
		}
	}

	logging.LogInfo("All frames read successfully")
	return dataBuffer.Bytes(), nil
}

// Close closes the connection, releasing any resources.
func (c *Conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close the underlying network connection
	return c.conn.Close()
}
