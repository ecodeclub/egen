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
package cmd

import (
	"github.com/gotomicro/egen/internal/generate"
	"github.com/gotomicro/egen/internal/model/ast"
	"os"
	"path/filepath"
)

func Cmd() {
	mg := generate.MySQLGenerator{}
	path, _ := filepath.Abs("./cmd/egen/data/data.go")
	for _, v := range ast.ParseModel(ast.LookUp(path, nil)) {
		mg.Generate(&v, os.Stdout)
	}
}
