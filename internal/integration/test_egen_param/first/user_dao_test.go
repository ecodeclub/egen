//go:build e2e

package first

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gotomicro/egen/internal/integration"
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
	users := []*integration.User{
		{1, "first", "123", "8.21"},
		{2, "second", "123", "8.22"},
	}
	ret, err := d.dao.Insert(d.ctx, users...)
	t := d.T()

	assert.Equal(t, int64(2), ret)
	assert.Equal(t, nil, err)

	user := integration.User{ID: 3, Username: "third", Password: "1234", Login: "8.23"}
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
	ret, err := d.dao.Insert(d.ctx, &integration.User{ID: 1, Username: "first", Password: "123", Login: "8.21"})
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	user, err := d.dao.SelectByWhere(d.ctx, "id=?", 1)
	assert.Equal(t, integration.User{ID: 1, Username: "first", Password: "123", Login: "8.21"}, *user)
	assert.Equal(t, nil, err)

	d.deleteAll()
}

func (d *UserDAOTestSuite) TestUserDAO_SelectBatchByRaw() {
	t := d.T()

	ret, err := d.dao.Insert(d.ctx, []*integration.User{
		{1, "first", "123", "8.21"},
		{2, "second", "123", "8.22"},
	}...)
	assert.Equal(t, int64(2), ret)
	assert.Equal(t, nil, err)

	users, err := d.dao.SelectBatchByWhere(d.ctx, "password=?", "123")
	assert.Equal(t, []*integration.User{
		{1, "first", "123", "8.21"},
		{2, "second", "123", "8.22"},
	}, users)
	assert.Equal(t, nil, err)

	d.deleteAll()
}

func (d *UserDAOTestSuite) TestUserDAO_UpdateNoneZeroColByWhere() {
	t := d.T()
	user := integration.User{ID: 1, Username: "first", Password: "123", Login: "8.21"}
	ret, err := d.dao.Insert(d.ctx, &user)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	ret, err = d.dao.UpdateNoneZeroColByWhere(d.ctx, &integration.User{Username: "second"}, "id=?", 1)
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

	user := integration.User{ID: 1, Username: "first", Password: "123", Login: "8.21"}
	ret, err := d.dao.Insert(d.ctx, &user)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	user = integration.User{ID: 2, Username: "second", Password: "123", Login: "8.21"}
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

	user := integration.User{ID: 1, Username: "first", Password: "123", Login: "8.21"}
	ret, err := d.dao.Insert(d.ctx, &user)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	ret, err = d.dao.UpdateSpecificColsByWhere(d.ctx, &integration.User{Username: "second"}, []string{"username"}, "id=?", 1)
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
	ret, err := d.dao.Insert(d.ctx, &integration.User{ID: 1, Username: "first", Password: "123", Login: "8.21"})
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)

	ret, err = d.dao.DeleteByWhere(d.ctx, "id=?", 1)
	assert.Equal(t, int64(1), ret)
	assert.Equal(t, nil, err)
}

func Test_All(t *testing.T) {
	suite.Run(t, new(UserDAOTestSuite))
}
