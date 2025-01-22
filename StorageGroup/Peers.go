package StorageGroup

// PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGetter。
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// Get() 方法用于从对应 group 查找缓存值。
type PeerGetter interface {
	Del(chunkId string) (bool, error)
	Upload(chunkId string, data []byte) (bool, error)
	Download(chunkId string) ([]byte, error)
}
