package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/abdullohsattorov/order-service/storage/postgres"
	"github.com/abdullohsattorov/order-service/storage/repo"
)

// IStorage ...
type IStorage interface {
	Order() repo.OrderStorageI
}

type storagePg struct {
	db        *sqlx.DB
	orderRepo repo.OrderStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:        db,
		orderRepo: postgres.NewOrderRepo(db),
	}
}

func (s storagePg) Order() repo.OrderStorageI {
	return s.orderRepo
}
