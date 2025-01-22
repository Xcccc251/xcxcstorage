package StorageGroup

import (
	"XcStorage/XcXcPanFileServer"
	"XcStorage/etcd"
	"context"
	"fmt"
	"time"
)

type StorageClient struct {
	serviceName string
}

func (sc *StorageClient) Del(chunkId string) (bool, error) {
	//创建一个etcd client
	cli, err := etcd.ClientInit()
	if err != nil {
		return false, err
	}
	defer cli.Close()
	conn, err := etcd.DiscoverFromEtcd(cli, sc.serviceName)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	client := XcXcPanFileServer.NewXcXcPanFileServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	rsp, err := client.DelChunk(ctx, &XcXcPanFileServer.DelChunkRequest{
		ChunkId: chunkId,
	})
	if err != nil {
		return false, err
	}
	return rsp.Success, nil
}

func (sc *StorageClient) Upload(chunkId string, data []byte) (bool, error) {
	cli, err := etcd.ClientInit()
	if err != nil {
		return false, err
	}
	defer cli.Close()
	conn, err := etcd.DiscoverFromEtcd(cli, sc.serviceName)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	client := XcXcPanFileServer.NewXcXcPanFileServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	rsp, err := client.UploadChunk(ctx, &XcXcPanFileServer.UploadChunkRequest{
		ChunkId: chunkId,
		Data:    data,
	})
	if err != nil {
		return false, err
	}
	return rsp.Success, nil
}

func (sc *StorageClient) Download(chunkId string) ([]byte, error) {
	cli, err := etcd.ClientInit()
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	conn, err := etcd.DiscoverFromEtcd(cli, sc.serviceName)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	client := XcXcPanFileServer.NewXcXcPanFileServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	rsp, err := client.DownloadChunk(ctx, &XcXcPanFileServer.DownloadChunkRequest{
		ChunkId: chunkId,
	})
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}
