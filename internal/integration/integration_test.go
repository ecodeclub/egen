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

//go:build e2e

package integration

import (
	"database/sql"
	"log"
	"testing"
	"time"
)

func TestSqlRun(t *testing.T) {
	var user User
	user.InitUser("admin", "12345")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:13306)/user")
	if err != nil {
		t.Fatal("连接失败:", err)
	}

	err = db.Ping()
	for err != nil {
		log.Println("等待数据库开启:", err)
		err = db.Ping()
		time.Sleep(1000)
	}

	if err = user.Insert(db); err != nil {
		t.Fatal("fail to insert:", err)
	}

	if err = user.Insert(db); err != nil {
		t.Fatal("fail to insert:", err)
	}

	user.ID = 1
	if err = user.Select(db); err != nil {
		t.Fatal("fail to select:", err)
	}

	user.Login = time.Now().Format("2006-01-02 15:04:05")
	if err = user.Update(db); err != nil {
		t.Fatal("fail to update:", err)
	}

	if err = user.Delete(db); err != nil {
		t.Fatal("fail to delete:", err)
	}
}
