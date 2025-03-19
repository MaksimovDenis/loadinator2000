package loader

import (
	"context"

	desc "github.com/MaksimovDenis/loadinator2000/pkg/loader_v1"
)

func (hdl *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	fileName, err := hdl.loaderSecrvice.Create(ctx, req.Filename, req.FilePath, req.Data)
	if err != nil {
		hdl.log.Error().Err(err).Msgf("failed to add file %s", req.Filename)
		return nil, err
	}

	return &desc.CreateResponse{
		Message: fileName,
	}, nil
}
