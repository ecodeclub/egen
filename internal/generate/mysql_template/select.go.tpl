

func (dao *{{.GoName}}DAO) NewOne(row *sql.Row) (*{{.GoName}}, error) {
	if err := row.Err(); err != nil {
		return nil, err
	}
	var val {{.GoName}}
	err := row.Scan({{.QuotedExecArgsWithParameter "&" "val" .QuotedAllCol}})
	return &val, err
}

func (dao *{{.GoName}}DAO) SelectByRaw(ctx context.Context, query string, args ...any) (*{{.GoName}}, error) {
	row := dao.DB.QueryRowContext(ctx, query, args...)
	return dao.NewOne(row)
}

func (dao *{{.GoName}}DAO) SelectByWhere(ctx context.Context, where string, args ...any) (*{{.GoName}}, error) {
	s := "SELECT {{.QuotedAllCol}} FROM {{.QuotedTableName}} WHERE " + where
	return dao.SelectByRaw(ctx, s, args...)
}

func (dao *{{.GoName}}DAO) NewBatch(rows *sql.Rows) ([]*{{.GoName}}, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var vals = make([]*{{.GoName}}, 0, {{len .Fields}})
	for rows.Next() {
		var val {{.GoName}}
		if err := rows.Scan({{.QuotedExecArgsWithParameter "&" "val" .QuotedAllCol}}); err != nil {
			return nil, err
		}
		vals = append(vals, &val)
	}
	return vals, nil
}

func (dao *{{.GoName}}DAO) SelectBatchByRaw(ctx context.Context, query string, args ...any) ([]*{{.GoName}}, error) {
	rows, err := dao.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return dao.NewBatch(rows)
}

func (dao *{{.GoName}}DAO) SelectBatchByWhere(ctx context.Context, where string, args ...any) ([]*{{.GoName}}, error) {
	s := "SELECT {{.QuotedAllCol}} FROM {{.QuotedTableName}} WHERE " + where
	return dao.SelectBatchByRaw(ctx, s, args...)
}
