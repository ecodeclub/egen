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

package generate

import (
	"bytes"
	"github.com/gotomicro/egen/internal/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMySQLGenerator_Generate(t *testing.T) {
	testCases := []struct {
		name     string
		model    *model.Model
		wantCode string
		wantErr  error
		testdata string
	}{
		{
			name: "user",
			model: &model.Model{
				TableName: "user",
				GoName:    "User",
				Fields: []model.Field{
					{ColName: "login_time", GoName: "LoginTime"},
					{ColName: "first_name", GoName: "FirstName", Order: true},
					{ColName: "last_name", GoName: "LastName", Order: true},
					{ColName: "user_id", GoName: "UserId", IsPrimaryKey: true},
				},
			},
			wantErr:  nil,
			testdata: "./testdata/user.go",
		},
		{
			name: "order",
			model: &model.Model{
				TableName: "order",
				GoName:    "Order",
				Fields: []model.Field{
					{ColName: "order_time", GoName: "OrderTime", Order: true},
					{ColName: "order_id", GoName: "OrderId", Order: true},
					{ColName: "user_id", GoName: "UserId", IsPrimaryKey: true},
				},
			},
			wantErr:  nil,
			testdata: "./testdata/order.go",
		},
	}

	mg := &MySQLGenerator{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := os.ReadFile(testCase.testdata)
			assert.Equal(t, nil, err)
			testCase.wantCode = string(data)
			w := &bytes.Buffer{}
			err = mg.Generate(testCase.model, w)
			assert.Equal(t, testCase.wantErr, err)
			assert.Equal(t, testCase.wantCode, w.String())
		})
	}
}
