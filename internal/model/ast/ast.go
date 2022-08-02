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

// 在这里提供基于 AST 的实现。

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type fileVisitor struct {
	ans   Annotations
	types []*typeVisitor
}

type typeVisitor struct {
	ans    Annotations
	fields []Field
	GoName string
}

type TypeNode struct {
	Annotations
	Fields []Field
	GoName string
}

type Field struct {
	GoName string
	GoType string
	Annotations
}

type File struct {
	Annotations
	TypeNodes []TypeNode
}

func (f *fileVisitor) Get() File {
	types := make([]TypeNode, 0, len(f.types))
	for _, t := range f.types {
		types = append(types, t.Get())
	}
	return File{
		Annotations: f.ans,
		TypeNodes:   types,
	}
}

func (f *fileVisitor) Visit(node ast.Node) ast.Visitor {
	file, ok := node.(*ast.File)
	if ok {
		f.ans = newAnnotations(file.Doc)
		return f
	}
	typ, ok := node.(*ast.TypeSpec)
	if ok {
		res := &typeVisitor{
			ans:    newAnnotations(typ.Doc),
			fields: make([]Field, 0),
			GoName: typ.Name.Name,
		}
		f.types = append(f.types, res)
		return res
	}
	return f
}

func (t *typeVisitor) Get() TypeNode {
	return TypeNode{
		Annotations: t.ans,
		Fields:      t.fields,
		GoName:      t.GoName,
	}
}

func (t *typeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	fd, ok := node.(*ast.Field)
	if ok {
		t.fields = append(t.fields, Field{
			Annotations: newAnnotations(fd.Doc),
			GoName:      fd.Names[0].Name,
			GoType:      fmt.Sprintf("%v", fd.Type),
		})
		return nil
	}
	return t
}

func LookUp(path string, src any) File {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	fv := &fileVisitor{}
	ast.Walk(fv, f)
	return fv.Get()
}
