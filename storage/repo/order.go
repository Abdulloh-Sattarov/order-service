package repo

import (
	pb "github.com/abdullohsattorov/order-service/genproto/order_service"
)

// OrderStorageI ...
type OrderStorageI interface {
	Create(pb.Order) (pb.Order, error)
	Get(id string) (pb.Order, error)
	List(page, limit int64) ([]*pb.Order, int64, error)
	Update(pb.Order) (pb.Order, error)
	Delete(id string) error
}
