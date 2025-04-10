package p2p

import (
	"crypto/rand"
)

func _SignByteSlice(bslice []byte) ([]byte, error) {
	return nil, nil
}

func _BuildRandomVIdValue() (NodeP2PConnectionValidationId, error) {
	return nil, nil
}

func _SignSteamPacketWSigPacket(packet interface{}) ([]byte, error) {
	return []byte("signature"), nil
}

func _GetSignerPublicKey() NodePublicSignatureKey {
	return NodePublicSignatureKey{}
}

func _GetEncryptionPublicKey() NodePublicEncryptionKey {
	return NodePublicEncryptionKey{}
}

func _GetCryptoMethodesStatements() NodeP2PCryptoMethode {
	return "ed25519#curve25519;"
}

func _GenerateRandom256BitValue() ([]byte, error) {
	bytes := make([]byte, 32) // 256 Bit = 32 Byte
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
