package p2p

import (
	"context"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/quic-go/quic-go"
)

func _TryOpenP2PConnectionControlStream(isIncommingConnection bool, conn quic.Connection, config *NodeP2PConnectionConfig, encpubkey NodePublicEncryptionKey, sigpub NodePublicSignatureKey, ctx context.Context) (*NodeP2PConnectionControlStream, error) {
	// Das Hello Packet wird erzeugt und in Bytes umgewandelt
	helloPacketWithoutSignature := HelloControlSteamPacketWSig{
		LocalVersion:      openkeyp2p.Version,
		SupportedVersions: openkeyp2p.SUPPORTED_VERSION,
		NodeConfigOptions: _GetStringSliceFromNodeP2PConnectionConfig(config),
		SignerKey:         sigpub,
		EncryptionKey:     encpubkey,
		CryptoKeyMethod: NodeP2PCryptoMethode{
			Name:       "ed25519#curve25519",
			Signature:  "ed25519",
			Encryption: "curve25519",
		},
	}

	// Das Paket wird Signiert und zurückgegeben
	signature, err := _SignHelloControlSteamPacketWSigPacket(&helloPacketWithoutSignature)
	if err != nil {
		return nil, err
	}

	// Die Signatur wird hinzugefügt
	helloPacket := &HelloControlSteamPacket{
		HelloControlSteamPacketWSig: helloPacketWithoutSignature,
		Signature:                   signature,
	}

	// Das Hello Packet wird in Bytes umgewandelt
	bytedHelloPacket, err := _SerializeHelloControlSteamPacket(helloPacket)
	if err != nil {
		return nil, err
	}

	// Die Streamverbindung wird aufgebaut und das Hello Packet wird übertragen
	streamConn, err := _TryOpenQuicBidirectionalStream(isIncommingConnection, conn, bytedHelloPacket, ctx)
	if err != nil {
		return nil, err
	}

	// Der Stream wird zu einem Controlstream Geupgradet
	controlStream, err := _TypeControlStreamFromBidirectionalStream(streamConn)
	if err != nil {
		return nil, err
	}

	return controlStream, nil
}

func _TypeControlStreamFromBidirectionalStream(bidstr *QuicBidirectionalStream) (*NodeP2PConnectionControlStream, error) {
	// Es wird versucht die Hello Stream Nachricht einzulesen
	helloStreamMessage, err := _DeserializeHelloControlSteamPacket(bidstr._recivedHelloBytePacket)
	if err != nil {
		return nil, err
	}

	return &NodeP2PConnectionControlStream{QuicBidirectionalStream: bidstr, destPeerHelloPacket: helloStreamMessage}, nil
}

func (o *NodeP2PConnectionControlStream) GetDestinationVersion() uint64 {
	return 0
}

func (o *NodeP2PConnectionControlStream) GetDestinationSupportedVersions() []uint64 {
	return nil
}
