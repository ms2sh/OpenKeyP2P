package crypto

import (
	"crypto/ed25519"
	"fmt"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
)

type OpenKeyP2PSignature []byte

func (o OpenKeyP2PSignature) GetRawSignature() []byte {
	return o
}

func (o OpenKeyP2PSignature) VerifySignatureHashDigest(dataHash openkeyp2p.HashSlice, pubKey openkeyp2p.OpenKeyP2PPublicKey) (bool, error) {
	if dataHash == nil {
		return false, fmt.Errorf("invalid hash digest")
	}
	if pubKey == nil {
		return false, fmt.Errorf("invalid public key")
	}
	if o == nil {
		return false, fmt.Errorf("nill signature")
	}
	return ed25519.Verify(ed25519.PublicKey(pubKey), dataHash, o), nil
}
