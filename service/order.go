package service

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gofrs/uuid"

	"github.com/abdullohsattorov/order-service/config"
	c "github.com/abdullohsattorov/order-service/genproto/catalog_service"
	pb "github.com/abdullohsattorov/order-service/genproto/order_service"
	l "github.com/abdullohsattorov/order-service/pkg/logger"
	grpcclient "github.com/abdullohsattorov/order-service/service/grpc_client"
	"github.com/abdullohsattorov/order-service/storage"
)

// OrderService ...
type OrderService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcclient.IServiceManager
	config  *config.Config
}

// NewOrderService ...
func NewOrderService(storage storage.IStorage, log l.Logger, client grpcclient.IServiceManager, config *config.Config) *OrderService {
	return &OrderService{
		storage: storage,
		logger:  log,
		client:  client,
		config:  config,
	}
}

func (s *OrderService) Create(ctxReq context.Context, req *pb.OrderReq) (*pb.OrderResp, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return &pb.OrderResp{}, status.Error(codes.Internal, "failed generate uuid")
	}

	req.OrderId = id.String()

	order, err := s.storage.Order().Create(*req)
	if err != nil {
		s.logger.Error("failed to create order", l.Error(err))
		return &pb.OrderResp{}, status.Error(codes.Internal, "failed to create order")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	catalogBook, err := s.client.CatalogService().GetBook(ctx, &c.ByIdReq{Id: order.BookId})
	if err != nil {
		s.logger.Error("failed while getting book", l.Error(err))
		return &pb.OrderResp{}, status.Error(codes.Internal, "failed while getting book")
	}

	catalogAuthor, err := s.client.CatalogService().GetAuthor(ctx, &c.ByIdReq{Id: catalogBook.AuthorId})
	if err != nil {
		s.logger.Error("failed while getting author", l.Error(err))
		return &pb.OrderResp{}, status.Error(codes.Internal, "failed to while getting author")
	}

	order.BookName = catalogBook.Name
	order.AuthorId = catalogAuthor.AuthorId
	order.AuthorName = catalogAuthor.Name
	return &order, nil
}

func (s *OrderService) Get(ctxReq context.Context, req *pb.ByIdReq) (*pb.OrderResp, error) {
	order, err := s.storage.Order().Get(req.GetId())
	if err != nil {
		s.logger.Error("failed to get order", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get order")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	catalogBook, err := s.client.CatalogService().GetBook(ctx, &c.ByIdReq{Id: order.BookId})
	if err != nil {
		s.logger.Error("failed while getting book", l.Error(err))
		return &pb.OrderResp{}, status.Error(codes.Internal, "failed while getting book")
	}

	catalogAuthor, err := s.client.CatalogService().GetAuthor(ctx, &c.ByIdReq{Id: catalogBook.AuthorId})
	if err != nil {
		s.logger.Error("failed while getting author", l.Error(err))
		return &pb.OrderResp{}, status.Error(codes.Internal, "failedto while getting author")
	}
	order.BookName = catalogBook.Name
	order.AuthorId = catalogAuthor.AuthorId
	order.AuthorName = catalogAuthor.Name

	return &order, nil
}

func (s *OrderService) List(ctxReq context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	orders, count, err := s.storage.Order().List(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to list orders", l.Error(err))
		return &pb.ListResp{}, status.Error(codes.Internal, "failed to list orders")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	for i, order := range orders {
		catalogBook, err := s.client.CatalogService().GetBook(ctx, &c.ByIdReq{Id: order.BookId})
		if err != nil {
			s.logger.Error("failed while getting book", l.Error(err))
			return &pb.ListResp{}, status.Error(codes.Internal, "failedto while getting book")
		}

		catalogAuthor, err := s.client.CatalogService().GetAuthor(ctx, &c.ByIdReq{Id: catalogBook.AuthorId})
		if err != nil {
			s.logger.Error("failed while getting author", l.Error(err))
			return &pb.ListResp{}, status.Error(codes.Internal, "failedto while getting author")
		}
		order.BookName = catalogBook.Name
		order.AuthorId = catalogAuthor.AuthorId
		order.AuthorName = catalogAuthor.Name

		orders[i] = order
	}
	return &pb.ListResp{
		Orders: orders,
		Count:  count,
	}, nil
}

func (s *OrderService) Update(ctxReq context.Context, req *pb.OrderReq) (*pb.OrderResp, error) {
	order, err := s.storage.Order().Update(*req)
	if err != nil {
		s.logger.Error("failed to update order", l.Error(err))
		return nil, status.Error(codes.Internal, "failedto update order")
	}

	return &order, nil
}

func (s *OrderService) Delete(ctxReq context.Context, req *pb.ByIdReq) (*pb.EmptyResp, error) {
	err := s.storage.Order().Delete(req.Id)
	if err != nil {
		s.logger.Error("failed to delete order", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete order")
	}

	return &pb.EmptyResp{}, nil
}
