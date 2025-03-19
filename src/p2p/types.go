package p2p

import (
	"context"
	"sync"

	"github.com/quic-go/quic-go"
)

type ConnectionId string
type AddressType string
type NodeP2PConfigEntry string
type NodePublicSignatureKey []byte
type NodePublicEncryptionKey []byte

type NodeP2PCryptoMethode struct {
	Name       string
	Signature  string
	Encryption string
}

type NodeP2PConnection struct {
	conn                  quic.Connection
	controlStream         *NodeP2PConnectionControlStream
	packageTrafficStream  *_NodeP2PConnectionPackageTrafficStream
	config                *NodeP2PConnectionConfig
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

type NodeP2PConnectionConfig struct {
	AllowAutoRouting       bool
	AllowTrafficForwarding bool
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
	_sendHelloBytePacket    []byte
	_recivedHelloBytePacket []byte
}

type NodeP2PConnectionControlStream struct {
	*QuicBidirectionalStream
	destPeerHelloPacket HelloControlSteamPacket
}

type _NodeP2PConnectionPackageTrafficStream struct {
	*QuicBidirectionalStream
}
