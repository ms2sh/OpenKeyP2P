package p2p

import (
	"context"

	"github.com/quic-go/quic-go"
)

type ConnectionId string
type AddressType string

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
	config *NodeP2PListenerConfig
}

type NodeP2PConnectionControlStream struct {
	inControlStream  quic.Stream
	outControlStream quic.Stream
}

type NodeP2PConnection struct {
	conn                  quic.Connection
	controlStream         *NodeP2PConnectionControlStream
	config                *NodeP2PConnectionConfig
	context               context.Context
	contextCancel         context.CancelCauseFunc
	isIncommingConnection bool
}

type _ChanStreamErrorResult struct {
	stream quic.Stream
	err    error
}
