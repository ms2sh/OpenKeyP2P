package encoding

import (
	"bytes"
	"crypto/sha256"
)

/*
EncodeAddressByteSlice encodes the provided type flag, public keys, and optional port into a byte slice.
It appends a 5-byte checksum at the end of the resulting byte slice.

Parameters:
  - typeFlag: a flag indicating the type of address.
  - encryppkey: the first public key (e.g., ED25519, typically 32 bytes).
  - sigpubkey: the second public key (e.g., BLS12-381 or SECP256k1; size depends on type).
  - port: optional port data (if provided, the first 4 bytes are used).

Returns:
  - A byte slice representing the encoded address.
  - An error, if any occurred during encoding (currently always nil).
*/
func EncodeAddressByteSlice(typeFlag uint8, encryppkey []byte, sigpubkey []byte, port []byte) ([]byte, error) {
	var buffer bytes.Buffer

	// Write the type flag (1 byte)
	buffer.WriteByte(typeFlag)

	// Append the first public key (e.g., ED25519, expected to be 32 bytes)
	buffer.Write(encryppkey)

	// Append the second public key:
	// For BLS12-381 this is 48 bytes (G1 compressed), for SECP256k1 it might be 33 bytes
	buffer.Write(sigpubkey)

	// If a port is provided, append the first 4 bytes of the port information
	if port != nil {
		buffer.Write(port[0:4])
	}

	// Compute the checksum: use the first 5 bytes of the SHA-256 hash of the current buffer
	hash := sha256.Sum256(buffer.Bytes())
	checksum := hash[:5]
	buffer.Write(checksum)

	// Return the complete byte slice
	return buffer.Bytes(), nil
}
