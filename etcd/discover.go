package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"time"
)

// 服务发现,获取指定地址
func DiscoverFromEtcd(c *clientv3.Client, service string) (*grpc.ClientConn, error) {
	// 获取 etcd 中的 key
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := c.Get(ctx, service, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var addr string

	// 检查 key 是否存在
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("no address found for key: %s", service)
	}
	addr = string(resp.Kvs[0].Value)
	log.Printf("[etcd] get address ok: %s \n", addr)

	// 创建 gRPC 客户端连接
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server at %s: %w", addr, err)
	}

	return conn, nil
}

func GetAllPeers(c *clientv3.Client, prefix string) ([]string, error) {
	resp, err := c.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var addresses []string

	for _, kv := range resp.Kvs {
		addr := string(kv.Value)
		addresses = append(addresses, addr)
	}
	return addresses, nil

}
