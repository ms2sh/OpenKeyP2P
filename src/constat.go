package openkeyp2p

import "encoding/base32"

const (
	Version                     OpenKeyP2PVesion  = OpenKeyP2PVesion(100000001)
	Prefix                      OpenKeyP2PPrefix  = OpenKeyP2PPrefix("okp2p")
	Base32DefaultBase32Alphabet Base32Alphabet    = Base32Alphabet("qpzry9x8gf2tvdw0s3jn54khce6mua7l")
	SHA_256                     HashAlgorithm     = 1
	SHA3_256                    HashAlgorithm     = 2
	Type_Ed25519                OpenKeyP2PKeyType = 0

	// Hash Methods
	DEFAULT_HASH_METHODE_256BIT = SHA_256

	// Log Levels
	LOG_LEVEL_P2P_QUIC LogLevel = 1
	LOG_LEVEL_P2P      LogLevel = 2
)

var (
	Base32Encoding    = base32.NewEncoding(string(Base32DefaultBase32Alphabet)).WithPadding(base32.NoPadding)
	SUPPORTED_VERSION = []OpenKeyP2PVesion{Version}
)
