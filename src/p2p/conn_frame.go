package p2p

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"

	"github.com/ms2sh/OpenKeyP2P/src/logging"
)

// readFrame liest ein einzelnes Frame von der Verbindung.
func (c *FWConn) readFrame() (*_Frame, error) {
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

// writeFrame sends a single frame over the connection.
func (c *FWConn) writeFrame(frame *_Frame) error {
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
