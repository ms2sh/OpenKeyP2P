package address

import (
	"crypto/ed25519"
	"log"
)

func GenerateEd25519KeyPair(seed []byte) (ed25519.PublicKey, ed25519.PrivateKey) {
	if len(seed) != ed25519.SeedSize {
		log.Fatalf("Der Seed muss %d Bytes lang sein", ed25519.SeedSize)
	}
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	return publicKey, privateKey
}
