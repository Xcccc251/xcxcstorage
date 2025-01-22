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

type UploadChunkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadChunkLogic {
	return &UploadChunkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadChunkLogic) UploadChunk(in *XcXcPanFileServer.UploadChunkRequest) (*XcXcPanFileServer.UploadChunkResponse, error) {
	chunkId := in.GetChunkId()
	data := in.GetData()
	log.Printf("[xcstorage %s] Recv RPC Request->UploadChunk - ChunkId: %s", StorageGroup.Server.Addr, chunkId)
	peer := StorageGroup.Server.Peers.Get(chunkId)
	if peer != StorageGroup.Server.Addr {
		log.Printf("[xcstorage %s] [Chunk: %s] should be stored in peer %s", StorageGroup.Server.Addr, chunkId, peer)
		success, err := StorageGroup.Server.GrpcGetters[peer].Upload(chunkId, data)
		if err != nil || !success {
			return &XcXcPanFileServer.UploadChunkResponse{
				Message: "Fail",
				Success: false,
			}, err
		}
		return &XcXcPanFileServer.UploadChunkResponse{
			Message: "OK",
			Success: true,
		}, nil
	}

	chunk, err := os.Create(define.FILE_DIR + "/" + chunkId)
	if err != nil {
		return &XcXcPanFileServer.UploadChunkResponse{
			Success: false,
		}, err
	}
	defer chunk.Close()
	_, err = chunk.Write(data)
	if err != nil {
		return &XcXcPanFileServer.UploadChunkResponse{
			Success: false,
			Message: "Fail",
		}, err
	}
	chunk.Seek(0, 0)

	return &XcXcPanFileServer.UploadChunkResponse{
		Success: true,
		Message: "上传成功",
	}, nil
}
