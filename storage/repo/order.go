package repo

import (
	pb "github.com/abdullohsattorov/order-service/genproto/order_service"
)

// OrderStorageI ...
type OrderStorageI interface {
	Create(pb.OrderReq) (pb.OrderResp, error)
	Get(id string) (pb.OrderResp, error)
	List(page, limit int64) ([]*pb.OrderResp, int64, error)
	Update(pb.OrderReq) (pb.OrderResp, error)
	Delete(id string) error
}
