package p2p

import (
	"context"

	"github.com/quic-go/quic-go"
)

func _TryOpenP2PConnectionTrafficStream(isIncommingConnection bool, conn quic.Connection, ctx context.Context) (*NodeP2PTrafficStream, error) {
	// Es wird ein Zufälliger Wert erzeugt
	randomValue, err := _BuildRandomVIdValue()
	if err != nil {
		return nil, err
	}

	// Das Paket wird Signiert und zurückgegeben
	signature, err := _SignByteSlice(randomValue)
	if err != nil {
		return nil, err
	}

	// Die Signatur wird hinzugefügt
	helloPacket := &HelloTrafficStreamPacket{
		ValId:     randomValue,
		Signature: signature,
	}

	// Das Hello Packet wird in Bytes umgewandelt
	bytedHelloPacket, err := _SerializeSteamPacket(helloPacket)
	if err != nil {
		return nil, err
	}

	// Die Streamverbindung wird aufgebaut und das Hello Packet wird übertragen
	streamConn, err := _TryOpenQuicBidirectionalStream(isIncommingConnection, conn, bytedHelloPacket, ctx)
	if err != nil {
		return nil, err
	}

	// Der Stream wird zu einem TrafficStream Geupgradet
	TrafficStream, err := _TypeTrafficStreamFromBidirectionalStream(streamConn)
	if err != nil {
		return nil, err
	}

	return TrafficStream, nil
}

func _TypeTrafficStreamFromBidirectionalStream(bidstr *QuicBidirectionalStream) (*NodeP2PTrafficStream, error) {
	// Es wird versucht die Hello Stream Nachricht einzulesen
	_, err := _DeserializeTrafficSteamPacket(bidstr._recivedHelloBytePacket)
	if err != nil {
		return nil, err
	}

	return &NodeP2PTrafficStream{QuicBidirectionalStream: bidstr}, nil
}
