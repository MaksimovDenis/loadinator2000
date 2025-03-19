package app

import (
	"context"
	"os"

	"github.com/MaksimovDenis/loadinator2000/internal/api/loader"
	db "github.com/MaksimovDenis/loadinator2000/internal/client"
	"github.com/MaksimovDenis/loadinator2000/internal/client/db/pg"
	"github.com/MaksimovDenis/loadinator2000/internal/client/db/transaction"
	"github.com/MaksimovDenis/loadinator2000/internal/closer"
	"github.com/MaksimovDenis/loadinator2000/internal/config"
	"github.com/MaksimovDenis/loadinator2000/internal/repository"
	loaderRepository "github.com/MaksimovDenis/loadinator2000/internal/repository/loader"
	loaderService "github.com/MaksimovDenis/loadinator2000/internal/service/loader"

	"github.com/MaksimovDenis/loadinator2000/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient         db.Client
	txManager        db.TxManager
	loaderRepository repository.LoaderRepository

	loaderService service.LoaderService

	log zerolog.Logger

	loaderImpl *loader.Implementation
}

func newServiceProvider() *serviceProvider {
	srv := &serviceProvider{}
	srv.log = srv.initLogger()

	return srv
}

func (srv *serviceProvider) initLogger() zerolog.Logger {
	logFile, err := os.OpenFile("./internal/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open log file")
	}

	logLevel, err := zerolog.ParseLevel("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	multiWriter := zerolog.MultiLevelWriter(os.Stdout, logFile)

	logger := zerolog.New(multiWriter).Level(logLevel).With().Timestamp().Logger()

	return logger
}

func (srv *serviceProvider) PGConfig() config.PGConfig {
	if srv.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get pg config")
		}

		srv.pgConfig = cfg
	}

	return srv.pgConfig
}

func (srv *serviceProvider) GRPCConfig() config.GRPCConfig {
	if srv.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get server config")
		}

		srv.grpcConfig = cfg
	}

	return srv.grpcConfig
}

func (srv *serviceProvider) DBClient(ctx context.Context) db.Client {
	if srv.dbClient == nil {
		client, err := pg.New(ctx, srv.PGConfig().DSN())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create db client")
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("ping error")
		}

		closer.Add(func() error {
			client.Close()
			return nil
		})

		srv.dbClient = client
	}

	return srv.dbClient
}

func (srv *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if srv.txManager == nil {
		srv.txManager = transaction.NewTransactionsManager(srv.DBClient(ctx).DB())
	}

	return srv.txManager
}

func (srv *serviceProvider) LoaderRepository(ctx context.Context) repository.LoaderRepository {
	if srv.loaderRepository == nil {
		srv.loaderRepository = loaderRepository.NewRepository(srv.DBClient(ctx), srv.log)
	}

	return srv.loaderRepository
}

func (srv *serviceProvider) LoaderService(ctx context.Context) service.LoaderService {
	if srv.loaderService == nil {
		srv.loaderService = loaderService.NewService(
			srv.LoaderRepository(ctx),
			srv.TxManager(ctx),
			srv.log)
	}

	return srv.loaderService
}

func (srv *serviceProvider) LoadImpl(ctx context.Context) *loader.Implementation {
	if srv.loaderImpl == nil {
		srv.loaderImpl = loader.NewImplementation(srv.LoaderService(ctx), srv.log)
	}

	return srv.loaderImpl
}
