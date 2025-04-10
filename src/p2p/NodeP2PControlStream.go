package p2p

import (
	"context"
	"net"
	"strconv"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
	"github.com/quic-go/quic-go"
)

func _TryOpenP2PConnectionControlStream(localhostNetworkInterface *net.Interface, isIncommingConnection bool, conn quic.Connection, config NodeP2PConnectionConfig, localSocketEp NodeP2PSocketAddress, remoteSocketEp NodeP2PSocketAddress, connCtx context.Context, connCtxCancel context.CancelCauseFunc) (*NodeP2PControlStream, error) {
	// IP und Port extrahieren
	host, portStr, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return nil, err
	}

	// Port in uint16 umwandeln
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, err
	}

	// IP-Adresse in Bytes umwandeln
	ip := net.ParseIP(host)
	if ip == nil {
		return nil, err
	}

	// IPv4 oder IPv6 als Bytes ausgeben
	var ipBytes []byte
	var mtu int
	if ip4 := ip.To4(); ip4 != nil {
		ipBytes = ip4 // IPv4 als 4-Byte Array
		mtu = openkeyp2p.CalculateQUICPayloadSize(localhostNetworkInterface.MTU, false)
	} else {
		ipBytes = ip.To16() // IPv6 als 16-Byte Array
		mtu = openkeyp2p.CalculateQUICPayloadSize(localhostNetworkInterface.MTU, true)
	}

	// Das Hello Packet wird erzeugt und in Bytes umgewandelt
	helloPacketWithoutSignature := L1HelloControlSteamPacketWSig{
		LocalVersion:       openkeyp2p.Version,
		SupportedVersions:  openkeyp2p.SUPPORTED_VERSION,
		NodeConfigOptions:  config,
		SignerKey:          _GetSignerPublicKey(),
		EncryptionKey:      _GetEncryptionPublicKey(),
		YourIpPort:         NodeP2PAdressPort(port),
		YourIpAddress:      NodeP2PIpAddress(ipBytes),
		CryptoKeyMethod:    _GetCryptoMethodesStatements(),
		CMTU:               uint16(mtu),
		ACKPerPackage:      false,
		MaxPacketPerSecond: 0,
	}

	// Das Paket wird Signiert und zurückgegeben
	signature, err := _SignSteamPacketWSigPacket(&helloPacketWithoutSignature)
	if err != nil {
		return nil, err
	}

	// Die Signatur wird hinzugefügt
	helloPacket := &L1HelloControlSteamPacket{
		L1HelloControlSteamPacketWSig: helloPacketWithoutSignature,
		Signature:                     signature,
	}

	// Das Hello Packet wird in Bytes umgewandelt
	bytedHelloPacket, err := _SerializeSteamPacket(helloPacket)
	if err != nil {
		return nil, err
	}

	// Die Streamverbindung wird aufgebaut und das Hello Packet wird übertragen
	streamConn, err := _TryOpenQuicBidirectionalStream(isIncommingConnection, conn, bytedHelloPacket, localSocketEp, remoteSocketEp, connCtx, connCtxCancel)
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

func _TypeControlStreamFromBidirectionalStream(bidstr *QuicBidirectionalStream) (*NodeP2PControlStream, error) {
	// Es wird versucht die Hello Stream Nachricht einzulesen
	helloStreamMessage, err := _DeserializeHelloControlSteamPacket(bidstr._recivedHelloBytePacket)
	if err != nil {
		return nil, err
	}

	return &NodeP2PControlStream{QuicBidirectionalStream: bidstr, destPeerHelloPacket: helloStreamMessage}, nil
}

func (o *NodeP2PControlStream) GetDestinationVersion() openkeyp2p.OpenKeyP2PVesion {
	return o.destPeerHelloPacket.LocalVersion
}

func (o *NodeP2PControlStream) GetDestinationSupportedVersions() []openkeyp2p.OpenKeyP2PVesion {
	return o.destPeerHelloPacket.SupportedVersions
}

func (o *NodeP2PControlStream) GetMyLocalIPByAnotherPeer() NodeP2PIpAddress {
	return o.destPeerHelloPacket.YourIpAddress
}

func (o *NodeP2PControlStream) GetMyLocalPortByAnotherPeer() NodeP2PAdressPort {
	return o.destPeerHelloPacket.YourIpPort
}

func (o *NodeP2PControlStream) GetMTU() uint16 {
	return o.destPeerHelloPacket.CMTU
}

func (o *NodeP2PControlStream) GetACKPeerPacket() bool {
	return o.destPeerHelloPacket.ACKPerPackage
}
