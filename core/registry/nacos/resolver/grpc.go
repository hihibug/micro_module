package resolver

import (
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"google.golang.org/grpc/resolver"
)

type GrpcResolverBuilder struct {
	Client *[]model.Instance
}

func NewGrpcResolverBuilder(c *[]model.Instance) *GrpcResolverBuilder {
	return &GrpcResolverBuilder{
		Client: c,
	}
}

func (r *GrpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	
}

func (r *GrpcResolverBuilder) Scheme() string {
	return "nacos"
}
