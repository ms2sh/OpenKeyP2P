package p2p

import (
	"context"
	"sync"
	"time"

	openkeyp2p "github.com/ms2sh/OpenKeyP2P/src"
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
type NodeP2PKeepaliveProcessId []byte
type NodeP2PSocketAddress string

type NodeP2PConfigEntry struct {
	Name  string
	Value string
}

type NodeP2PConnection struct {
	conn                    quic.Connection
	controlStream           *NodeP2PControlStream
	packageTrafficStream    *NodeP2PTrafficStream
	config                  NodeP2PConnectionConfig
	ctx                     context.Context
	contextCancel           context.CancelCauseFunc
	writerControlBuffer     *openkeyp2p.ThreadSafeContextBuffer
	localKeepalivePacketIds *sync.Map
	isIncommingConnection   bool
	keepaliveTime           time.Duration
	localSocketAddress      NodeP2PSocketAddress
	remoteSocketAddress     NodeP2PSocketAddress
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
	writeMutex              *sync.Mutex
	readMutex               *sync.Mutex
	_sendHelloBytePacket    []byte
	_recivedHelloBytePacket []byte
	_localSocketEp          NodeP2PSocketAddress
	_remoteSocketEp         NodeP2PSocketAddress
}

type NodeP2PControlStream struct {
	*QuicBidirectionalStream
	destPeerHelloPacket L1HelloControlSteamPacket
}

type NodeP2PTrafficStream struct {
	*QuicBidirectionalStream
}

type _NodeP2pKeepaliveProcess struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	LMutex *sync.Mutex
	Finish bool
}
