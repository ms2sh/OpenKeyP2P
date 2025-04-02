package p2p

import (
	"context"
	"sync"

	"github.com/quic-go/quic-go"
)

type ConnectionId string
type AddressType string
type NodePublicSignatureKey []byte
type NodePublicEncryptionKey []byte
type NodeP2PIpAddress []byte
type NodeP2PAdressPort uint16
type NodeP2PCryptoMethode string
type NodeP2PConnectionConfig string
type NodeP2PConnectionValidationId []byte

type NodeP2PConfigEntry struct {
	Name  string
	Value string
}

type NodeP2PConnection struct {
	conn                  quic.Connection
	controlStream         *NodeP2PControlStream
	packageTrafficStream  *NodeP2PTrafficStream
	config                NodeP2PConnectionConfig
	context               context.Context
	contextCancel         context.CancelCauseFunc
	isIncommingConnection bool
}

type NodeP2PListenerConfig struct {
	AllowInternetConnection       bool
	AllowPrivateNetworkConnection bool
	AllowAutoRouting              bool
	AllowTrafficForwarding        bool
}

type NodeP2Listener struct {
	config   *NodeP2PListenerConfig
	listener *quic.Listener
	lock     *sync.Mutex
}

type QuicBidirectionalStream struct {
	inStream                quic.Stream
	outStream               quic.Stream
	ctxCancle               context.CancelCauseFunc
	lock                    *sync.Mutex
	ctx                     context.Context
	quicConn                quic.Connection
	_sendHelloBytePacket    []byte
	_recivedHelloBytePacket []byte
}

type NodeP2PControlStream struct {
	*QuicBidirectionalStream
	destPeerHelloPacket HelloControlSteamPacket
}

type NodeP2PTrafficStream struct {
	*QuicBidirectionalStream
}
