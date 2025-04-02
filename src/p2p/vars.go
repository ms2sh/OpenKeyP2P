package p2p

import "sync"

var (
	nodeConnections map[ConnectionId]*NodeP2PConnection
	controlLock     *sync.Mutex = new(sync.Mutex)
	wasSetuped      bool        = false
)

func _VarsAddNodeConnection(nodeConn *NodeP2PConnection) error {
	controlLock.Lock()
	defer controlLock.Unlock()
	nodeConnections[nodeConn.GetConnectionId()] = nodeConn
	return nil
}

func _VarsDeleteNodeConnection(nodeConn *NodeP2PConnection) {
	controlLock.Lock()
	defer controlLock.Unlock()
	nodeConnections[nodeConn.GetConnectionId()] = nodeConn
}

func _VarsWasSetuped() bool {
	controlLock.Lock()
	reval := wasSetuped
	controlLock.Unlock()
	return reval
}
