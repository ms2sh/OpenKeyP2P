package main

import (
	"crypto/rand"
	"fmt"
	"log"

	blst "github.com/supranational/blst/bindings/go"
)

func main() {
	// Erzeuge einen 32-Byte Seed
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		log.Fatalf("Fehler bei der Zufallsgenerierung: %v", err)
	}

	// Erzeuge ed25519-Schlüsselpaar mit dem gleichen Seed
	ed25519Pub, _ := GenerateEd25519KeyPair(seed)
	_, _ = GenerateSecp256k1KeyPair(seed)

	// Erzeuge BLS-Schlüsselpaar mit dem gleichen Seed
	secretKey := blst.KeyGen(seed, nil)
	if secretKey == nil {
		log.Fatal("Fehler: SecretKey konnte nicht generiert werden")
	}

	// Leite den zugehörigen BLS Public Key ab
	blsPubKey := new(blst.P1Affine).From(secretKey)

	// Führe die Kodierung aus – dabei wird der ed25519 Public Key als kryponPubKey verwendet
	t := blsPubKey.Compress()
	encodedBls, err := EncodeAddressToBase32StringWPrefix(3, ed25519Pub, t, []byte("abcd"))
	if err != nil {
		log.Fatalf("Fehler: %v", err)
	}

	fmt.Println(encodedBls)

	ft, epkey, spkey, port, err := DecodeAddressBase32StringWPrefix(encodedBls)
	if err != nil {
		fmt.Println(err)
		return
	}

	tta, _ := EncodeAddressToBase32StringWPrefix(ft, epkey, spkey, port)
	fmt.Println(tta)
	fmt.Println(tta == encodedBls)
}
