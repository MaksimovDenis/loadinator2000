package loader

import (
	"context"

	desc "github.com/MaksimovDenis/loadinator2000/pkg/loader_v1"
)

func (hdl *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	data, err := hdl.loaderSecrvice.Get(ctx, req.Filename)
	if err != nil {
		hdl.log.Error().Err(err).Msgf("failed to get file %s", req.Filename)
		return nil, err
	}

	return &desc.GetResponse{
		Data: data,
	}, nil
}
