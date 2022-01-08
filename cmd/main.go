package main

import (
	"net"

	"github.com/abdullohsattorov/order-service/config"
	pb "github.com/abdullohsattorov/order-service/genproto/order_service"
	"github.com/abdullohsattorov/order-service/pkg/db"
	"github.com/abdullohsattorov/order-service/pkg/logger"
	"github.com/abdullohsattorov/order-service/service"
	grpcClient "github.com/abdullohsattorov/order-service/service/grpc_client"
	"github.com/abdullohsattorov/order-service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "order-service")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			log.Fatal("failed cleanup logger", logger.Error(err))
		}
	}(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	client, err := grpcClient.New(cfg)
	if err != nil {
		log.Fatal("grpc client connection error", logger.Error(err))
	}

	pgStorage := storage.NewStoragePg(connDB)

	orderService := service.NewOrderService(pgStorage, log, client, &cfg)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, orderService)

	reflection.Register(s)

	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
