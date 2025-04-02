package p2p

import "fmt"

func Setup() error {
	controlLock.Lock()
	defer controlLock.Unlock()

	if wasSetuped {
		return fmt.Errorf("was always setup")
	}

	nodeConnections = make(map[ConnectionId]*NodeP2PConnection)

	wasSetuped = true

	return nil
}
