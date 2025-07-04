package rpc_test

import (
	"fmt"
	rpc2 "github.com/hihibug/microdule/v2/Framework/rpc"
	rpc "github.com/hihibug/microdule/v2/Framework/rpc/config"
	etcd2 "github.com/hihibug/microdule/v2/core/etcd"
	"github.com/hihibug/microdule/v2/core/grpc"
	"github.com/hihibug/microdule/v2/core/registry/etcd"
	"github.com/hihibug/microdule/v2/core/registry/nacos"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	//"github.com/hihibug/Framework/v2/core/grpc"

	"testing"
)

var global rpc2.Rpc

func TestEtcdGrpc(t *testing.T) {
	global = rpc2.NewRpc(grpc.NewGrpc(&rpc.Config{
		Port:        12344,
		ServiceName: "test",
		GroupName:   "def",
	}))

	etcdService := etcd.NewEtcd(etcd.Config{
		Etcd: etcd2.Config{
			Addr:     "172.16.102.119:2379",
			Password: "",
			TimeOut:  5,
		},
	})
	etcdService.SubscribeService(func(m []rpc.Config) {
		fmt.Println(m)
	})

	global.SetRegisterInterface(etcdService)

	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println(global.DiscoveryService("test"))
	}()

	defer func() {
		fmt.Println(global.Close())
	}()

	// 监听signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println(global.Close())
		os.Exit(1)
	}()
	// 获取当前进程的 id
	pid := os.Getpid()
	fmt.Println("process id:", pid)

	global.Run()
}

func TestNacosGrpc1(t *testing.T) {
	global = rpc2.NewRpc(grpc.NewGrpc(&rpc.Config{
		Port:        12346,
		ServiceName: "test1",
		GroupName:   "test1",
		Weight:      1,
	}))

	na := nacos.NewNacos(nacos.Config{
		Addr: "172.16.102.119:8848",
		ClientConfig: constant.ClientConfig{
			TimeoutMs: 5000,
			//NamespaceId: "59128e50-8ef7-4c54-9660-a88e1167cb05",
		},
	})
	na.SubscribeService(func(m []rpc.Config) {
		fmt.Println(m)
	})

	global.SetRegisterInterface(na)
	//
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println(global.DiscoveryService("test1"))
		//fmt.Println(n.DeregisterService())

	}()
	defer func() {
		fmt.Println(global.Close())
	}()

	// 监听signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println(global.Close())
		os.Exit(1)
	}()

	fmt.Println(global.Run())
}

func TestNacosGrpc(t *testing.T) {
	global = rpc2.NewRpc(grpc.NewGrpc(&rpc.Config{
		Port:        12345,
		ServiceName: "test1",
		GroupName:   "test1",
		Weight:      1,
	}))

	na := nacos.NewNacos(nacos.Config{
		Addr: "172.16.102.119:8848",
		ClientConfig: constant.ClientConfig{
			TimeoutMs: 5000,
			//NamespaceId: "59128e50-8ef7-4c54-9660-a88e1167cb05",
		},
	})
	na.SubscribeService(func(m []rpc.Config) {
		fmt.Println(m)
	})

	global.SetRegisterInterface(na)
	//
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println(global.DiscoveryService("test1"))
		//fmt.Println(n.DeregisterService())

	}()
	defer func() {
		fmt.Println(global.Close())
	}()

	// 监听signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println(global.Close())
		os.Exit(1)
	}()

	fmt.Println(global.Run())
}

func TestNacosGrpcClient(t *testing.T) {
	//conn, err := grpcA.Dial("127.0.0.1:12344", grpcA.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer conn.Close()
	//fmt.Println(_global)
	//conn, err := _global.ServiceDiscovery("127.0.0.1:12344")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(conn)
	//defer conn.(grpc.DiscoveryGrpcConnect).Close()
	//
	//fmt.Println(conn.(grpc.DiscoveryGrpcConnect))
}

func TestA(t *testing.T) {
	a := []int{1, 2, 3}
	var b = make([]int, 3)
	copy(b, a[:2])

	fmt.Println(b)
	fmt.Println(cap(b))
	b[2] = 6
	fmt.Println(unsafe.Pointer(&b))
	//b = append(b, 4)
	fmt.Println(a)

	b[2] = 7
	b = append(b, 5)
	fmt.Println(a)

	b[0] = 10
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(unsafe.Pointer(&a))
	fmt.Println(unsafe.Pointer(&b))
}
