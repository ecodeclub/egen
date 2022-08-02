package code

import (
	"context"
	"database/sql"
	"strings"
)

type UserDAO struct {
	DB *sql.DB
}

func (dao *UserDAO) Insert(ctx context.Context, vals ...*User) (int64, error) {
	var args = make([]interface{}, len(vals)*(5))
	var str = ""
	for k, v := range vals {
		if k != 0 {
			str += ", "
		}
		str += "(?,?,?,?,?)"
		args = append(args, v.LoginTime, v.FirstName, v.LastName, v.UserId, v.Password)
	}
	sqlSen := "INSERT INTO `user`(`login_time`,`first_name`,`last_name`,`user_id`,`password`) VALUES" + str
	res, err := dao.DB.ExecContext(ctx, sqlSen, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (dao *UserDAO) NewOne(row *sql.Row) (*User, error) {
	if err := row.Err(); err != nil {
		return nil, err
	}
	var val User
	err := row.Scan(&val.LoginTime, &val.FirstName, &val.LastName, &val.UserId, &val.Password)
	return &val, err
}

func (dao *UserDAO) SelectByRaw(ctx context.Context, query string, args ...any) (*User, error) {
	row := dao.DB.QueryRowContext(ctx, query, args...)
	return dao.NewOne(row)
}

func (dao *UserDAO) SelectByWhere(ctx context.Context, where string, args ...any) (*User, error) {
	s := "SELECT `login_time`,`first_name`,`last_name`,`user_id`,`password` FROM `user` WHERE " + where
	return dao.SelectByRaw(ctx, s, args...)
}

func (dao *UserDAO) NewBatch(rows *sql.Rows) ([]*User, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var vals = make([]*User, 0, 5)
	for rows.Next() {
		var val User
		if err := rows.Scan(&val.LoginTime, &val.FirstName, &val.LastName, &val.UserId, &val.Password); err != nil {
			return nil, err
		}
		vals = append(vals, &val)
	}
	return vals, nil
}

func (dao *UserDAO) SelectBatchByRaw(ctx context.Context, query string, args ...any) ([]*User, error) {
	rows, err := dao.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return dao.NewBatch(rows)
}

func (dao *UserDAO) SelectBatchByWhere(ctx context.Context, where string, args ...any) ([]*User, error) {
	s := "SELECT `login_time`,`first_name`,`last_name`,`user_id`,`password` FROM `user` WHERE " + where
	return dao.SelectBatchByRaw(ctx, s, args...)
}

func (dao *UserDAO) UpdateColsByWhere(ctx context.Context, val *User, cols []string, where string, args ...any) (int64, error) {
	newArgs, colAfter := dao.quotedSpecificCol(val, cols...)
	newArgs = append(newArgs, args...)
	s := "UPDATE `user` SET " + colAfter + " WHERE " + where
	return dao.UpdateColsByRaw(ctx, s, newArgs...)
}

func (dao *UserDAO) quotedNoneZero(val *User) ([]interface{}, string) {
	var cols = make([]string, 0, 5)
	var args = make([]interface{}, 0, 5)
	if val.LoginTime != "" {
		args = append(args, val.LoginTime)
		cols = append(cols, "`login_time`")
	}
	if val.FirstName != "" {
		args = append(args, val.FirstName)
		cols = append(cols, "`first_name`")
	}
	if val.LastName != "" {
		args = append(args, val.LastName)
		cols = append(cols, "`last_name`")
	}
	if val.UserId != 0 {
		args = append(args, val.UserId)
		cols = append(cols, "`user_id`")
	}
	if val.Password != nil {
		args = append(args, val.Password)
		cols = append(cols, "`password`")
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *UserDAO) quotedNonePK(val *User) ([]interface{}, string) {
	var cols = []string{`login_time`, `first_name`, `last_name`, `password`}
	var args = []interface{}{val.LoginTime, val.FirstName, val.LastName, val.Password}
	return args, strings.Join(cols, "=?")
}

func (dao *UserDAO) quotedSpecificCol(val *User, cols ...string) ([]interface{}, string) {
	var relation = make(map[string]interface{}, 5)
	var args = make([]interface{}, 0, 5)
	relation["first_name"] = val.FirstName
	relation["last_name"] = val.LastName
	relation["login_time"] = val.LoginTime
	relation["password"] = val.Password
	relation["user_id"] = val.UserId
	for i := 0; i < len(cols); i++ {
		args = append(args, relation[cols[i]])
		cols[i] = "`" + cols[i] + "`"
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *UserDAO) UpdateColsByRaw(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := dao.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (dao *UserDAO) DeleteByWhere(ctx context.Context, where string, args ...any) (int64, error) {
	s := "DELETE FROM `user` WHERE " + where
	return dao.DeleteByRaw(ctx, s, args...)
}

func (dao *UserDAO) DeleteByRaw(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := dao.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
