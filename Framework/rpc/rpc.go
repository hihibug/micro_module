package rpc

type Rpc interface {
	Client() any
	SetRegisterInterface(registry Registry)
	DiscoveryService(sName string) (any, error)
	Run() error
	Close() error
}

func NewRpc(result Rpc) Rpc {
	return result
}
