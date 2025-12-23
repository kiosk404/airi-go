package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/data/upload/domain/entity"
)

type UploadService interface {
	UploadFile(ctx context.Context, req *UploadFileRequest) (resp *UploadFileResponse, err error)
	UploadFiles(ctx context.Context, req *UploadFilesRequest) (resp *UploadFilesResponse, err error)
	GetFiles(ctx context.Context, req *GetFilesRequest) (resp *GetFilesResponse, err error)
	GetFile(ctx context.Context, req *GetFileRequest) (resp *GetFileResponse, err error)
}

type UploadFileRequest struct {
	File *entity.File `json:"file"`
}
type UploadFileResponse struct {
	File *entity.File `json:"file"`
}
type UploadFilesRequest struct {
	Files []*entity.File `json:"files"`
}

type UploadFilesResponse struct {
	Files []*entity.File `json:"files"`
}

type GetFilesRequest struct {
	IDs []int64 `json:"ids"`
}

type GetFilesResponse struct {
	Files []*entity.File `json:"files"`
}

type GetFileRequest struct {
	ID int64 `json:"id"`
}

type GetFileResponse struct {
	File *entity.File `json:"file"`
}
