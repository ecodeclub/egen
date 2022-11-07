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

package daocmd

import (
	"testing"

	"github.com/gotomicro/egen/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUpdateByParams(t *testing.T) {
	var m = model.Model{
		PkgName:   "model.",
		TableName: "into",
		GoName:    "Into",
		Fields: []model.Field{
			{ColName: "countryside", IsPrimaryKey: false, GoName: "CountrySide", GoType: "string"},
			{ColName: "suburb", IsPrimaryKey: true, GoName: "Suburb", GoType: "string"},
		},
		Methods: []model.Method{
			{
				FuncName:    "Hello",
				SqlType:     "select",
				SqlSentence: "select * from `user` where `name` in #names and `status`=#status",
				Params: []model.Parameter{
					{GoName: "ctx", GoType: "context.Context", Exist: false, HasLen: true},
					{GoName: "names", GoType: "[5]string", Exist: false, HasLen: true},
					{GoName: "status", GoType: "bool", Exist: false, HasLen: false},
				},
				Results: []string{"int64", "error"},
			},
		},
	}
	var data = m
	data.SelfPkgName = "code"
	updateByParams(&m, "./dao", "./dao/code")
	assert.Equal(t, m, data)
}
