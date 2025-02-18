package encoding

import "encoding/base32"

/*
EncodeAddressToBase32StringWPrefix encodes an address into a Base32 string with a predefined prefix.
It first encodes the address components into a byte slice using EncodeAddressByteSlice,
then converts this byte slice into a Base32 string using a Bech32 alphabet (without padding),
and finally prepends the "okp2p" prefix to the resulting string.

Parameters:
  - typeFlag: a flag indicating the address type.
  - encryppkey: the first public key (e.g., ED25519, typically 32 bytes).
  - sigpubkey: the second public key (e.g., BLS12-381 or SECP256k1; size depends on the type).
  - port: optional port data (if provided, the first 4 bytes are used).

Returns:
  - A Base32 encoded string with the "okp2p" prefix.
  - An error if the encoding process fails.
*/
func EncodeAddressToBase32StringWPrefix(typeFlag uint8, encryppkey []byte, sigpubkey []byte, port []byte) (string, error) {
	// Encode the address components into a byte slice.
	addressByteSlice, err := EncodeAddressByteSlice(typeFlag, encryppkey, sigpubkey, port)
	if err != nil {
		// Return an error if the encoding of the byte slice fails.
		return "", err
	}

	// Create a Base32 encoder using the Bech32 alphabet without padding.
	encoder := base32.NewEncoding(bech32Alphabet).WithPadding(base32.NoPadding)
	// Encode the byte slice into a Base32 string.
	encoded := encoder.EncodeToString(addressByteSlice)

	// Prepend the "okp2p" prefix to the encoded string and return the result.
	return "okp2p" + encoded, nil
}
