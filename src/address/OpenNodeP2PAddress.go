package address

import "github.com/ms2sh/OpenKeyP2P/src/address/encoding"

func (o *OpenNodeP2PAddress) SerializeAdressByteSlice() ([]byte, error) {
	byteSlice, err := encoding.EncodeAddressByteSlice(o.typeFlag, o.encryptionPublicKey, o.signaturePublicKey, nil)
	if err != nil {
		return nil, err
	}
	return byteSlice, nil
}

func (o *OpenNodeP2PAddress) SerializeAdressString() (string, error) {
	addrstr, err := encoding.EncodeAddressToBase32StringWPrefix(o.typeFlag, o.encryptionPublicKey, o.signaturePublicKey, nil)
	if err != nil {
		return "", err
	}
	return addrstr, nil
}

func (o *OpenNodeP2PAddress) GetEncryptionKey() []byte {
	return o.encryptionPublicKey
}

func (o *OpenNodeP2PAddress) GetSignatureKey() []byte {
	return o.signaturePublicKey
}

func (o *OpenNodeP2PAddress) GetType() uint8 {
	return o.typeFlag
}

func OpenNodeP2PAddressFromPublicKeys(keyCType KeyChainType, encryptionKey []byte, signatureKey []byte) *OpenNodeP2PAddress {
	return &OpenNodeP2PAddress{encryptionPublicKey: encryptionKey, signaturePublicKey: signatureKey, typeFlag: uint8(keyCType)}
}
