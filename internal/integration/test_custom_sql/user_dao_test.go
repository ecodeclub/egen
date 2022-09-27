//go:build e2e

package genDAO

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserDAOTestSuite struct {
	suite.Suite
	ctx context.Context
	dao *UserGenDAO
}

func connectDatabase(driver, config string) *sql.DB {
	common, err := sql.Open(driver, config)
	if err != nil {
		log.Fatal("连接失败:", err)
	}

	err = common.Ping()
	for err != nil {
		log.Println("等待数据库开启:", err)
		err = common.Ping()
		time.Sleep(1 * time.Second)
	}

	return common
}

func (d *UserDAOTestSuite) deleteAll() (int64, error) {
	return d.dao.DeleteByWhere(d.ctx, "1=1")
}

func (d *UserDAOTestSuite) SetupSuite() {
	d.ctx = context.Background()
	d.dao, _ = NewUserGenDAO(connectDatabase("mysql", "root:root@tcp(127.0.0.1:13306)/user_infor"))
}

func (d *UserDAOTestSuite) TearDownSuite() {
	d.dao.session.(*sql.DB).Close()
}

func (d *UserDAOTestSuite) TestUserDAO_Insert() {
	users := []*User{
		{1, "first", "123", "8.21", true, 0},
		{2, "second", "123", "8.22", true, 0},
	}
	ret, err := d.dao.Insert(d.ctx, users...)
	t := d.T()

	assert.Equal(t, int64(2), ret)
	assert.Equal(t, nil, err)

	user := User{ID: 3, Username: "third", Password: "1234", Login: "8.23"}
	ret, err = d.dao.Insert(d.ctx, &user)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	ret, err = d.dao.Insert(d.ctx)
	assert.Equal(t, int64(0), ret)
	assert.Equal(t, nil, err)

	d.deleteAll()
}

func (d *UserDAOTestSuite) TestUserDAO_SelectByWhere() {
	t := d.T()
	ret, err := d.dao.Insert(d.ctx, &User{ID: 1, Username: "first", Password: "123", Login: "8.21"})
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	user, err := d.dao.SelectByWhere(d.ctx, "id=?", 1)
	assert.Equal(t, User{ID: 1, Username: "first", Password: "123", Login: "8.21"}, *user)
	assert.Equal(t, nil, err)

	d.deleteAll()
}

func (d *UserDAOTestSuite) TestUserDAO_SelectBatchByRaw() {
	t := d.T()

	ret, err := d.dao.Insert(d.ctx, []*User{
		{1, "first", "123", "8.21", true, 0},
		{2, "second", "123", "8.22", true, 0},
	}...)
	assert.Equal(t, int64(2), ret)
	assert.Equal(t, nil, err)

	users, err := d.dao.SelectBatchByWhere(d.ctx, "password=?", "123")
	assert.Equal(t, []*User{
		{1, "first", "123", "8.21", true, 0},
		{2, "second", "123", "8.22", true, 0},
	}, users)
	assert.Equal(t, nil, err)

	d.deleteAll()
}

func (d *UserDAOTestSuite) TestUserDAO_UpdateNoneZeroColByWhere() {
	t := d.T()
	user := User{ID: 1, Username: "first", Password: "123", Login: "8.21", Money: 0}
	ret, err := d.dao.Insert(d.ctx, &user)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	ret, err = d.dao.UpdateNoneZeroColByWhere(d.ctx, &User{Username: "second"}, "id=?", 1)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	user.Username = "second"

	userPt, err := d.dao.SelectByWhere(d.ctx, "id=?", user.ID)
	assert.Equal(t, user, *userPt)
	assert.Equal(t, nil, err)

	d.deleteAll()
}

func (d *UserDAOTestSuite) TestUserDAO_UpdateNonePKColByWhere() {
	t := d.T()

	user := User{ID: 1, Username: "first", Password: "123", Login: "8.21", Money: 0}
	ret, err := d.dao.Insert(d.ctx, &user)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	user = User{ID: 2, Username: "second", Password: "123", Login: "8.21"}
	ret, err = d.dao.UpdateNonePKColByWhere(d.ctx, &user, "id=?", 1)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	user.ID = 1

	userPt, err := d.dao.SelectByWhere(d.ctx, "id=?", 1)
	assert.Equal(t, user, *userPt)
	assert.Equal(t, nil, err)

	d.deleteAll()
}

func (d *UserDAOTestSuite) TestUserDAO_UpdateSpecificColsByWhere() {
	t := d.T()

	user := User{ID: 1, Username: "first", Password: "123", Login: "8.21"}
	ret, err := d.dao.Insert(d.ctx, &user)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	ret, err = d.dao.UpdateSpecificColsByWhere(d.ctx, &User{Username: "second"}, []string{"username"}, "id=?", 1)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	user.Username = "second"

	userPt, err := d.dao.SelectByWhere(d.ctx, "id=?", user.ID)
	assert.Equal(t, user, *userPt)
	assert.Equal(t, nil, err)

	d.deleteAll()
}

func (d *UserDAOTestSuite) TestUserDAO_DeleteByWhere() {
	t := d.T()
	ret, err := d.dao.Insert(d.ctx, &User{ID: 1, Username: "first", Password: "123", Login: "8.21"})
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	ret, err = d.dao.DeleteByWhere(d.ctx, "id=?", 1)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
}

func (d *UserDAOTestSuite) TestUserDAO_Method() {
	t := d.T()
	ret, err := d.dao.Insert(d.ctx, []*User{
		{1, "first", "123", "8.21", true, 15.6},
		{2, "second", "123", "8.22", false, 17.1},
		{3, "second", "123", "8.22", true, 12.9},
	}...)
	assert.Equal(t, int64(3), ret)
	assert.Equal(t, nil, err)

	nums, err := d.dao.GetCount(d.ctx)
	assert.Equal(t, int(3), nums)
	assert.Equal(t, nil, err)

	affNum, err := d.dao.UpdateOne(d.ctx, "third", false)
	assert.Equal(t, int64(1), affNum)
	assert.Equal(t, nil, err)

	user := User{2, "third", "123", "8.22", false, 17.1}
	rets, err := d.dao.FindExist(d.ctx, []uint32{1, 2}, false)
	assert.Equal(t, []*User{
		&user,
	}, rets)
	assert.Equal(t, nil, err)

	affNum, err = d.dao.DeleteOne(d.ctx, true, []uint32{1, 3})
	assert.Equal(t, int64(0), affNum)
	assert.Equal(t, nil, err)

	money, err := d.dao.GetTotalMoney(d.ctx)
	money, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", money), 64)
	assert.Equal(t, 45.6, money)
	assert.Equal(t, nil, err)

	status, err := d.dao.GetAllStatus(d.ctx)
	assert.Equal(t, []bool{true, false, true}, status)
	assert.Equal(t, nil, err)

	d.deleteAll()
}

func Test_All(t *testing.T) {
	suite.Run(t, new(UserDAOTestSuite))
}
