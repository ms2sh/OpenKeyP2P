package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ms2sh/OpenKeyP2P/src/crypto"
	"github.com/ms2sh/OpenKeyP2P/src/p2p"
)

func main() {
	if err := p2p.Setup(); err != nil {
		panic(err)
	}

	tlsConfig, err := crypto.GenerateTempTLSConfig()
	if err != nil {
		panic(err)
	}

	listenerConfig := &p2p.NodeP2PListenerConfig{AllowInternetConnection: true, AllowPrivateNetworkConnection: true, AllowAutoRouting: true, AllowTrafficForwarding: true}
	err = p2p.AddListener("0.0.0.0", 995, tlsConfig, listenerConfig)
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
