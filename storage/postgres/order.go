package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	pb "github.com/abdullohsattorov/order-service/genproto/order_service"
)

type orderRepo struct {
	db *sqlx.DB
}

// New Repo
func NewOrderRepo(db *sqlx.DB) *orderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(order pb.OrderReq) (pb.OrderResp, error) {
	var orderId string
	err := r.db.QueryRow(`
		INSERT INTO orders (order_id, book_uuid, description)
		VALUES ($1, $2, $3) returning order_id`, order.OrderId, order.BookId, order.Description).Scan(&orderId)
	if err != nil {
		return pb.OrderResp{}, err
	}

	var NewOrder pb.OrderResp

	NewOrder, err = r.Get(orderId)

	if err != nil {
		return pb.OrderResp{}, err
	}

	return NewOrder, nil
}

func (r *orderRepo) Get(orderId string) (pb.OrderResp, error) {
	var order pb.OrderResp

	err := r.db.QueryRow(`
		SELECT order_id, book_uuid, description, created_at, updated_at FROM orders WHERE deleted_at IS NULL AND order_id = $1`,
		orderId).Scan(&order.OrderId, &order.BookId, &order.Description, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return pb.OrderResp{}, err
	}

	return order, err
}

func (r *orderRepo) List(page, limit int64) ([]*pb.OrderResp, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Queryx(`
		SELECT order_id, book_uuid, description, created_at, updated_at FROM orders where deleted_at is null LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:err check

	var (
		orders []*pb.OrderResp
		count  int64
	)
	for rows.Next() {
		var order pb.OrderResp
		err = rows.Scan(&order.OrderId, &order.BookId, &order.Description, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		orders = append(orders, &order)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM orders where deleted_at is null`).Scan(&count)

	if err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

func (r *orderRepo) Update(order pb.OrderReq) (pb.OrderResp, error) {
	result, err := r.db.Exec(`UPDATE orders SET book_uuid=$1, description=$2, updated_at=$3 WHERE order_id=$4`, order.BookId, order.Description, time.Now().UTC(), order.OrderId)
	if err != nil {
		return pb.OrderResp{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.OrderResp{}, sql.ErrNoRows
	}

	var NewOrder pb.OrderResp

	NewOrder, err = r.Get(order.OrderId)

	fmt.Println(result, NewOrder)

	if err != nil {
		return pb.OrderResp{}, err
	}

	return NewOrder, nil
}

func (r *orderRepo) Delete(id string) error {
	result, err := r.db.Exec(`UPDATE orders SET deleted_at=$1 WHERE order_id=$2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
