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

package ast

import (
	"testing"

	"github.com/gotomicro/egen/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestParseModel(t *testing.T) {
	testCases := []struct {
		src  string
		want []model.Model
	}{
		{
			src: `
package model

type Into struct {
	// @ColName countryside
	CountrySide string
	// @PrimaryKey true
	Suburb  string
}

// @TableName Order
type Order struct{
	// @PrimaryKey true
	// @ColName user_id
	UserId uint32
}
`, want: []model.Model{
				model.Model{
					PkgName:   "model.",
					TableName: "into",
					GoName:    "Into",
					Fields: []model.Field{
						model.Field{ColName: "countryside", IsPrimaryKey: false, GoName: "CountrySide", GoType: "string"},
						model.Field{ColName: "suburb", IsPrimaryKey: true, GoName: "Suburb", GoType: "string"},
					},
				},
				model.Model{
					PkgName:   "model.",
					TableName: "Order",
					GoName:    "Order",
					Fields: []model.Field{
						model.Field{ColName: "user_id", IsPrimaryKey: true, GoName: "UserId", GoType: "uint32"},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.want, ParseModel(LookUp("", tc.src)))
	}
}

func TestConvert(t *testing.T) {
	assert.Equal(t, "first_name", Convert("FirstName"))
	assert.Equal(t, "order", Convert("Order"))
	assert.Equal(t, "into", Convert("into"))
}
