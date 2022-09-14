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
	"flag"
	"fmt"
	"os"

	daocmd "github.com/gotomicro/egen/cmd/egen/dao"
)

var (
	longHelp  = flag.Bool("help", false, "提供帮助")
	shortHelp = flag.Bool("h", false, "提供帮助")
)

func Execute() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		switch flag.Args()[0] {
		case "dao":
			daocmd.ExecDao(os.Args[2:])
		default:
			usage()
		}
	} else {
		usage()
	}
}

func usage() {
	fmt.Println("提供以下几种命令:")
	daocmd.DaoFlagSet.Usage()
}

func Help() {
	if *shortHelp || *longHelp {
		usage()
	}
}
