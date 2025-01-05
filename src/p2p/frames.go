package p2p

import (
	"encoding/binary"
	"hash/crc32"
)

// _Header represents the header of a frame.
// It contains metadata about the frame such as data length, checksum,
// process ID, frame number, and a flag indicating if it's the last frame.
type _Header struct {
	DataLength uint32 // 4 Bytes: Size of the body data
	Checksum   uint32 // 4 Bytes: Checksum of the body data
	ProcessId  uint32 // 4 Bytes: Process ID
	FrameNo    uint64 // 8 Bytes: Frame number (uint64 -> up to 10 PB)
	LastFrame  bool   // 1 Byte: Is this the last frame?
}

// _Frame represents a data frame consisting of a header and a body.
type _Frame struct {
	Header *_Header // Pointer to the frame's header
	Body   []byte   // Byte slice containing the actual data
}

// GetBodyBytes returns the relevant data bytes of the frame.
// It slices the body according to the DataLength specified in the header.
func (f *_Frame) GetBodyBytes() []byte {
	return f.Body[:f.Header.DataLength]
}

// splitIntoFrames splits the given data into multiple frames.
// Each frame has a maximum size of frameSize bytes and is associated with processId.
func splitIntoFrames(data []byte, frameSize uint32, processId uint32) []*_Frame {
	const headerSize = 21 // 21 Bytes (4+4+4+8+1)

	// Ensure that the frame size is larger than the header size
	if frameSize <= headerSize {
		// Frame size must be larger than the header,
		// otherwise, no body can fit.
		panic("frameSize must be larger than header size")
	}

	// Calculate the maximum number of data bytes for the body per frame.
	bodySize := frameSize - headerSize

	var frames []*_Frame // Slice to store the generated frames

	var frameNo uint64 = 0           // Initialize frame number
	var offset uint64 = 0            // Starting position in the data stream
	totalLength := uint64(len(data)) // Total length of the data

	for offset < totalLength || totalLength == 0 {
		// If no data is present (totalLength == 0), create exactly one frame with an empty body.
		if totalLength == 0 {
			singleFrame := &_Frame{
				Header: &_Header{
					DataLength: 0,         // No data
					Checksum:   0,         // Checksum is 0 for empty body
					ProcessId:  processId, // Process ID
					FrameNo:    frameNo,   // Frame number
					LastFrame:  true,      // Mark as last frame
				},
				Body: []byte{}, // Empty body
			}
			frames = append(frames, singleFrame) // Add the frame to the list
			break                                // Exit the loop
		}

		// Calculate how many bytes are remaining
		remaining := totalLength - offset

		// Determine the size of the next chunk
		var chunkSize uint32
		if remaining < uint64(bodySize) {
			chunkSize = uint32(remaining) // Last chunk, smaller than bodySize
		} else {
			chunkSize = bodySize // Full bodySize
		}

		// Extract the chunk data from the total dataset
		body := data[offset : offset+uint64(chunkSize)]

		// Calculate the checksum (CRC32 from the standard package)
		checksum := crc32.ChecksumIEEE(body)

		// Determine if this is the last frame
		lastFrame := (offset+uint64(chunkSize) >= totalLength)

		// Build the header for the current frame
		header := _Header{
			DataLength: chunkSize, // Size of the data in the body
			Checksum:   checksum,  // Calculated checksum
			ProcessId:  processId, // Process ID
			FrameNo:    frameNo,   // Current frame number
			LastFrame:  lastFrame, // Flag indicating if it's the last frame
		}

		// Create the frame with the header and the body
		frame := &_Frame{
			Header: &header,
			Body:   body,
		}

		// Append the frame to the list of frames
		frames = append(frames, frame)

		// Update the offset and frame number for the next iteration
		offset += uint64(chunkSize)
		frameNo++

		// If the last frame is reached, exit the loop
		if lastFrame {
			break
		}
	}

	return frames // Return the list of frames
}

// frameToBytes converts a _Frame object into a byte slice.
// The resulting byte slice consists of the header followed by the body.
func frameToBytes(frame *_Frame) []byte {
	// We know that our header is 21 bytes long:
	// - 4 Bytes DataLength
	// - 4 Bytes Checksum
	// - 4 Bytes ProcessId
	// - 8 Bytes FrameNo
	// - 1 Byte LastFrame
	headerSize := 21
	headerBuf := make([]byte, headerSize) // Buffer for the header

	// 1) DataLength (4 Bytes, BigEndian)
	binary.BigEndian.PutUint32(headerBuf[0:4], frame.Header.DataLength)

	// 2) Checksum (4 Bytes, BigEndian)
	binary.BigEndian.PutUint32(headerBuf[4:8], frame.Header.Checksum)

	// 3) ProcessId (4 Bytes, BigEndian)
	binary.BigEndian.PutUint32(headerBuf[8:12], frame.Header.ProcessId)

	// 4) FrameNo (8 Bytes, BigEndian)
	binary.BigEndian.PutUint64(headerBuf[12:20], frame.Header.FrameNo)

	// 5) LastFrame (1 Byte -> 0 or 1)
	if frame.Header.LastFrame {
		headerBuf[20] = 1 // Last frame
	} else {
		headerBuf[20] = 0 // Not the last frame
	}

	// Append the body
	// Complete frame = Header + Body
	result := append(headerBuf, frame.Body...)

	return result
}
