package main

import "fmt"

func main() {
	pkg := InitService()

	fmt.Println(
		pkg.Srv.Config().Data.Name,
		pkg.Srv.Modules("http"),
	)

	fmt.Println(pkg.Srv.Start())
}
