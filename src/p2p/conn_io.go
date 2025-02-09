package p2p

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

// Write sends the provided byte slice over the connection by splitting it into frames.
// Enhanced error handling: avoids panic, retries on transient errors, and provides detailed error information.
func (c *FWConn) Write(b []byte) error {
	return c._Write(b)
}
func (c *FWConn) _Write(b []byte) error {
	// Generate Process ID
	procid, err := randomUint32()
	if err != nil {
		// Log the error and return it
		logging.LogError("Process ID generation failed: %v", err)
		return fmt.Errorf("process ID generation failed: %w", err)
	}

	logging.LogDebug("Generated Process ID: %d", procid)

	// Split the payload into frames
	frames := splitIntoFrames(b, fragmentSize, procid)
	logging.LogInfo("Split data into %d frames", len(frames))
	if len(frames) > int(c.maxFrameSize) {
		logging.LogError("To many number of frames (%d)", c.maxFrameSize)
		return fmt.Errorf("to many number of frames (%d)", c.maxFrameSize)
	}

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

// Read liest Daten von der Verbindung in Fragmenten.
func (c *FWConn) Read() ([]byte, error) {
	return c._Read()
}
func (c *FWConn) _Read() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	logging.LogInfo("Starting to read data")

	var (
		framesRead uint64
		dataBuffer = bytes.Buffer{}
	)

	for {
		if framesRead >= c.maxFrameSize {
			logging.LogError("Maximum number of frames (%d) exceeded", c.maxFrameSize)
			return nil, fmt.Errorf("maximum number of frames (%d) exceeded", c.maxFrameSize)
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
func (c *FWConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close the underlying network connection
	return c.conn.Close()
}
