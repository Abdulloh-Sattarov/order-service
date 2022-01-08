package grpcclient

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/abdullohsattorov/order-service/config"
	pb "github.com/abdullohsattorov/order-service/genproto/catalog_service"
)

// IServiceManager ...
type IServiceManager interface {
	CatalogService() pb.CatalogServiceClient
}

// IServiceManager ...
type serviceManager struct {
	cfg            config.Config
	catalogService pb.CatalogServiceClient
}

// New ...
func New(cfg config.Config) (IServiceManager, error) {
	connCatalog, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.CatalogServiceHost, cfg.CatalogServicePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("catalog service dial host: %s port: %d", cfg.CatalogServiceHost, cfg.CatalogServicePort)
	}

	serviceManager := &serviceManager{
		cfg:            cfg,
		catalogService: pb.NewCatalogServiceClient(connCatalog),
	}

	return serviceManager, nil
}

// CatalogService ...
func (s *serviceManager) CatalogService() pb.CatalogServiceClient {
	return s.catalogService
}
