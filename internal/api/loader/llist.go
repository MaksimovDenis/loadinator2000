package loader

import (
	"context"

	"github.com/MaksimovDenis/loadinator2000/internal/models"
	desc "github.com/MaksimovDenis/loadinator2000/pkg/loader_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (hdl *Implementation) List(ctx context.Context, req *desc.ListRequest) (*desc.ListResponse, error) {
	files, err := hdl.loaderSecrvice.List(ctx, req.Limit, req.Offset)
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to get files")
		return nil, err
	}

	return &desc.ListResponse{
		Files: convertListResp(files),
	}, nil
}

func convertListResp(
	files []models.FileInfo,
) []*desc.FilesInfo {
	filesListResp := make([]*desc.FilesInfo, len(files))

	for idx, val := range files {
		filesListResp[idx] = convertModelFileToListResp(val)
	}

	return filesListResp
}

func convertModelFileToListResp(
	val models.FileInfo,
) *desc.FilesInfo {
	res := desc.FilesInfo{
		Filename:  val.FileName,
		FilePath:  val.FilePath,
		CreatedAt: timestamppb.New(val.CreatedAt),
		UpdatedAt: timestamppb.New(val.UpdatedAt.Time),
	}

	return &res
}
