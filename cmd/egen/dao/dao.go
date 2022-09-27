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
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotomicro/egen/internal/utils"
)

var (
	src, dst, dataModel, path string
	DaoFlagSet                = initDaoFlagSet()
	tips                      = `
  -src file/dir -dst file/dir -type string  -> 将src中的指定的type生成到./dst中
  -src file/dir -dst dir                    -> 将src中的所有type,生成对应的./dst/type_dao.go
  -src file/dir -type string                -> 若src中存在该type,则生成./dao/type_dao.go
  -src file                                 -> 将src中的所有type,生成对应的./dao/type_dao.go
  -src dir                                  -> 扫描src下的所有go文件,若存在type,则生成type_dao.go.不会递归往下查找
  -dst file/dir -type string                -> 若当前目录下存在该type,则生成./dst/type_dao.go
  -type string                              -> 在当前目录下的go文件,若存在type,则生成./dao/type_dao.go
  
`
)

const (
	allModel    = ""
	defaultDst  = "./dao"
	defaultSrc  = "."
	defaultPath = ""
)

func ExecDao(args []string) {
	if len(args) < 1 {
		log.Println("将扫描当前目录下的所有go文件,并生成对应的type_dao.go")
	}
	if err := DaoFlagSet.Parse(args); err != nil {
		log.Println(err)
	}
	if err := initDao(src, dst, dataModel, path); err != nil {
		log.Println(err)
	}
}

func initDao(src, dst, dataModel, path string) error {
	if dst == defaultDst || !utils.IsExist(dst) && !strings.HasSuffix(dst, ".go") {
		if err := os.MkdirAll(dst, 0666); err != nil {
			return err
		}
	} else if strings.HasSuffix(dst, ".go") && !utils.IsExist(dst) {
		dir, _ := filepath.Split(dst)
		if err := os.MkdirAll(dir, 0666); err != nil {
			return err
		}
	}
	return execWrite(src, dst, dataModel, path)
}

func initDaoFlagSet() *flag.FlagSet {
	daoFlagSet := flag.NewFlagSet("dao", flag.ExitOnError)
	daoFlagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of dao:\n")
		daoFlagSet.PrintDefaults()
		fmt.Print(tips)
	}

	daoFlagSet.SetOutput(os.Stdout)
	daoFlagSet.StringVar(&dst, "dst", defaultDst, "生成的代码写入的文件或目录")
	daoFlagSet.StringVar(&src, "src", defaultSrc, "读取结构体的文件或目录")
	daoFlagSet.StringVar(&dataModel, "type", allModel, "结构体名称")
	daoFlagSet.StringVar(&path, "import", defaultPath, "import时的路径")

	return daoFlagSet
}
