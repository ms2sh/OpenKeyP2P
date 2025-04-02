package p2p

import openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"

type HelloControlSteamPacketWSig struct {
	LocalVersion       openkeyp2p.OpenKeyP2PVesion   `cbor:"2"`
	SupportedVersions  []openkeyp2p.OpenKeyP2PVesion `cbor:"3"`
	CryptoKeyMethod    NodeP2PCryptoMethode          `cbor:"4"`
	NodeConfigOptions  NodeP2PConnectionConfig       `cbor:"5"`
	YourIpAddress      NodeP2PIpAddress              `cbor:"6"`
	YourIpPort         NodeP2PAdressPort             `cbor:"7"`
	SignerKey          NodePublicSignatureKey        `cbor:"8"`
	EncryptionKey      NodePublicEncryptionKey       `cbor:"9"`
	CMTU               uint16                        `cbor:"10"`
	ACKPerPackage      bool                          `cbor:"11"`
	MaxPacketPerSecond uint16                        `cbor:"12"`
}

type HelloControlSteamPacket struct {
	HelloControlSteamPacketWSig
	Signature []byte `cbor:"1"`
}

type HelloTrafficStreamPacket struct {
	ValId     NodeP2PConnectionValidationId `cbor:"2"`
	Signature []byte                        `cbor:"1"`
}
