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

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func InitModel() *Model {
	return &Model{
		TableName: "user",
		GoName:    "User",
		Fields: []Field{
			{ColName: "first_name", GoName: "FirstName"},
			{ColName: "last_name", GoName: "LastName"},
			{ColName: "user_id", GoName: "UserId"},
		},
	}
}

func TestModel_QuotedExecArgsWithParameter(t *testing.T) {
	args := InitModel().QuotedExecArgsWithParameter("&", "user", "`first_name`,`last_name`")
	assert.Equal(t, "&user.FirstName, &user.LastName", args)
}

func TestModel_QuotedTableName(t *testing.T) {
	table := InitModel().QuotedTableName()
	assert.Equal(t, "`user`", table)
}

func TestModel_QuotedAllCol(t *testing.T) {
	cols := InitModel().QuotedAllCol()
	assert.Equal(t, "`first_name`,`last_name`,`user_id`", cols)
}

func TestModel_QuotedExecArgsWithAll(t *testing.T) {
	args := InitModel().QuotedExecArgsWithAll()
	assert.Equal(t, "v.FirstName, v.LastName, v.UserId", args)
}

func TestModel_InsertWithReplaceParameter(t *testing.T) {
	para := InitModel().InsertWithReplaceParameter()
	assert.Equal(t, "?,?,?", para)
}
