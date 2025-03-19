package loader

import (
	db "github.com/MaksimovDenis/loadinator2000/internal/client"
	"github.com/MaksimovDenis/loadinator2000/internal/repository"
	"github.com/MaksimovDenis/loadinator2000/internal/service"
	"github.com/rs/zerolog"
)

const (
	maxDownloadRequests = 10
	maxListRequests     = 100
)

type serv struct {
	loaderRepository repository.LoaderRepository
	txManager        db.TxManager
	log              zerolog.Logger
	downloadLimiter  chan struct{}
	listLimiter      chan struct{}
}

func NewService(
	loaderRepository repository.LoaderRepository,
	txManager db.TxManager,
	log zerolog.Logger,
) service.LoaderService {
	return &serv{
		loaderRepository: loaderRepository,
		txManager:        txManager,
		log:              log,
		downloadLimiter:  make(chan struct{}, maxDownloadRequests),
		listLimiter:      make(chan struct{}, maxListRequests),
	}
}
