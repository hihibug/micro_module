package grpc

import (
	"errors"
	"fmt"
	"github.com/hihibug/microdule/v2/Framework/rpc"
	"github.com/hihibug/microdule/v2/Framework/rpc/config"
	"log"
	"net"
	"strings"

	"github.com/hihibug/microdule/v2/core/utils"
	"google.golang.org/grpc/credentials/insecure"

	grpcA "google.golang.org/grpc"
)

type DiscoveryGrpcConnect = *grpcA.ClientConn

type Grpc struct {
	Grpc     *grpcA.Server
	Config   *config.Config
	Register rpc.Registry
}

func NewGrpc(c *config.Config, opt ...grpcA.ServerOption) *Grpc {
	grpcServer := grpcA.NewServer(opt...)
	return &Grpc{
		Grpc:   grpcServer,
		Config: c,
	}
}

func (g *Grpc) Client() any {
	return g
}

func (g *Grpc) SetRegisterInterface(registry rpc.Registry) {
	g.Register = registry
}

func (g *Grpc) DiscoveryService(sName string) (any, error) {
	if g.Register == nil && net.ParseIP(strings.Split(sName, ":")[0]) == nil {
		return nil, errors.New("service not register")
	}
	if g.Register != nil {
		addr, err := g.Register.DiscoveryService(sName, g.Config.GroupName, g.Config.ClusterName)
		if err != nil {
			return nil, err
		}
		sName = addr
	}
	return grpcA.Dial(
		sName,
		grpcA.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func (g *Grpc) Run() error {
	ip, _ := utils.ExternalIP()
	g.Config.Ip = ip
	log.Printf("rpc   name: %s", g.Config.ServiceName)
	log.Printf("rpc   mod: %s", "grpc")
	if g.Register != nil {
		err := g.Register.RegisterService(g.Config)
		if err != nil {
			return err
		}
	}

	address := fmt.Sprintf(":%d", g.Config.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic("grpc net listen error :" + err.Error())
	}

	err = g.Grpc.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func (g *Grpc) Close() error {
	return g.Register.DeregisterService()
}
