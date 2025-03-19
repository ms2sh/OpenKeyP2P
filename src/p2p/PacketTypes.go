package p2p

import openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"

// Hello Stream Packet

type HelloControlSteamPacketWSig struct {
	LocalVersion      openkeyp2p.OpenKeyP2PVesion   `cbor:"1"`
	SupportedVersions []openkeyp2p.OpenKeyP2PVesion `cbor:"2"`
	CryptoKeyMethod   NodeP2PCryptoMethode          `cbor:"3"`
	NodeConfigOptions []NodeP2PConfigEntry          `cbor:"3"`
	SignerKey         NodePublicSignatureKey        `cbor:"4"`
	EncryptionKey     NodePublicEncryptionKey       `cbor:"5"`
}

type HelloControlSteamPacket struct {
	HelloControlSteamPacketWSig
	Signature []byte `cbor:"6"`
}
