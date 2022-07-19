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
)

// Model 模型定义
type Model struct {
	TableName string
	GoName    string
	Fields    []Field
}

type Field struct {
	ColName      string
	IsPrimaryKey bool
	OfWhere      bool
	GoName       string
}

func (m *Model) InsertWithReplaceParameter() string {
	var str strings.Builder
	for k := range m.Fields {
		if k != 0 {
			str.WriteByte(',')
		}
		str.WriteByte('?')
	}
	return str.String()
}

func (m *Model) QuotedTableName() string {
	return "`" + m.TableName + "`"
}

func (m *Model) QuotedExecArgsWithAll() string {
	var str strings.Builder
	for k, v := range m.Fields {
		if k != 0 {
			str.WriteByte(',')
		}
		str.WriteString("v." + v.GoName)
	}
	return str.String()
}

func (m *Model) QuotedAllCol() string {
	var str strings.Builder
	for k, v := range m.Fields {
		if k != 0 {
			str.WriteByte(',')
		}
		str.WriteString("`" + v.ColName + "`")
	}
	return str.String()
}
