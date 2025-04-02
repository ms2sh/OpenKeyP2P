package p2p

func (o *NodeP2PListenerConfig) GetConnectionConfig() NodeP2PConnectionConfig {
	return NewNodeP2PConnectionConfig()
}
