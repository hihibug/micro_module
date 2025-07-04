package etcd

type PoolData struct {
	Address string
	Weight  int
	Current int
}

type WeightedRoundRobin struct {
	servers []PoolData
}

func NewWeightedRoundRobin(servers []PoolData) *WeightedRoundRobin {
	return &WeightedRoundRobin{servers: servers}
}

func (wrr *WeightedRoundRobin) Next() string {
	var maxWeight int
	var serverIndex int

	for i, server := range wrr.servers {
		server.Current += server.Weight
		if server.Current > maxWeight {
			maxWeight = server.Current
			serverIndex = i
		}
	}

	server := wrr.servers[serverIndex]
	wrr.servers[serverIndex].Current -= wrr.servers[serverIndex].Weight
	return server.Address
}
