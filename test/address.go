package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/ms2sh/OpenKeyP2P/src/crypto"
)

func main() {
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		log.Fatalf("Fehler bei der Zufallsgenerierung: %v", err)
	}

	priv, pub := crypto.GenerateRandomKeyPair(seed)
	addr, err := crypto.OpenKeyP2PAddressFromPublicKey(pub)
	if err != nil {
		panic(err)
	}

	dataHash := crypto.ComputeSha256BitHash([]byte("hallo welt"))

	sig, err := crypto.AddressSign(priv, dataHash)
	if err != nil {
		panic(err)
	}

	result, err := addr.VerifySignature(sig, dataHash)
	if err != nil {
		panic(err)
	}

	fmt.Println(addr.ToString(), result)

	decodedAddr, err := crypto.OpenKeyP2PAddressDecodeFromString(addr.ToString())
	if err != nil {
		panic(err)
	}
	fmt.Println(decodedAddr.ToString())

	_, err = crypto.OpenKeyP2PAddressDecodeFromByteSlice(decodedAddr.ToByteSlice())
	if err != nil {
		panic(err)
	}
}
