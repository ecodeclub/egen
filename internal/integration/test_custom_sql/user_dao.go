package genDAO

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type UserGenDAO struct {
	session interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
}

type UserTxGenDAO struct {
	*UserGenDAO
}

func (dao *UserTxGenDAO) Rollback() error {
	tx, ok := dao.session.(*sql.Tx)
	if !ok {
		return errors.New("非事务")
	}
	return tx.Rollback()
}

func (dao *UserTxGenDAO) Commit() error {
	tx, ok := dao.session.(*sql.Tx)
	if !ok {
		return errors.New("非事务")
	}
	return tx.Commit()
}

func (dao *UserGenDAO) Begin(ctx context.Context, opts *sql.TxOptions) (*UserTxGenDAO, error) {
	db, ok := dao.session.(*sql.DB)
	if !ok {
		return nil, errors.New("不能在事务中开启事务")
	}
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &UserTxGenDAO{
		UserGenDAO: &UserGenDAO{tx},
	}, nil
}

func NewUserGenDAO(db *sql.DB) (*UserGenDAO, error) {
	return &UserGenDAO{db}, nil
}

func (dao *UserGenDAO) Insert(ctx context.Context, vals ...*User) (int64, error) {
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
		args = append(args, v.ID, v.Username, v.Password, v.Login, v.Status, v.Money)
	}
	sqlSen := "INSERT INTO `user_dao`(`id`,`username`,`password`,`login`,`status`,`money`) VALUES" + str
	res, err := dao.session.ExecContext(ctx, sqlSen, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (dao *UserGenDAO) NewOne(row *sql.Row) (*User, error) {
	if err := row.Err(); err != nil {
		return nil, err
	}
	var val User
	err := row.Scan(&val.ID, &val.Username, &val.Password, &val.Login, &val.Status, &val.Money)
	return &val, err
}

func (dao *UserGenDAO) SelectByRaw(ctx context.Context, query string, args ...any) (*User, error) {
	row := dao.session.QueryRowContext(ctx, query, args...)
	return dao.NewOne(row)
}

func (dao *UserGenDAO) SelectByWhere(ctx context.Context, where string, args ...any) (*User, error) {
	s := "SELECT `id`,`username`,`password`,`login`,`status`,`money` FROM `user_dao` WHERE " + where
	return dao.SelectByRaw(ctx, s, args...)
}

func (dao *UserGenDAO) NewBatch(rows *sql.Rows) ([]*User, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var vals = make([]*User, 0, 6)
	for rows.Next() {
		var val User
		if err := rows.Scan(&val.ID, &val.Username, &val.Password, &val.Login, &val.Status, &val.Money); err != nil {
			return nil, err
		}
		vals = append(vals, &val)
	}
	return vals, nil
}

func (dao *UserGenDAO) SelectBatchByRaw(ctx context.Context, query string, args ...any) ([]*User, error) {
	rows, err := dao.session.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return dao.NewBatch(rows)
}

func (dao *UserGenDAO) SelectBatchByWhere(ctx context.Context, where string, args ...any) ([]*User, error) {
	s := "SELECT `id`,`username`,`password`,`login`,`status`,`money` FROM `user_dao` WHERE " + where
	return dao.SelectBatchByRaw(ctx, s, args...)
}

func (dao *UserGenDAO) UpdateSpecificColsByWhere(ctx context.Context, val *User, cols []string, where string, args ...any) (int64, error) {
	newArgs, colAfter := dao.quotedSpecificCol(val, cols...)
	newArgs = append(newArgs, args...)
	s := "UPDATE `user_dao` SET " + colAfter + " WHERE " + where
	return dao.UpdateColsByRaw(ctx, s, newArgs...)
}

func (dao *UserGenDAO) UpdateNoneZeroColByWhere(ctx context.Context, val *User, where string, args ...any) (int64, error) {
	newArgs, colAfter := dao.quotedNoneZero(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `user_dao` SET " + colAfter + " WHERE " + where
	return dao.UpdateColsByRaw(ctx, s, newArgs...)
}

func (dao *UserGenDAO) UpdateNonePKColByWhere(ctx context.Context, val *User, where string, args ...any) (int64, error) {
	newArgs, colAfter := dao.quotedNonePK(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `user_dao` SET " + colAfter + " WHERE " + where
	return dao.UpdateColsByRaw(ctx, s, newArgs...)
}

func (dao *UserGenDAO) quotedNoneZero(val *User) ([]interface{}, string) {
	var cols = make([]string, 0, 6)
	var args = make([]interface{}, 0, 6)
	if val.ID != 0 {
		args = append(args, val.ID)
		cols = append(cols, "`id`")
	}
	if val.Username != "" {
		args = append(args, val.Username)
		cols = append(cols, "`username`")
	}
	if val.Password != "" {
		args = append(args, val.Password)
		cols = append(cols, "`password`")
	}
	if val.Login != "" {
		args = append(args, val.Login)
		cols = append(cols, "`login`")
	}
	if val.Status {
		args = append(args, val.Status)
		cols = append(cols, "`status`")
	}
	if val.Money != 0 {
		args = append(args, val.Money)
		cols = append(cols, "`money`")
	}
	if len(cols) == 1 {
		cols[0] = cols[0] + "=?"
	} else {
		cols[len(cols)-1] = cols[len(cols)-1] + "=?"
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *UserGenDAO) quotedNonePK(val *User) ([]interface{}, string) {
	var cols = []string{"`username`", "`password`", "`login`", "`status`", "`money`"}
	var args = []interface{}{val.Username, val.Password, val.Login, val.Status, val.Money}
	if len(cols) == 1 {
		cols[0] = cols[0] + "=?"
	} else {
		cols[len(cols)-1] = cols[len(cols)-1] + "=?"
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *UserGenDAO) quotedSpecificCol(val *User, cols ...string) ([]interface{}, string) {
	var relation = make(map[string]interface{}, 6)
	var args = make([]interface{}, 0, 6)
	relation["id"] = val.ID
	relation["login"] = val.Login
	relation["money"] = val.Money
	relation["password"] = val.Password
	relation["status"] = val.Status
	relation["username"] = val.Username
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

func (dao *UserGenDAO) UpdateColsByRaw(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := dao.session.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (dao *UserGenDAO) DeleteByWhere(ctx context.Context, where string, args ...any) (int64, error) {
	s := "DELETE FROM `user_dao` WHERE " + where
	return dao.DeleteByRaw(ctx, s, args...)
}

func (dao *UserGenDAO) DeleteByRaw(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := dao.session.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (dao *UserGenDAO) FindExist(ctx context.Context, ids []uint32, status bool) ([]*User, error) {
	var params = make([]any, 0)
	Sql := "SELECT `id`, `username`, `password`, `login`, `status`, `money` FROM `user_dao` WHERE `id` IN #ids AND `status` = #status"
	replaceParams := ""
	for i := 0; i < len(ids); i++ {
		if i != 0 {
			replaceParams += ","
		}
		replaceParams += "?"
	}
	Sql = strings.Replace(Sql, "#ids", "("+replaceParams+")", 1)
	for _, v := range ids {
		params = append(params, v)
	}
	Sql = strings.Replace(Sql, "#status", "?", 1)
	params = append(params, status)
	ret := make([]*User, 0, 20)
	rows, err := dao.session.QueryContext(ctx, Sql, params...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var one User
		if err = rows.Scan(&one.ID, &one.Username, &one.Password, &one.Login, &one.Status, &one.Money); err != nil {
			return nil, err
		}
		ret = append(ret, &one)
	}
	return ret, err
}

func (dao *UserGenDAO) UpdateOne(ctx context.Context, name string, status bool) (int64, error) {
	var params = make([]any, 0)
	Sql := "UPDATE `user_dao` SET `username`=#name WHERE `status`=#status"
	Sql = strings.Replace(Sql, "#name", "?", 1)
	params = append(params, name)
	Sql = strings.Replace(Sql, "#status", "?", 1)
	params = append(params, status)
	ret, err := dao.UpdateColsByRaw(ctx, Sql, params...)
	return ret, err
}

func (dao *UserGenDAO) DeleteOne(ctx context.Context, status bool, ids []uint32) (int64, error) {
	var params = make([]any, 0)
	Sql := "DELETE FROM `user_dao` WHERE `status`=#status AND `id` NOT IN #ids"
	Sql = strings.Replace(Sql, "#status", "?", 1)
	params = append(params, status)
	replaceParams := ""
	for i := 0; i < len(ids); i++ {
		if i != 0 {
			replaceParams += ","
		}
		replaceParams += "?"
	}
	Sql = strings.Replace(Sql, "#ids", "("+replaceParams+")", 1)
	for _, v := range ids {
		params = append(params, v)
	}
	ret, err := dao.DeleteByRaw(ctx, Sql, params...)
	return ret, err
}

func (dao *UserGenDAO) GetCount(ctx context.Context) (int, error) {
	var params = make([]any, 0)
	Sql := "SELECT COUNT(*) FROM `user_dao`"
	var ret int
	row := dao.session.QueryRowContext(ctx, Sql, params...)
	err := row.Scan(&ret)
	return ret, err
}

func (dao *UserGenDAO) GetTotalMoney(ctx context.Context) (float64, error) {
	var params = make([]any, 0)
	Sql := "SELECT SUM(`money`) FROM `user_dao`"
	var ret float64
	row := dao.session.QueryRowContext(ctx, Sql, params...)
	err := row.Scan(&ret)
	return ret, err
}

func (dao *UserGenDAO) GetAllStatus(ctx context.Context) ([]bool, error) {
	var params = make([]any, 0)
	Sql := "SELECT `status` FROM `user_dao`"
	ret := make([]bool, 0, 10)
	rows, err := dao.session.QueryContext(ctx, Sql, params...)
	if err != nil {
		return ret, err
	}
	for rows.Next() {
		var one bool
		if err = rows.Scan(&one); err != nil {
			return ret, err
		}
		ret = append(ret, one)
	}
	return ret, err
}
