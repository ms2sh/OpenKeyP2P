package main

import (
	"log"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
)

// RoutingTable verwaltet Peers und Whitelist.
type RoutingTable struct {
	peers     map[peer.ID]struct{}
	whitelist map[string]struct{}
	mutex     sync.RWMutex
}

// ListPeers gibt die Liste aller Peers in der Routing-Tabelle zurück.
func (rt *RoutingTable) ListPeers() []peer.ID {
	rt.mutex.RLock()
	defer rt.mutex.RUnlock()

	var peerIDs []peer.ID
	for p := range rt.peers {
		peerIDs = append(peerIDs, p)
	}
	return peerIDs
}

// NewRoutingTable erstellt eine Routing-Tabelle.
func NewRoutingTable() *RoutingTable {
	return &RoutingTable{
		peers:     make(map[peer.ID]struct{}),
		whitelist: make(map[string]struct{}),
	}
}

// WhitelistPeer fügt einen Public Key zur Whitelist hinzu.
func (rt *RoutingTable) WhitelistPeer(pubKey string) {
	rt.mutex.Lock()
	defer rt.mutex.Unlock()
	rt.whitelist[pubKey] = struct{}{}
	log.Printf("[Whitelist] Peer hinzugefügt: %s\n", pubKey)
}
