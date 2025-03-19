package crypto

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"strings"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
)

type OpenKeyP2PAddress struct {
	Prefix openkeyp2p.OpenKeyP2PPrefix
	PubKey openkeyp2p.OpenKeyP2PPublicKey
}

func (o *OpenKeyP2PAddress) ToString() string {
	checksum := o.ComputeChecksum()
	bytesSlice := make([]byte, 0)
	bytesSlice = append(bytesSlice, 0)
	bytesSlice = append(bytesSlice, o.PubKey...)
	bytesSlice = append(bytesSlice, checksum...)
	encodedPublicKey := openkeyp2p.Base32Encoding.EncodeToString(bytesSlice)
	return fmt.Sprintf("%s%s", o.Prefix, encodedPublicKey)
}

func (o *OpenKeyP2PAddress) ToByteSlice() []byte {
	checksum := o.ComputeChecksum()
	bytesSlice := make([]byte, 0)
	bytesSlice = append(bytesSlice, []byte(openkeyp2p.Prefix)...)
	bytesSlice = append(bytesSlice, 0)
	bytesSlice = append(bytesSlice, o.PubKey...)
	bytesSlice = append(bytesSlice, checksum...)
	return bytesSlice
}

func (o *OpenKeyP2PAddress) ComputeChecksum() []byte {
	bytesSlice := make([]byte, 0)
	bytesSlice = append(bytesSlice, []byte(openkeyp2p.Prefix)...)
	bytesSlice = append(bytesSlice, 0)
	bytesSlice = append(bytesSlice, o.PubKey...)
	return ComputeChecksumCRC32(bytesSlice)
}

func (o *OpenKeyP2PAddress) ComputeHash() (openkeyp2p.HashSlice, error) {
	return ComputeHash(openkeyp2p.DEFAULT_HASH_METHODE_256BIT, o.ToByteSlice())
}

func (o *OpenKeyP2PAddress) VerifySignature(signature OpenKeyP2PSignature, dataHash openkeyp2p.HashSlice) (bool, error) {
	addrHash, err := o.ComputeHash()
	if err != nil {
		return false, err
	}

	dataAddressHashCombination, err := ComputeHash(openkeyp2p.DEFAULT_HASH_METHODE_256BIT, addrHash, dataHash)
	if err != nil {
		return false, err
	}

	verifySigResult, err := signature.VerifySignatureHashDigest(dataAddressHashCombination, o.PubKey)
	if err != nil {
		return false, err
	}

	return verifySigResult, nil
}

func OpenKeyP2PAddressDecodeFromByteSlice(adrString []byte) (*OpenKeyP2PAddress, error) {
	if !bytes.HasPrefix(adrString, []byte(openkeyp2p.Prefix)) {
		return nil, fmt.Errorf("string has no valid prefix")
	}

	plainAddrByteSlice := bytes.ReplaceAll(adrString, []byte(openkeyp2p.Prefix), []byte{})

	newAddr := &OpenKeyP2PAddress{Prefix: openkeyp2p.Prefix, PubKey: plainAddrByteSlice[1:33]}
	if !bytes.Equal(newAddr.ComputeChecksum(), plainAddrByteSlice[33:]) {
		return nil, fmt.Errorf("checksum invalid")
	}

	return newAddr, nil
}

func OpenKeyP2PAddressDecodeFromString(adrString string) (*OpenKeyP2PAddress, error) {
	if !strings.HasPrefix(adrString, string(openkeyp2p.Prefix)) {
		return nil, fmt.Errorf("string has no valid prefix")
	}

	base32Str := strings.ReplaceAll(adrString, string(openkeyp2p.Prefix), "")
	decodedAddress, err := openkeyp2p.Base32Encoding.DecodeString(base32Str)
	if err != nil {
		return nil, err
	}

	newAddr := &OpenKeyP2PAddress{Prefix: openkeyp2p.Prefix, PubKey: decodedAddress[1:33]}
	if !bytes.Equal(newAddr.ComputeChecksum(), decodedAddress[33:]) {
		return nil, fmt.Errorf("checksum invalid")
	}

	return newAddr, nil
}

func OpenKeyP2PAddressFromPublicKey(pubKey ed25519.PublicKey) (*OpenKeyP2PAddress, error) {
	adrStruct := &OpenKeyP2PAddress{
		Prefix: openkeyp2p.Prefix,
		PubKey: openkeyp2p.OpenKeyP2PPublicKey(pubKey),
	}
	return adrStruct, nil
}

func AddressSign(privKey ed25519.PrivateKey, dataHash openkeyp2p.HashSlice) (OpenKeyP2PSignature, error) {
	pubkey := PrivateKeyToPublicKey(privKey)

	addr, err := OpenKeyP2PAddressFromPublicKey(pubkey)
	if err != nil {
		return nil, err
	}

	addrHash, err := addr.ComputeHash()
	if err != nil {
		return nil, err
	}

	dataAddressHashCombination, err := ComputeHash(openkeyp2p.DEFAULT_HASH_METHODE_256BIT, addrHash, dataHash)
	if err != nil {
		return nil, err
	}

	signSlcie := ed25519.Sign(privKey, dataAddressHashCombination)

	return OpenKeyP2PSignature(signSlcie), nil
}
