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
	"github.com/gotomicro/egen/internal/model"
	"regexp"
	"strings"
)

func ParseModel(contents File, options ...model.Option) []model.Model {
	var models = make([]model.Model, 0, len(contents.TypeNodes))
	for _, v := range contents.TypeNodes {
		models = append(models, parseNode(v))
	}
	for k := range models {
		for _, option := range options {
			option(&models[k])
		}
		models[k].PkgName = contents.PkgName + "."
	}
	return models
}

func parseNode(typeNode TypeNode) model.Model {
	fields := make([]model.Field, 0, len(typeNode.Fields))
	tableName := Convert(typeNode.GoName)
	
	for _, v := range typeNode.Ans {
		if v.Key == "TableName" {
			tableName = v.Value
		}
	}
	for _, v := range typeNode.Fields {
		fields = append(fields, parseField(v))
	}
	
	return model.Model{
		GoName:    typeNode.GoName,
		TableName: tableName,
		Fields:    fields,
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

func Convert(name string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
