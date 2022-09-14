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

package daocmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotomicro/egen/internal/generate"
	"github.com/gotomicro/egen/internal/model"
	"github.com/gotomicro/egen/internal/model/ast"
	"github.com/gotomicro/egen/internal/utils"
)

func execWrite(src, dst, name, path string) error {
	var (
		dstDir = utils.IsDir(dst)
		srcDir = utils.IsDir(src)
	)

	if !dstDir && name == allModel {
		return errors.New("-dst 应该是一个目录，或者使用 -type 指定了单个类型")
	} else if srcDir && !dstDir && name == allModel {
		return errors.New("-src为目录的情况下-dst也应为目录 或者使用 -type指定了单个类型")
	}

	srcFiles := make([]string, 0, 10)
	if srcDir {
		files, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".go") {
				src, err = filepath.Abs(src)
				if err != nil {
					return err
				}
				srcFiles = append(srcFiles, src+"/"+file.Name())
			}
		}
	} else {
		src, err := filepath.Abs(src)
		if err != nil {
			return err
		}
		srcFiles = append(srcFiles, src)
	}

	models := make([]model.Model, 0, len(srcFiles))
	for _, name := range srcFiles {
		models = append(models, ast.ParseModel(ast.LookUp(name, nil), model.WithImports(path))...)
	}

	return WriteToFile(models, dst, name)
}

func WriteToFile(models []model.Model, dst, name string) error {
	var mg generate.MySQLGenerator
	for _, v := range models {
		if name != "" && v.GoName != name {
			continue
		}

		if utils.IsDir(dst) {
			// 可能要对多个文件进行写入 写入完成后直接close
			f, err := os.Create(dst + fmt.Sprintf("/%s_dao.go", ast.Convert(v.TableName)))
			if err != nil {
				f.Close() // 防止内存泄露
				return err
			}
			if err = mg.Generate(v, f); err != nil {
				f.Close() // 防止内存泄露
				return err
			}
			f.Close()
			fmt.Println(f.Name(), "已完成")
		} else {
			f, err := os.Create(dst)
			if err != nil {
				return err
			}

			if err = mg.Generate(v, f); err != nil {
				return err
			}
			f.Close() // 只有单个文件进行写入 直接defer
			fmt.Println(f.Name(), "已完成")
		}
	}
	return nil
}
