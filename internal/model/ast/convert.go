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
	"regexp"
	"strings"

	"github.com/gotomicro/egen/internal/model"
)

func ParseModel(contents File, options ...model.Option) []model.Model {
	// 一个结构体 + 一个接口
	var models = make([]model.Model, 0, len(contents.TypeNodes))
	for _, value := range contents.TypeNodes {
		one := parseNode(value)
		one.PkgName = contents.PkgName + "."

		for _, option := range options {
			option(&one)
		}

		if len(one.Methods) <= 0 {
			models = append(models, one)
			continue
		}

		for i := 0; i < len(models); i++ {
			if models[i].GoName+"DAO" == one.GoName {
				models[i].Methods = make([]model.Method, len(one.Methods))
				copy(models[i].Methods, one.Methods)
			}
		}
	}

	return models
}

func parseNode(typeNode TypeNode) model.Model {
	var methods []model.Method
	var fields []model.Field
	tableName := Convert(typeNode.GoName)

	for _, v := range typeNode.Ans {
		if v.Key == "TableName" {
			tableName = v.Value
		}
	}

	if len(typeNode.Fields) > 0 {
		fields = make([]model.Field, 0, len(typeNode.Fields))
		for _, v := range typeNode.Fields {
			fields = append(fields, parseField(v))
		}
	}

	if len(typeNode.Methods) > 0 {
		methods = make([]model.Method, 0, len(typeNode.Methods))
		for _, v := range typeNode.Methods {
			methods = append(methods, ParseMethods(v))
		}
		return model.Model{
			GoName:  typeNode.GoName,
			Methods: methods,
		}
	}

	return model.Model{
		GoName:    typeNode.GoName,
		TableName: tableName,
		Fields:    fields,
		Methods:   methods,
	}
}

func parseField(field Field) model.Field {
	colName := Convert(field.GoName)
	isPrimaryKey := false
	for _, v := range field.Ans {
		switch v.Key {
		case "ColName":
			colName = v.Value
		case "PrimaryKey":
			if v.Value != "false" {
				isPrimaryKey = true
			}
		}
	}
	return model.Field{
		ColName:      colName,
		IsPrimaryKey: isPrimaryKey,
		GoName:       field.GoName,
		GoType:       field.GoType,
	}
}

func ParseMethods(method Method) model.Method {
	if len(method.Ans) > 0 {
		return model.Method{
			FuncName:    method.FuncName,
			SqlType:     method.Ans[0].Key,
			SqlSentence: method.Ans[0].Value,
			Params:      method.Params,
			Results:     method.Results,
		}
	}
	return model.Method{
		FuncName: method.FuncName,
		Params:   method.Params,
		Results:  method.Results,
	}
}

func Convert(name string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
