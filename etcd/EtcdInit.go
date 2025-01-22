package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var defaultEtcdConfig = clientv3.Config{
	Endpoints:   []string{"localhost:2379"},
	DialTimeout: time.Second * 10,
}

func ClientInit() (*clientv3.Client, error) {
	client, err := clientv3.New(defaultEtcdConfig)
	if err != nil {
		log.Printf("[etcd] init client error: %v", err)
		return nil, err
	} else {
		log.Printf("[etcd] init client OK")
	}
	return client, err
}
