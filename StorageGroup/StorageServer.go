package StorageGroup

import (
	"XcStorage/consistentHash"
	"sync"
)

const serverPrefix = "xcstorage/"

const defaultReplicas = 50

type StorageServer struct {
	Addr        string
	mu          sync.Mutex
	Peers       *consistentHash.Map
	GrpcGetters map[string]*StorageClient
}

var Server *StorageServer

func NewStorageServer(addr string) *StorageServer {
	return &StorageServer{
		Addr: addr,
	}
}

// 设置节点
func (s *StorageServer) SetPeers(peers ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Peers = consistentHash.New(defaultReplicas, nil)
	s.Peers.Add(peers...)
	s.GrpcGetters = make(map[string]*StorageClient, len(peers))
	for _, peer := range peers {
		s.GrpcGetters[peer] = &StorageClient{serviceName: serverPrefix + peer}
	}
}

// 删除节点
func (s *StorageServer) DelPeers(peers ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, peer := range peers {
		if _, ok := s.GrpcGetters[peer]; ok {
			delete(s.GrpcGetters, peer)
			s.Peers.Remove(peer)
		}
	}
}

// 添加节点
func (s *StorageServer) AddPeers(peers ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, peer := range peers {
		if _, ok := s.GrpcGetters[peer]; !ok {
			s.GrpcGetters[peer] = &StorageClient{serviceName: serverPrefix + peer}
			s.Peers.Add(peer)
		}
	}
}

func (s *StorageServer) PickPeer(key string) (PeerGetter, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if peer := s.Peers.Get(key); peer != "" && peer != s.Addr {
		return s.GrpcGetters[peer], true
	}
	return nil, false
}
