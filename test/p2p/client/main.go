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

	connectionConfig := p2p.NewNodeP2PConnectionConfig()
	connectionConfig.Add("auto-routing", "yes")

	err = p2p.ConnectTo("quic://152.53.118.14:995", tlsConfig, connectionConfig)
	if err != nil {
		panic(err.Error())
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
