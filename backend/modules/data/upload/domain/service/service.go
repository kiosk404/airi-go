package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/data/upload/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/data/upload/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
)

type uploadSVC struct {
	fileRepo repo.FilesRepo
	idgen    idgen.IDGenerator
	oss      storage.Storage
}

func NewUploadSVC(rdb rdb.Provider, idgen idgen.IDGenerator, oss storage.Storage) UploadService {
	db := rdb.NewSession(context.Background()).DB()
	return &uploadSVC{fileRepo: repo.NewFilesRepo(db), idgen: idgen, oss: oss}
}

func (u *uploadSVC) UploadFile(ctx context.Context, req *UploadFileRequest) (resp *UploadFileResponse, err error) {
	resp = &UploadFileResponse{}
	if req.File.ID == 0 {
		req.File.ID, err = u.idgen.GenID(ctx)
		if err != nil {
			return nil, errorx.New(errno.ErrIDGenError)
		}
	}
	err = u.fileRepo.Create(ctx, req.File)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrUploadSystemErrorCode)
	}
	resp.File = req.File
	return
}

func (u *uploadSVC) UploadFiles(ctx context.Context, req *UploadFilesRequest) (resp *UploadFilesResponse, err error) {
	resp = &UploadFilesResponse{}
	for _, file := range req.Files {
		if file.ID == 0 {
			file.ID, err = u.idgen.GenID(ctx)
			if err != nil {
				return nil, errorx.New(errno.ErrIDGenError)
			}
		}
	}
	err = u.fileRepo.BatchCreate(ctx, req.Files)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrUploadSystemErrorCode)
	}
	resp.Files = req.Files
	return
}

func (u *uploadSVC) GetFiles(ctx context.Context, req *GetFilesRequest) (resp *GetFilesResponse, err error) {
	resp = &GetFilesResponse{}
	resp.Files, err = u.fileRepo.MGetByIDs(ctx, req.IDs)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrUploadSystemErrorCode)
	}
	return
}

func (u *uploadSVC) GetFile(ctx context.Context, req *GetFileRequest) (resp *GetFileResponse, err error) {
	resp = &GetFileResponse{}
	resp.File, err = u.fileRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrUploadSystemErrorCode)
	}
	if resp.File != nil {
		url, err := u.oss.GetObjectUrl(ctx, resp.File.TosURI)
		if err == nil {
			resp.File.Url = url
		}
	}
	return
}
