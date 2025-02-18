package encoding

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

/*
DecodeAddressByteSlice decodes an encoded address byte slice and extracts the type flag, public keys, and optional port.
It validates the integrity of the data using a 5-byte checksum.

Parameters:
  - byteSliceAddress: the byte slice representing the encoded address.

Returns:
  - typeFlag: the type flag (first byte) of the encoded address.
  - ed25519Pub: the first public key (ED25519, 32 bytes).
  - sigPub: the second public key (size depends on the type flag).
  - port: optional port data (if included, 4 bytes).
  - err: an error if the data is invalid or if the checksum does not match.
*/
func DecodeAddressByteSlice(byteSliceAddress []byte) (typeFlag uint8, ed25519Pub, sigPub, port []byte, err error) {
	// The first byte represents the type flag
	typeFlag = byteSliceAddress[0]

	// The ED25519 public key length is always 32 bytes
	ed25519Len := 32
	var sigLen, expectedLen int
	var hasPort bool

	// Determine the length of the second public key and whether a port is included based on the type flag
	switch typeFlag {
	case 1:
		// ED25519 + BLS12-381 (without port)
		sigLen = 48
		hasPort = false
		expectedLen = 1 + ed25519Len + sigLen + 5 // type flag + ED25519 key + second key + checksum
	case 2:
		// ED25519 + SECP256k1 (without port)
		sigLen = 33
		hasPort = false
		expectedLen = 1 + ed25519Len + sigLen + 5
	case 3:
		// ED25519 + BLS12-381 (with port)
		sigLen = 48
		hasPort = true
		expectedLen = 1 + ed25519Len + sigLen + 4 + 5 // extra 4 bytes for port
	case 4:
		// ED25519 + SECP256k1 (with port)
		sigLen = 33
		hasPort = true
		expectedLen = 1 + ed25519Len + sigLen + 4 + 5
	default:
		return 0, nil, nil, nil, fmt.Errorf("unknown type flag: %d", typeFlag)
	}

	// Check if the provided byte slice has the expected length
	if len(byteSliceAddress) != expectedLen {
		return 0, nil, nil, nil, fmt.Errorf("invalid data length: expected %d, got %d", expectedLen, len(byteSliceAddress))
	}

	// The last 5 bytes represent the checksum
	payload := byteSliceAddress[:len(byteSliceAddress)-5]
	checksum := byteSliceAddress[len(byteSliceAddress)-5:]

	// Validate the checksum: compute SHA-256 over the payload and compare the first 5 bytes
	hash := sha256.Sum256(payload)
	if !bytes.Equal(checksum, hash[:5]) {
		return 0, nil, nil, nil, fmt.Errorf("checksum does not match")
	}

	// Extract the ED25519 public key (bytes 1 to 32)
	ed25519Pub = payload[1 : 1+ed25519Len]

	// Extract the second public key (depending on the type flag)
	sigPub = payload[1+ed25519Len : 1+ed25519Len+sigLen]

	// If a port is included, extract the 4-byte port information
	if hasPort {
		port = payload[1+ed25519Len+sigLen : 1+ed25519Len+sigLen+4]
	}

	return typeFlag, ed25519Pub, sigPub, port, nil
}
