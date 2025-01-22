package main

import (
	"XcStorage/XcXcPanFileServer"
	"XcStorage/etcd"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	cli, err := etcd.ClientInit()
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()
	conn, err := etcd.DiscoverFromEtcd(cli, "xcstorage/127.0.0.1:9001")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	client := XcXcPanFileServer.NewXcXcPanFileServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//path := "./test2.mp4"
	//file, err := os.Open(path)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//data, err := Server_Helper.FileToBytes(file)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer file.Close()
	rsp, err := client.UploadChunk(ctx, &XcXcPanFileServer.UploadChunkRequest{
		ChunkId: "12345",
		Data:    []byte("123"),
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp.Success)
}
func TestDel(t *testing.T) {
	cli, err := etcd.ClientInit()
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()
	conn, err := etcd.DiscoverFromEtcd(cli, "xcstorage/127.0.0.1:9001")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	client := XcXcPanFileServer.NewXcXcPanFileServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	rsp, err := client.DelChunk(ctx, &XcXcPanFileServer.DelChunkRequest{
		ChunkId: "12345",
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp.Success)
}

func TestDownload(t *testing.T) {
	cli, err := etcd.ClientInit()
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()
	conn, err := etcd.DiscoverFromEtcd(cli, "xcstorage/127.0.0.1:9001")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	client := XcXcPanFileServer.NewXcXcPanFileServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	rsp, err := client.DownloadChunk(ctx, &XcXcPanFileServer.DownloadChunkRequest{
		ChunkId: "12345",
	})
	//create, err := os.Create("downloadtest.mp4")
	//if err != nil {
	//	fmt.Println(err)
	//}
	fmt.Println(rsp.Data)

	//defer create.Close()
	//create.Write(rsp.Data)
	fmt.Println(rsp.Success)
}
