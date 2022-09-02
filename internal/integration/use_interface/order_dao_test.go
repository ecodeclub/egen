//go:build e2e

package use_interface

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type OrderDAOTestSuite struct {
	suite.Suite
	ctx context.Context
	dao OrderDAO
}

func connetDatabase(driver, config string) handler {
	db, err := sql.Open(driver, config)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	
	err = db.Ping()
	for err != nil {
		log.Println("等待数据库开启:", err)
		err = db.Ping()
		time.Sleep(1 * time.Second)
	}
	// tx, err := db.Begin()
	// return tx
	return db
}

func (d *OrderDAOTestSuite) deleteAll() (int64, error) {
	return d.dao.DeleteByWhere(d.ctx, "1=1")
}

func (d *OrderDAOTestSuite) SetupSuite() {
	ownHandler := connetDatabase("mysql", "root:root@tcp(127.0.0.1:13306)/user_infor")
	switch x := ownHandler.(type) {
	case handler:
		d.dao.handler = x
	case handlerTx:
		d.dao.handlerTx = x
		ret := x.(handler)
		d.dao.handler = ret
	}
	d.ctx = context.Background()
}

func (d *OrderDAOTestSuite) TestOrderDAO_Insert() {
	t := d.T()
	orders := []*Order{
		{1, 1, 10},
		{2, 2, 10},
	}
	ret, err := d.dao.Insert(d.ctx, orders...)
	assert.Equal(t, int64(2), ret)
	assert.Equal(t, nil, err)
	
	ret, err = d.dao.Insert(d.ctx, &Order{3, 3, 60})
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
	
	ret, err = d.dao.Insert(d.ctx)
	assert.Equal(t, int64(0), ret)
	assert.Equal(t, nil, err)
	
	d.deleteAll()
	
	// 这种方式来判断 或者直接封装成方法
	if d.dao.handlerTx != nil {
		d.dao.handlerTx.Commit()
	}
	//d.commit()
}

func (d *OrderDAOTestSuite) TestOrderDAO_SelectByWhere() {
	t := d.T()
	
	orderInsert := Order{OrderId: 1, UserId: 1, Price: 10}
	ret, err := d.dao.Insert(d.ctx, &orderInsert)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
	
	order, err := d.dao.SelectByWhere(d.ctx, "user_id=?", 1)
	assert.Equal(t, Order{1, 1, 10}, *order)
	assert.Equal(t, nil, err)
	
	assert.Equal(t, orderInsert, *order)
	
	d.deleteAll()
}

func (d *OrderDAOTestSuite) TestOrderDAO_SelectBatchByRaw() {
	t := d.T()
	
	ordersInsert := []*Order{
		{1, 1, 10},
		{2, 2, 10},
	}
	ret, err := d.dao.Insert(d.ctx, ordersInsert...)
	assert.Equal(t, int64(2), ret)
	assert.Equal(t, nil, err)
	
	orders, err := d.dao.SelectBatchByWhere(d.ctx, "price=?", 10)
	assert.Equal(t, ordersInsert, orders)
	assert.Equal(t, nil, err)
	
	d.deleteAll()
}

func (d *OrderDAOTestSuite) TestOrderDAO_UpdateNoneZeroColByWhere() {
	t := d.T()
	
	orderInsert := Order{OrderId: 1, UserId: 1, Price: 10}
	ret, err := d.dao.Insert(d.ctx, &orderInsert)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
	
	orderInsert.Price = 20
	
	ret, err = d.dao.UpdateNoneZeroColByWhere(d.ctx, &Order{Price: 20}, "user_id=?", 1)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
	
	order, err := d.dao.SelectByWhere(d.ctx, "user_id=?", orderInsert.UserId)
	assert.Equal(t, nil, err)
	assert.Equal(t, orderInsert, *order)
	
	d.deleteAll()
}

func (d *OrderDAOTestSuite) TestOrderDAO_UpdateNonePKColByWhere() {
	t := d.T()
	
	orderInsert := Order{OrderId: 1, UserId: 1, Price: 10}
	ret, err := d.dao.Insert(d.ctx, &orderInsert)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
	
	orderInsert.Price = 20
	
	ret, err = d.dao.UpdateNonePKColByWhere(d.ctx, &Order{OrderId: 1, UserId: 2, Price: 20}, "user_id=?", orderInsert.OrderId)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
	
	order, err := d.dao.SelectByWhere(d.ctx, "user_id=?", orderInsert.UserId)
	assert.Equal(t, nil, err)
	assert.Equal(t, orderInsert, *order)
	
	d.deleteAll()
}

func (d *OrderDAOTestSuite) TestOrderDAO_UpdateSpecificColsByWhere() {
	t := d.T()
	
	orderInsert := Order{OrderId: 1, UserId: 1, Price: 10}
	ret, err := d.dao.Insert(d.ctx, &orderInsert)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
	
	orderInsert.Price = 20
	
	ret, err = d.dao.UpdateSpecificColsByWhere(d.ctx, &Order{Price: 20}, []string{"price"}, "user_id=?", 1)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
	
	order, err := d.dao.SelectByWhere(d.ctx, "user_id=?", orderInsert.UserId)
	assert.Equal(t, nil, err)
	assert.Equal(t, orderInsert, *order)
	
	d.deleteAll()
}

func (d *OrderDAOTestSuite) TestOrderDAO_DeleteByWhere() {
	t := d.T()
	orderInsert := Order{OrderId: 1, UserId: 1, Price: 10}
	ret, err := d.dao.Insert(d.ctx, &orderInsert)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
	
	ret, err = d.dao.DeleteByWhere(d.ctx, "user_id=?", 1)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
}

func Test_All(t *testing.T) {
	suite.Run(t, new(OrderDAOTestSuite))
}
