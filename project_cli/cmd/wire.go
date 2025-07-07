//go:build wireinject
// +build wireinject

package main

import (
	"project_cli/pkg"

	"github.com/google/wire"
)

func InitService() *pkg.Pkg {
	wire.Build(
		pkg.ProviderPkg,
	)

	return new(pkg.Pkg)
}
