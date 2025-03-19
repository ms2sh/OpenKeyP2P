package crypto

import (
	"crypto/ed25519"

	"filippo.io/edwards25519"
)

func Ed25519ToCurve25519PublicKey(edPub ed25519.PublicKey) ([]byte, error) {
	// Ed25519 Public Key in Edwards25519-Form
	var edPubPoint edwards25519.Point
	if _, err := edPubPoint.SetBytes(edPub); err != nil {
		return nil, err
	}

	// In Montgomery-Form (Curve25519) konvertieren
	curvePub := edPubPoint.BytesMontgomery()
	return curvePub[:], nil
}

func GenerateKeyPairFromSeed(seed []byte) (ed25519.PrivateKey, ed25519.PublicKey) {
	// Privaten Schlüssel aus dem Seed generieren
	privKey := ed25519.NewKeyFromSeed(seed)

	// Öffentlichen Schlüssel extrahieren
	pubKey := privKey.Public().(ed25519.PublicKey)

	return privKey, pubKey
}

func GenerateRandomKeyPair(seed []byte) (ed25519.PrivateKey, ed25519.PublicKey) {
	// Privaten Schlüssel aus dem Seed generieren
	privKey := ed25519.NewKeyFromSeed(seed)

	// Öffentlichen Schlüssel extrahieren
	pubKey := privKey.Public().(ed25519.PublicKey)

	return privKey, pubKey
}

func PrivateKeyToPublicKey(privKey ed25519.PrivateKey) ed25519.PublicKey {
	pubKey := privKey.Public().(ed25519.PublicKey)
	return pubKey
}
