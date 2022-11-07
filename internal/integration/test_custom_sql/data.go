package genDAO

import "context"

// User @TableName user_dao
//go:generate egen dao -src ./data.go -type User -dst user_dao.go
type User struct {
	// @PrimaryKey true
	// @ColName id
	ID       uint64
	Username string
	Password string
	Login    string
	Status   bool
	Money    float64
}

type UserDAO interface {
	// @select SELECT * FROM `user_dao` WHERE `id` IN #ids AND `status` = #status
	FindExist(ctx context.Context, status bool, ids []uint32) ([]*User, error)
	
	// @update UPDATE `user_dao` SET `username`=#name WHERE `status`=#status
	UpdateOne(ctx context.Context, status bool, name string) (int64, error)
	
	// @delete DELETE FROM `user_dao` WHERE `status`=#status AND `id` NOT IN #ids
	DeleteOne(ctx context.Context, status bool, ids []uint32) (int64, error)
	
	// @select SELECT COUNT(*) FROM `user_dao`
	GetCount(ctx context.Context) (int, error)
	
	// @select SELECT SUM(`money`) FROM `user_dao`
	GetTotalMoney(ctx context.Context) (float64, error)
	
	// @select SELECT `status` FROM `user_dao`
	GetAllStatus(ctx context.Context) ([]bool, error)
}
