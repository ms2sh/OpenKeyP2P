package p2p

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
