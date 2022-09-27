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
	TableName    string
	GoName       string
	Fields       []Field
	SelfPkgName  string
	PkgName      string
	ExtralImport string
	Methods      []Method
}

type Method struct {
	FuncName    string
	SqlType     string
	SqlSentence string
	Params      []Parameter
	Results     []string
}

type Parameter struct {
	GoName string
	GoType string // []byte
	Exist  bool
	HasLen bool
}

type Field struct {
	ColName      string
	IsPrimaryKey bool
	GoName       string
	GoType       string
}

type PkgInfor struct {
	PkgName      string
	ExtralImport string
}

type Option func(m *Model)

func WithImports(imports string) Option {
	return func(m *Model) {
		m.ExtralImport = imports
	}
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

func (m Method) IsInteger(t string) bool {
	switch t {
	case "int64", "int32", "int16", "int8", "int", "uint64", "uint32", "uint16", "uint8", "uint":
		return true
	case "byte":
		return true
	}
	return false
}

func (m Method) IsFloat(t string) bool {
	return t == "float32" || t == "float64"
}

func (m Method) IsError(t string) bool {
	return t == "error"
}

func (m Method) IsSlice(t string) bool {
	return strings.Contains(t, "[]")
}

func (md *Model) QuotedGoNameOfSqlParam(Sql string) []string {
	gos := make([]string, 0, 10)
	Sql = strings.ToUpper(strings.Replace(Sql, " ", "", -1))

	if i := strings.Index(Sql, "FROM"); i != -1 {
		Sql = Sql[:i]
	}
	if i := strings.Index(Sql, "`"); i != -1 {
		Sql = Sql[i:]
	}

	params := strings.Split(Sql, ",")
	for _, val := range params {
		for _, v := range md.Fields {
			if strings.ToLower(val) == "`"+v.ColName+"`" {
				gos = append(gos, "&"+"one."+v.GoName)
				break
			}
		}
	}
	return gos
}

func (m Method) QuotedColOfSqlParams() []string {
	params := make([]string, 0, len(m.Params))
	sql := m.SqlSentence
	for index1 := strings.Index(sql, "#"); index1 != -1; index1 = strings.Index(sql, "#") {
		index2 := strings.Index(sql[index1:], " ")
		if index2 != -1 {
			params = append(params, sql[index1+1:index2+index1])
			sql = strings.Replace(sql, sql[index1:index2+index1], "", 1)
		} else {
			params = append(params, sql[index1+1:])
			sql = strings.Replace(sql, sql[index1:], "", 1)
		}
	}
	return params
}

func (md *Model) QuotedColOfSql(m Method) string {
	if m.SqlType == "select" && strings.Contains(m.SqlSentence, "*") && !strings.Contains(m.SqlSentence, "(*)") && !strings.Contains(m.SqlSentence, "( * )") {
		cols := make([]string, 0, len(m.Params))
		for _, v := range md.Fields {
			cols = append(cols, "`"+v.ColName+"`")
		}
		// *替换所有列
		sqlSentence := strings.Replace(m.SqlSentence, "*", strings.Join(cols, ", "), 1)
		return sqlSentence
	}
	return m.SqlSentence
}

func (m Method) QuotedFunc() string {
	var method strings.Builder
	method.WriteString(m.FuncName)
	method.WriteByte('(')
	for k, v := range m.Params {
		if k != 0 {
			method.WriteByte(',')
			method.WriteByte(' ')
		}
		method.WriteString(v.GoName + " " + v.GoType)
	}
	method.WriteByte(')')
	method.WriteByte(' ')

	if len(m.Results) > 0 {
		method.WriteByte('(')
	}
	method.WriteString(strings.Join(m.Results, ", "))
	if len(m.Results) > 0 {
		method.WriteByte(')')
	}

	return method.String()
}

func (m *Model) WrapData(method Method) map[string]interface{} {
	return map[string]interface{}{
		"method": method,
		"model":  m,
	}
}
