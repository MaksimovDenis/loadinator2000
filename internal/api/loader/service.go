package loader

import (
	"github.com/MaksimovDenis/loadinator2000/internal/service"
	desc "github.com/MaksimovDenis/loadinator2000/pkg/loader_v1"
	"github.com/rs/zerolog"
)

type Implementation struct {
	desc.UnimplementedLoaderV1Server
	loaderSecrvice service.LoaderService
	log            zerolog.Logger
}

func NewImplementation(loaderSecrvice service.LoaderService, log zerolog.Logger) *Implementation {
	return &Implementation{
		loaderSecrvice: loaderSecrvice,
		log:            log,
	}
}
