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
	"fmt"
	"go/ast"

	"github.com/gotomicro/egen/internal/model"
)

func getType(value ast.Expr) string {
	typeName := ""
	switch x := value.(type) {
	case *ast.ArrayType:
		eType := ""
		switch xc := x.Elt.(type) {
		case *ast.Ident:
			eType = xc.Name
		case *ast.StarExpr:
			switch xcc := xc.X.(type) {
			case *ast.Ident:
				eType = "*" + xcc.Name
			case *ast.SelectorExpr:
				eType = fmt.Sprintf("*%v.%v", xcc.X, xcc.Sel.Name)
			}
		case *ast.SelectorExpr:
			eType = fmt.Sprintf("%v.%v", xc.X, xc.Sel.Name)
		}
		if x.Len == nil {
			typeName = fmt.Sprintf("[]%v", eType)
		} else {
			t, ok := x.Len.(*ast.BasicLit)
			if ok {
				typeName = fmt.Sprintf("[%v]%v", t.Value, eType)
			}
		}
	case *ast.Ident:
		typeName = x.Name
	case *ast.SelectorExpr:
		typeName = fmt.Sprintf("%v.%v", x.X, x.Sel.Name)
	}
	return typeName
}

func getParam(node *ast.FuncType) []model.Parameter {
	var methods = make([]model.Parameter, 0, len(node.Params.List))
	for _, value := range node.Params.List {
		var method = model.Parameter{
			GoName: value.Names[0].String(),
			GoType: getType(value.Type),
			Exist:  false,
			HasLen: false,
		}
		methods = append(methods, method)
	}
	if methods[0].GoType != "context.Context" {
		panic("自定义方法的首个参数必须要context.Context 类型")
	}
	return methods
}

func getResult(node *ast.FuncType) []string {
	var method = make([]string, 0, len(node.Params.List))

	for _, value := range node.Results.List {
		method = append(method, getType(value.Type))
	}

	return method
}
