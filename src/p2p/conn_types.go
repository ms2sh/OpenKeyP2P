package p2p

import (
	"net"
	"sync"
)

// Conn wraps a net.Conn with a Mutex and a Condition to synchronize read/write operations.
type FWConn struct {
	conn         net.Conn   // Underlying network connection
	mu           sync.Mutex // Mutex to protect concurrent access
	cond         *sync.Cond // Condition variable for synchronization
	connType     _ConnType  // Connection Protocol
	maxFrameSize uint64
	//reading  bool       // Flag to check if a read operation is active
	//writing  bool       // Flag to check if a write operation is active
}
