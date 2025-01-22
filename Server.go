package main

import (
	"XcStorage/StorageGroup"
	"XcStorage/XcXcPanFileServer"
	"XcStorage/common/define"
	"XcStorage/etcd"
	"XcStorage/internal/config"
	"XcStorage/internal/server"
	"XcStorage/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Dir string

func main() {
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	var etcdHost = "127.0.0.1:2379"
	var addr string
	fmt.Println("Before Starting This Server, Make Sure That The Etcd Server Is Running")
	fmt.Println("Please Input The Etcd Host (127.0.0.1:2379) ：")
	fmt.Scanln(&etcdHost)
	fmt.Println("Please Input The Server Port (127.0.0.1:xxxx) ：")
	fmt.Scanln(&addr)
	fmt.Println("Please Input The Storage Directory (/home/xcstorage) ：")
	fmt.Scanln(&Dir)

	define.FILE_DIR = "./" + Dir

	os.MkdirAll(Dir, os.ModePerm)

	addr = "127.0.0.1:" + addr

	go func() {
		StorageGroup.Server = StorageGroup.NewStorageServer(addr)
		cli, err := etcd.ClientInit()
		if err != nil {
			log.Fatalf("etcd client init failed: %v", err)
		}
		otherPeers, err := etcd.GetAllPeers(cli, "xcstorage")
		if err != nil {
			log.Fatalf("get other peers from etcd failed: %v", err)
		}
		fmt.Println("other peers:", otherPeers)
		peers := append(otherPeers, addr)
		StorageGroup.Server.SetPeers(peers...)
		fmt.Println("all peers:", peers)
		go func() {
			time.Sleep(5 * time.Second)
			for {
				currentPeers, _ := etcd.GetAllPeers(cli, "xcstorage")
				closedPeers := GetClosedPeers(peers, currentPeers)
				newPeers := GetNewPeers(peers, currentPeers)
				if len(closedPeers) != 0 {
					StorageGroup.Server.DelPeers(closedPeers...)
					log.Println("Nodes have been closed :", closedPeers)
					log.Println("Current nodes :", currentPeers)
				}
				if len(newPeers) != 0 {
					StorageGroup.Server.AddPeers(newPeers...)
					log.Println("New nodes have been added :", newPeers)
					log.Println("Current nodes :", currentPeers)
				}
				peers = currentPeers
				time.Sleep(5 * time.Second)

			}
		}()

		flag.Parse()
		var c config.Config
		c.Etcd.Hosts = append([]string{}, etcdHost)
		c.Timeout = 300000
		c.Etcd.Key = "xcstorage/" + addr
		c.ListenOn = addr
		c.Name = "xcstorage_" + addr
		c.Log = logx.LogConf{
			ServiceName: "xcstorage", // 日志服务名称
			Mode:        "file",      // 日志模式，可选：console, file
			Path:        "./logs",
			Level:       "severe", // 日志级别，可选：debug, info, warn, error, severe
		}

		ctx := svc.NewServiceContext(c)

		fmt.Println("c.Etcd.Key:", c.Etcd.Key)

		s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
			XcXcPanFileServer.RegisterXcXcPanFileServiceServer(grpcServer, server.NewXcXcPanFileServiceServer(ctx))

			if c.Mode == service.DevMode || c.Mode == service.TestMode {
				reflection.Register(grpcServer)
			}
		})
		defer s.Stop()

		fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
		s.Start()
	}()

	<-exitCh
	log.Println("Shutdown signal received, exiting...")
}

func GetClosedPeers(peers, newpeers []string) []string {
	var closedPeers []string
	for _, peer := range peers {
		if !IsContain(newpeers, peer) {
			closedPeers = append(closedPeers, peer)
		}
	}
	return closedPeers
}
func GetNewPeers(peers, newpeers []string) []string {
	var newPeers []string
	for _, peer := range newpeers {
		if !IsContain(peers, peer) {
			newPeers = append(newPeers, peer)
		}
	}
	return newPeers
}

func IsContain(peers []string, peer string) bool {
	for _, p := range peers {
		if p == peer {
			return true
		}
	}
	return false
}
