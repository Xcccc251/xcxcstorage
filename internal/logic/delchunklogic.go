package logic

import (
	"XcStorage/StorageGroup"
	"XcStorage/XcXcPanFileServer"
	"XcStorage/common/define"
	"log"
	"os"

	"XcStorage/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelChunkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelChunkLogic {
	return &DelChunkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelChunkLogic) DelChunk(in *XcXcPanFileServer.DelChunkRequest) (*XcXcPanFileServer.DelChunkResponse, error) {
	chunkId := in.GetChunkId()
	log.Printf("[xcstorage %s] Recv RPC Request->DelChunk - ChunkId: %s", StorageGroup.Server.Addr, chunkId)
	peer := StorageGroup.Server.Peers.Get(chunkId)
	if peer != StorageGroup.Server.Addr {
		log.Printf("[xcstorage %s] [Chunk: %s] is be stored in peer %s", StorageGroup.Server.Addr, chunkId, peer)
		success, err := StorageGroup.Server.GrpcGetters[peer].Del(chunkId)
		if err != nil || !success {
			return &XcXcPanFileServer.DelChunkResponse{
				Success: false,
			}, err
		}
		return &XcXcPanFileServer.DelChunkResponse{
			Success: true,
		}, nil
	}

	path := define.FILE_DIR + "/" + chunkId
	err := os.Remove(path)
	if err != nil {
		return &XcXcPanFileServer.DelChunkResponse{
			Success: false,
		}, err
	}
	log.Printf("[xcstorage %s] [Chunk: %s] has be deleted", StorageGroup.Server.Addr, chunkId)
	return &XcXcPanFileServer.DelChunkResponse{
		Message: "删除成功",
		Success: true,
	}, nil
}
