package address

type OpenNodeP2PAddress struct {
	encryptionPublicKey []byte
	signaturePublicKey  []byte
	typeFlag            uint8
}

type OpenNodeP2PSocketAddress struct {
	*OpenNodeP2PAddress
	port []byte
}

type KeyChainType uint8
