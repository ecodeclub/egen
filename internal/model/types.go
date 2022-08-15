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
	IsPrimaryKey bool `mapstructure:"PrimaryKey"`
	GoName       string
	GoType       string
}

func (f *Field) IsInteger() bool {
	switch f.GoType {
	case "int64", "int32", "int16", "int8", "int", "uint64", "uint32", "uint16", "uint8", "uint":
		return true
	case "byte", "rune":
		return true
	}
	return false
}

func (f *Field) IsFloat() bool {
	switch f.GoType {
	case "float32", "float64":
		return true
	}
	return false
}

func (f *Field) IsString() bool {
	return f.GoType == "string"
}

func (f *Field) IsBool() bool {
	return f.GoType == "bool"
}

func (f *Field) IsSlice() bool {
	return f.GoType[0:2] == "[]"
}

func (f *Field) IsPtr() bool {
	return f.GoType[0] == '*'
}

func (m *Model) QuotedTableName() string {
	return "`" + m.TableName + "`"
}

func (m *Model) QuotedExecArgsWithParameter(col []string, flag, owner string) string {
	var str = make([]string, 0, len(m.Fields))
	var strMap = make(map[string]int, len(m.Fields))
	for k, v := range col {
		strMap[v] = k
	}
	for _, v := range m.Fields {
		if _, exist := strMap["`"+v.ColName+"`"]; exist {
			str = append(str, flag+owner+v.GoName)
		}
	}
	return strings.Join(str, ", ")
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

func (m *Model) QuotedAllCol() []string {
	var cols = make([]string, 0, len(m.Fields))
	for _, v := range m.Fields {
		cols = append(cols, "`"+v.ColName+"`")
	}
	return cols
}

func (*Model) AddToString(cols []string) string {
	return strings.Join(cols, ",")
}

func (m *Model) QuotedRelationship() map[string]string {
	relation := make(map[string]string, len(m.Fields))
	for _, v := range m.Fields {
		relation[v.ColName] = v.GoName
	}
	return relation
}
