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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func InitType() *Model {
	return &Model{
		TableName: "user",
		GoName:    "User",
		Fields: []Field{
			{GoType: "uint32"},
			{GoType: "string"},
			{GoType: "*int"},
			{GoType: "bool"},
			{GoType: "[]byte"},
			{GoType: "float32"},
			{GoType: "byte"},
		},
	}
}

func TestWithImports(t *testing.T) {
	model := InitModel()
	WithImports("imports")(model)
	assert.Equal(t, "imports", model.ExtralImport)
}

func TestModel_QuotedExecArgsWithParameter(t *testing.T) {
	dao := InitModel()
	args := dao.QuotedExecArgsWithParameter(dao.QuotedAllCol(), "&", "user.")
	assert.Equal(t, "&user.FirstName, &user.LastName, &user.UserId", args)
}

func TestModel_AddToString(t *testing.T) {
	dao := InitModel()
	args := dao.AddToString(dao.QuotedAllCol())
	assert.Equal(t, "`first_name`,`last_name`,`user_id`", args)
}

func TestModel_QuotedRelationshipIndex(t *testing.T) {
	args := InitModel().QuotedRelationship()["first_name"]
	assert.Equal(t, "FirstName", args)
}

func TestModel_QuotedTableName(t *testing.T) {
	table := InitModel().QuotedTableName()
	assert.Equal(t, "`user`", table)
}

func TestModel_QuotedAllCol(t *testing.T) {
	cols := strings.Join(InitModel().QuotedAllCol(), ",")
	assert.Equal(t, "`first_name`,`last_name`,`user_id`", cols)
}

func TestModel_InsertWithReplaceParameter(t *testing.T) {
	para := InitModel().InsertWithReplaceParameter()
	assert.Equal(t, "?,?,?", para)
}

func TestField_GoType(t *testing.T) {
	m := InitType()
	assert.True(t, m.Fields[0].IsInteger())
	assert.True(t, m.Fields[1].IsString())
	assert.True(t, m.Fields[2].IsPtr())
	assert.True(t, m.Fields[3].IsBool())
	assert.True(t, m.Fields[4].IsSlice())
	assert.True(t, m.Fields[5].IsFloat())
	assert.True(t, m.Fields[6].IsInteger())
}
