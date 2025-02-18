package encoding

import (
	"encoding/base32"
	"fmt"
	"strings"
)

/*
DecodeAddressBase32StringWPrefix decodes an address from a Base32 encoded string with a predefined prefix.
It performs the following steps:
 1. Checks if the encoded string starts with the "okp2p" prefix.
 2. Removes the prefix from the string.
 3. Decodes the remaining string using Base32 with the Bech32 alphabet (without padding).
 4. Validates the decoded data and passes it to DecodeAddressByteSlice to extract address components.

Parameters:
  - encoded: a Base32 encoded string that starts with the "okp2p" prefix.

Returns:
  - typeFlag: the type flag indicating the address type.
  - ed25519Pub: the first public key (ED25519, 32 bytes).
  - sigPub: the second public key (size depends on the type flag).
  - port: optional port data (if included, 4 bytes).
  - err: an error if any step of the decoding process fails.
*/
func DecodeAddressBase32StringWPrefix(encoded string) (typeFlag uint8, ed25519Pub, sigPub, port []byte, err error) {
	const prefix = "okp2p"

	// Check if the encoded string has the correct prefix.
	if !strings.HasPrefix(encoded, prefix) {
		return 0, nil, nil, nil, fmt.Errorf("invalid prefix")
	}

	// Remove the prefix from the encoded string.
	encoded = encoded[len(prefix):]

	// Create a Base32 decoder using the Bech32 alphabet without padding.
	decoder := base32.NewEncoding(bech32Alphabet).WithPadding(base32.NoPadding)

	// Decode the Base32 string into a byte slice.
	data, err := decoder.DecodeString(encoded)
	if err != nil {
		return 0, nil, nil, nil, fmt.Errorf("error during Base32 decoding: %w", err)
	}

	// Ensure the decoded data is not empty.
	if len(data) < 1 {
		return 0, nil, nil, nil, fmt.Errorf("data too short")
	}

	// Decode the byte slice into address components using DecodeAddressByteSlice.
	return DecodeAddressByteSlice(data)
}
