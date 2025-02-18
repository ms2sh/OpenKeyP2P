package address

import "github.com/ms2sh/OpenKeyP2P/src/address/encoding"

func (o *OpenNodeP2PSocketAddress) SerializeAdressByteSlice() ([]byte, error) {
	byteSlice, err := encoding.EncodeAddressByteSlice(o.typeFlag, o.encryptionPublicKey, o.signaturePublicKey, o.port)
	if err != nil {
		return nil, err
	}
	return byteSlice, nil
}

func (o *OpenNodeP2PSocketAddress) SerializeAdressString() (string, error) {
	addrstr, err := encoding.EncodeAddressToBase32StringWPrefix(o.typeFlag, o.encryptionPublicKey, o.signaturePublicKey, o.port)
	if err != nil {
		return "", err
	}
	return addrstr, nil
}

func (o *OpenNodeP2PSocketAddress) GetPort() []byte {
	return o.port
}
