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
				srcFiles = append(srcFiles, filepath.Join(src, file.Name()))
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
		for _, val := range ast.ParseModel(ast.LookUp(name, nil), model.WithImports(path)) {
			updateByParams(&val, src, dst)
			models = append(models, val)
		}
	}

	return WriteToFile(models, dst, name)
}

func WriteToFile(models []model.Model, dst, name string) error {
	var mg generate.MySQLGenerator
	for _, v := range models {
		if name != "" && v.GoName != name {
			continue
		}
		var f *os.File
		var err error

		if utils.IsDir(dst) {
			// 可能要对多个文件进行写入 写入完成后直接close
			f, err = os.Create(filepath.Join(dst, fmt.Sprintf("%s_dao.go", ast.Convert(v.TableName))))
		} else {
			f, err = os.Create(dst)
		}
		if err != nil {
			return err
		}

		if err = mg.Generate(v, f); err != nil {
			return err
		}

		f.Close()
		fmt.Println(f.Name(), "已完成")
	}
	return nil
}

func updateByParams(v *model.Model, src, dst string) {
	if filepath.Dir(src) != filepath.Dir(dst) {
		v.SelfPkgName = filepath.Base(dst)
	} else {
		v.SelfPkgName = v.PkgName[:len(v.PkgName)-1]
		v.PkgName = ""
	}

	for k := range v.Methods {
		for i := range v.Methods[k].Params {
			v.Methods[k].Params[i].GoType = strings.ReplaceAll(v.Methods[k].Params[i].GoType, v.GoName, v.PkgName+v.GoName)
		}
		for i := range v.Methods[k].Results {
			v.Methods[k].Results[i] = strings.ReplaceAll(v.Methods[k].Results[i], v.GoName, v.PkgName+v.GoName)
		}

		sqlSentence := v.Methods[k].SqlSentence
		// 保证params是有序的，模板中遍历append参数是按照sql中参数顺序来的
		params := make([]model.Parameter, 0, len(v.Methods[k].Params))
		for indexSpecial := strings.Index(sqlSentence, "#"); indexSpecial != -1; indexSpecial = strings.Index(sqlSentence, "#") {
			name := ""
			indexSpace := strings.Index(sqlSentence[indexSpecial:], " ")
			if indexSpace == -1 {
				name = sqlSentence[indexSpecial+1:]
			} else {
				name = sqlSentence[indexSpecial+1 : indexSpecial+indexSpace]
			}
			sqlSentence = strings.Replace(sqlSentence, "#", "", 1)
			for _, val := range v.Methods[k].Params {
				if name == val.GoName {
					// 为false时，不需要再进行append了
					val.Exist = true
					// 判断是[]byte之外的切片类型
					val.HasLen = strings.Contains(val.GoType, "[") &&
						strings.Contains(val.GoType, "]") &&
						!strings.Contains(val.GoType, "byte")
					// 根据sql中的参数进行append，保证顺序正确
					params = append(params, val)
				}
			}
		}
		// 把除第一个参数之外的替换成出现在sql里面的参数
		v.Methods[k].Params = append(v.Methods[k].Params[0:1], params...)
	}
}
