package ed25519

import (
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
)

func GenerateSecp256k1KeyPair(seed []byte) (*btcec.PrivateKey, *btcec.PublicKey) {
	if len(seed) != 32 {
		log.Fatalf("Der Seed muss 32 Bytes lang sein, erhalten: %d", len(seed))
	}
	// btcec/v2 bietet PrivKeyFromBytes, das direkt ein Schlüsselpaar aus den Bytes erzeugt.
	priv, pub := btcec.PrivKeyFromBytes(seed)
	return priv, pub
}
