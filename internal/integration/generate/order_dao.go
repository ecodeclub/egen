//go:build e2e

package generate

import (
	"context"
	"database/sql"
	"strings"
)

type handler interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type handlerTx interface {
	Rollback() error
	Commit() error
}

type OrderDAO struct {
	handler
	handlerTx
}

func (dao *OrderDAO) initDatabase(operator interface{}) {
	switch x := operator.(type) {
	case handler:
		dao.handlerTx = nil
		dao.handler = x
	case handlerTx:
		// 同时赋值
		dao.handlerTx = x
		dao.handler = x.(handler)
	}
}

// 调用提交及回滚全权交给用户
func (dao *OrderDAO) commit() error {
	if dao.handlerTx != nil {
		return dao.handlerTx.Commit()
	}
	// 不存在调用 不应报错
	return nil
}

func (dao *OrderDAO) rollback() error {
	if dao.handlerTx != nil {
		return dao.handlerTx.Rollback()
	}
	return nil
}

func (dao *OrderDAO) Insert(ctx context.Context, vals ...*Order) (int64, error) {
	if len(vals) == 0 || vals == nil {
		return 0, nil
	}
	var args = make([]interface{}, 0, len(vals)*(3))
	var str = ""
	for k, v := range vals {
		if k != 0 {
			str += ", "
		}
		str += "(?,?,?)"
		args = append(args, v.UserId, v.OrderId, v.Price)
	}
	sqlSen := "INSERT INTO `order`(`user_id`,`order_id`,`price`) VALUES" + str
	res, err := dao.handler.ExecContext(ctx, sqlSen, args...)
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
	err := row.Scan(&val.UserId, &val.OrderId, &val.Price)
	return &val, err
}

func (dao *OrderDAO) SelectByRaw(ctx context.Context, query string, args ...any) (*Order, error) {
	row := dao.handler.QueryRowContext(ctx, query, args...)
	return dao.NewOne(row)
}

func (dao *OrderDAO) SelectByWhere(ctx context.Context, where string, args ...any) (*Order, error) {
	s := "SELECT `user_id`,`order_id`,`price` FROM `order` WHERE " + where
	return dao.SelectByRaw(ctx, s, args...)
}

func (dao *OrderDAO) NewBatch(rows *sql.Rows) ([]*Order, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var vals = make([]*Order, 0, 3)
	for rows.Next() {
		var val Order
		if err := rows.Scan(&val.UserId, &val.OrderId, &val.Price); err != nil {
			return nil, err
		}
		vals = append(vals, &val)
	}
	return vals, nil
}

func (dao *OrderDAO) SelectBatchByRaw(ctx context.Context, query string, args ...any) ([]*Order, error) {
	rows, err := dao.handler.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return dao.NewBatch(rows)
}

func (dao *OrderDAO) SelectBatchByWhere(ctx context.Context, where string, args ...any) ([]*Order, error) {
	s := "SELECT `user_id`,`order_id`,`price` FROM `order` WHERE " + where
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
	var cols = make([]string, 0, 3)
	var args = make([]interface{}, 0, 3)
	if val.UserId != 0 {
		args = append(args, val.UserId)
		cols = append(cols, "`user_id`")
	}
	if val.OrderId != 0 {
		args = append(args, val.OrderId)
		cols = append(cols, "`order_id`")
	}
	if val.Price != 0 {
		args = append(args, val.Price)
		cols = append(cols, "`price`")
	}
	if len(cols) == 1 {
		cols[0] = cols[0] + "=?"
	} else {
		cols[len(cols)-1] = cols[len(cols)-1] + "=?"
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) quotedNonePK(val *Order) ([]interface{}, string) {
	var cols = []string{
		"`order_id`",
		"`price`"}
	var args = []interface{}{val.OrderId, val.Price}
	if len(cols) == 1 {
		cols[0] = cols[0] + "=?"
	} else {
		cols[len(cols)-1] = cols[len(cols)-1] + "=?"
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) quotedSpecificCol(val *Order, cols ...string) ([]interface{}, string) {
	var relation = make(map[string]interface{}, 3)
	var args = make([]interface{}, 0, 3)
	relation["order_id"] = val.OrderId
	relation["price"] = val.Price
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
	res, err := dao.handler.ExecContext(ctx, query, args...)
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
	res, err := dao.handler.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
