// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	ID       uint64
	Username string
	Password string
	Login    string
}

func (user *User) InitUser(username, password string) {
	user.Username = username
	user.Password = password
	user.Login = time.Now().Format("2006-01-02 15:04:05")
}

// insert
func (user *User) Insert(DB *sql.DB) error {
	sqlSentence := `INSERT INTO user_account(username,login,password) VALUE(?,?,?)`
	_, err := DB.Exec(sqlSentence, user.Username, user.Login, user.Password)
	return err
}

// select
func (user *User) Select(DB *sql.DB) error {
	sqlSentence := `SELECT * FROM user_account WHERE id = ?`
	row := DB.QueryRow(sqlSentence, user.ID)
	return row.Scan(&user.ID, &user.Password, &user.Login, &user.Username)
}

// delete
func (user *User) Delete(DB *sql.DB) error {
	sqlSentence := `DELETE FROM user_account WHERE id = ?`
	_, err := DB.Exec(sqlSentence, user.ID)
	return err
}

// update
func (user *User) Update(DB *sql.DB) error {
	sqlSentence := `UPDATE user_account SET username=?,login=?,password=? WHERE id=?`
	_, err := DB.Exec(sqlSentence, user.Username, user.Login, user.Password, user.ID)
	return err
}
