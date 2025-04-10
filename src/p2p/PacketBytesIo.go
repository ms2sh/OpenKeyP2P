package p2p

import (
	"github.com/fxamacker/cbor/v2"
)

// Serialize serialisiert die Struktur in CBOR
func _SerializeSteamPacket(packet interface{}) ([]byte, error) {
	return cbor.Marshal(packet)
}

// Deserialize deserialisiert CBOR-Daten zurück in die Struktur
func _DeserializeHelloControlSteamPacket(data []byte) (L1HelloControlSteamPacket, error) {
	var packet L1HelloControlSteamPacket
	err := cbor.Unmarshal(data, &packet)
	return packet, err
}

// Deserialize deserialisiert CBOR-Daten zurück in die Struktur
func _DeserializeTrafficSteamPacket(data []byte) (L1HelloTrafficStreamPacket, error) {
	var packet L1HelloTrafficStreamPacket
	err := cbor.Unmarshal(data, &packet)
	return packet, err
}
