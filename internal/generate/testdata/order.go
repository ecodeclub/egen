package code

import (
	"context"
	"database/sql"
)

type OrderDAO struct {
	DB *sql.DB
}

func (dao *OrderDAO) Insert(ctx context.Context, vals ...*Order) (int64, error) {
	var args = make([]interface{}, len(vals)*(3))
	var str = ""
	for k, v := range vals {
		if k != 0 {
			str += ", "
		}
		str += "(?,?,?)"
		args = append(args, v.OrderTime, v.OrderId, v.UserId)
	}
	sqlSen := "INSERT INTO `order`(`order_time`,`order_id`,`user_id`) VALUES" + str
	res, err := dao.DB.ExecContext(ctx, sqlSen, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (dao *OrderDAO) NewOne(row *sql.Row) (*Order, error) {
	if err := row.Err(); err != nil {
		return nil, err
	}
	var val Order
	err := row.Scan(&val.OrderTime, &val.OrderId, &val.UserId)
	return &val, err
}

func (dao *OrderDAO) SelectByRaw(ctx context.Context, query string, args ...any) (*Order, error) {
	row := dao.DB.QueryRowContext(ctx, query, args...)
	return dao.NewOne(row)
}

func (dao *OrderDAO) SelectByWhere(ctx context.Context, where string, args ...any) (*Order, error) {
	s := "SELECT `order_time`,`order_id`,`user_id` FROM `order` WHERE " + where
	return dao.SelectByRaw(ctx, s, args...)
}

func (dao *OrderDAO) NewBatch(rows *sql.Rows) ([]*Order, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var vals = make([]*Order, 0, 3)
	for rows.Next() {
		var val Order
		if err := rows.Scan(&val.OrderTime, &val.OrderId, &val.UserId); err != nil {
			return nil, err
		}
		vals = append(vals, &val)
	}
	return vals, nil
}

func (dao *OrderDAO) SelectBatchByRaw(ctx context.Context, query string, args ...any) ([]*Order, error) {
	rows, err := dao.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return dao.NewBatch(rows)
}

func (dao *OrderDAO) SelectBatchByWhere(ctx context.Context, where string, args ...any) ([]*Order, error) {
	s := "SELECT `order_time`,`order_id`,`user_id` FROM `order` WHERE " + where
	return dao.SelectBatchByRaw(ctx, s, args...)
}
