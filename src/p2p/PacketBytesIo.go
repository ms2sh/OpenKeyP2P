package p2p

import (
	"github.com/fxamacker/cbor/v2"
)

// Serialize serialisiert die Struktur in CBOR
func _SerializeHelloControlSteamPacket(packet *HelloControlSteamPacket) ([]byte, error) {
	return cbor.Marshal(packet)
}

// Deserialize deserialisiert CBOR-Daten zurück in die Struktur
func _DeserializeHelloControlSteamPacket(data []byte) (HelloControlSteamPacket, error) {
	var packet HelloControlSteamPacket
	err := cbor.Unmarshal(data, &packet)
	return packet, err
}
