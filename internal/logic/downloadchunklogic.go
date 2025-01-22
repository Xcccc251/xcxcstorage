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

type DownloadChunkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDownloadChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadChunkLogic {
	return &DownloadChunkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DownloadChunkLogic) DownloadChunk(in *XcXcPanFileServer.DownloadChunkRequest) (*XcXcPanFileServer.DownloadChunkResponse, error) {
	chunkId := in.GetChunkId()
	log.Printf("[xcstorage %s] Recv RPC Request->DownloadChunk - ChunkId: %s", StorageGroup.Server.Addr, chunkId)
	peer := StorageGroup.Server.Peers.Get(chunkId)
	if peer != StorageGroup.Server.Addr {
		log.Printf("[xcstorage %s] [Chunk: %s] is be stored in peer %s", StorageGroup.Server.Addr, chunkId, peer)
		data, err := StorageGroup.Server.GrpcGetters[peer].Download(chunkId)
		if err != nil {
			return &XcXcPanFileServer.DownloadChunkResponse{
				Data:    nil,
				Success: false,
			}, err
		}
		return &XcXcPanFileServer.DownloadChunkResponse{
			Data:    data,
			Success: true,
		}, nil
	}

	path := define.FILE_DIR + "/" + chunkId
	if data, err := os.ReadFile(path); err != nil {
		return &XcXcPanFileServer.DownloadChunkResponse{
			Data:    nil,
			Success: false,
		}, err
	} else {
		return &XcXcPanFileServer.DownloadChunkResponse{
			Data:    data,
			Success: true,
		}, nil
	}
}
