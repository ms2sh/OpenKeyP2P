package bls12

import (
	"crypto/ed25519"
	"log"

	blst "github.com/supranational/blst/bindings/go"
)

func GenerateBLS12KeyPair(seed []byte) (*blst.P1Affine, blst.SecretKey) {
	if len(seed) != ed25519.SeedSize {
		log.Fatalf("Der Seed muss %d Bytes lang sein", ed25519.SeedSize)
	}

	// Erzeuge BLS-Schlüsselpaar mit dem gleichen Seed
	secretKey := blst.KeyGen(seed, nil)
	if secretKey == nil {
		log.Fatal("Fehler: SecretKey konnte nicht generiert werden")
	}

	// Leite den zugehörigen BLS Public Key ab
	blsPubKey := new(blst.P1Affine).From(secretKey)

	return blsPubKey, *secretKey
}
