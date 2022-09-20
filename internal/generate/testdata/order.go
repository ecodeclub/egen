package code

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type OrderDAO struct {
	session interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
}

type OrderTxDAO struct {
	*OrderDAO
}

func (dao *OrderTxDAO) Rollback() error {
	tx, ok := dao.session.(*sql.Tx)
	if !ok {
		return errors.New("非事务")
	}
	return tx.Rollback()
}

func (dao *OrderTxDAO) Commit() error {
	tx, ok := dao.session.(*sql.Tx)
	if !ok {
		return errors.New("非事务")
	}
	return tx.Commit()
}

func (dao *OrderDAO) Begin(ctx context.Context, opts *sql.TxOptions) (*OrderTxDAO, error) {
	db, ok := dao.session.(*sql.DB)
	if !ok {
		return nil, errors.New("不能在事务中开启事务")
	}
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &OrderTxDAO{
		OrderDAO: &OrderDAO{tx},
	}, nil
}

func NewOrderDAO(db *sql.DB) (*OrderDAO, error) {
	return &OrderDAO{db}, nil
}

func (dao *OrderDAO) Insert(ctx context.Context, vals ...*Order) (int64, error) {
	if len(vals) == 0 || vals == nil {
		return 0, nil
	}
	var args = make([]interface{}, 0, len(vals)*(6))
	var str = ""
	for k, v := range vals {
		if k != 0 {
			str += ", "
		}
		str += "(?,?,?,?,?,?)"
		args = append(args, v.OrderTime, v.OrderId, v.UserId, v.HasBuy, v.Price, v.Seller)
	}
	sqlSen := "INSERT INTO `order`(`order_time`,`order_id`,`user_id`,`has_buy`,`price`,`seller`) VALUES" + str
	res, err := dao.session.ExecContext(ctx, sqlSen, args...)
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
	err := row.Scan(&val.OrderTime, &val.OrderId, &val.UserId, &val.HasBuy, &val.Price, &val.Seller)
	return &val, err
}

func (dao *OrderDAO) SelectByRaw(ctx context.Context, query string, args ...any) (*Order, error) {
	row := dao.session.QueryRowContext(ctx, query, args...)
	return dao.NewOne(row)
}

func (dao *OrderDAO) SelectByWhere(ctx context.Context, where string, args ...any) (*Order, error) {
	s := "SELECT `order_time`,`order_id`,`user_id`,`has_buy`,`price`,`seller` FROM `order` WHERE " + where
	return dao.SelectByRaw(ctx, s, args...)
}

func (dao *OrderDAO) NewBatch(rows *sql.Rows) ([]*Order, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var vals = make([]*Order, 0, 6)
	for rows.Next() {
		var val Order
		if err := rows.Scan(&val.OrderTime, &val.OrderId, &val.UserId, &val.HasBuy, &val.Price, &val.Seller); err != nil {
			return nil, err
		}
		vals = append(vals, &val)
	}
	return vals, nil
}

func (dao *OrderDAO) SelectBatchByRaw(ctx context.Context, query string, args ...any) ([]*Order, error) {
	rows, err := dao.session.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return dao.NewBatch(rows)
}

func (dao *OrderDAO) SelectBatchByWhere(ctx context.Context, where string, args ...any) ([]*Order, error) {
	s := "SELECT `order_time`,`order_id`,`user_id`,`has_buy`,`price`,`seller` FROM `order` WHERE " + where
	return dao.SelectBatchByRaw(ctx, s, args...)
}

func (dao *OrderDAO) UpdateSpecificColsByWhere(ctx context.Context, val *Order, cols []string, where string, args ...any) (int64, error) {
	newArgs, colAfter := dao.quotedSpecificCol(val, cols...)
	newArgs = append(newArgs, args...)
	s := "UPDATE `order` SET " + colAfter + " WHERE " + where
	return dao.UpdateColsByRaw(ctx, s, newArgs...)
}

func (dao *OrderDAO) UpdateNoneZeroColByWhere(ctx context.Context, val *Order, where string, args ...any) (int64, error) {
	newArgs, colAfter := dao.quotedNoneZero(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `order` SET " + colAfter + " WHERE " + where
	return dao.UpdateColsByRaw(ctx, s, newArgs...)
}

func (dao *OrderDAO) UpdateNonePKColByWhere(ctx context.Context, val *Order, where string, args ...any) (int64, error) {
	newArgs, colAfter := dao.quotedNonePK(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `order` SET " + colAfter + " WHERE " + where
	return dao.UpdateColsByRaw(ctx, s, newArgs...)
}

func (dao *OrderDAO) quotedNoneZero(val *Order) ([]interface{}, string) {
	var cols = make([]string, 0, 6)
	var args = make([]interface{}, 0, 6)
	if val.OrderTime != "" {
		args = append(args, val.OrderTime)
		cols = append(cols, "`order_time`")
	}
	if val.OrderId != 0 {
		args = append(args, val.OrderId)
		cols = append(cols, "`order_id`")
	}
	if val.UserId != 0 {
		args = append(args, val.UserId)
		cols = append(cols, "`user_id`")
	}
	if val.HasBuy != false {
		args = append(args, val.HasBuy)
		cols = append(cols, "`has_buy`")
	}
	if val.Price != 0 {
		args = append(args, val.Price)
		cols = append(cols, "`price`")
	}
	if val.Seller != nil {
		args = append(args, val.Seller)
		cols = append(cols, "`seller`")
	}
	if len(cols) == 1 {
		cols[0] = cols[0] + "=?"
	} else {
		cols[len(cols)-1] = cols[len(cols)-1] + "=?"
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) quotedNonePK(val *Order) ([]interface{}, string) {
	var cols = []string{"`order_time`", "`order_id`", "`has_buy`", "`price`", "`seller`"}
	var args = []interface{}{val.OrderTime, val.OrderId, val.HasBuy, val.Price, val.Seller}
	if len(cols) == 1 {
		cols[0] = cols[0] + "=?"
	} else {
		cols[len(cols)-1] = cols[len(cols)-1] + "=?"
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) quotedSpecificCol(val *Order, cols ...string) ([]interface{}, string) {
	var relation = make(map[string]interface{}, 6)
	var args = make([]interface{}, 0, 6)
	relation["has_buy"] = val.HasBuy
	relation["order_id"] = val.OrderId
	relation["order_time"] = val.OrderTime
	relation["price"] = val.Price
	relation["seller"] = val.Seller
	relation["user_id"] = val.UserId
	for i := 0; i < len(cols); i++ {
		args = append(args, relation[cols[i]])
		cols[i] = "`" + cols[i] + "`"
	}
	if len(cols) == 1 {
		cols[0] = cols[0] + "=?"
	} else {
		cols[len(cols)-1] = cols[len(cols)-1] + "=?"
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) UpdateColsByRaw(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := dao.session.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (dao *OrderDAO) DeleteByWhere(ctx context.Context, where string, args ...any) (int64, error) {
	s := "DELETE FROM `order` WHERE " + where
	return dao.DeleteByRaw(ctx, s, args...)
}

func (dao *OrderDAO) DeleteByRaw(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := dao.session.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
