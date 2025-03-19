package crypto

import (
	"crypto/sha256"
	"fmt"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"golang.org/x/crypto/sha3"
)

func ComputeSha256BitHash(data []byte) openkeyp2p.HashSlice {
	hash := sha256.Sum256(data) // Gibt ein [32]byte-Array zurück
	return hash[:]              // Umwandlung in []byte
}

func ComputeSha3_256BitHash(data []byte) openkeyp2p.HashSlice {
	hash := sha3.Sum256(data) // Gibt ein [32]byte-Array zurück
	return hash[:]            // Umwandlung in []byte
}

func ComputeHash(algo openkeyp2p.HashAlgorithm, data ...[]byte) (openkeyp2p.HashSlice, error) {
	byteBlock := make([]byte, 0)
	for _, item := range data {
		byteBlock = append(byteBlock, item...)
	}
	switch algo {
	case openkeyp2p.SHA_256:
		return ComputeSha256BitHash(byteBlock), nil
	case openkeyp2p.SHA3_256:
		return ComputeSha3_256BitHash(byteBlock), nil
	default:
		return nil, fmt.Errorf("unsupported hashing methode")
	}
}
