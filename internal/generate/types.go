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

package generate

import (
	"io"

	"github.com/gotomicro/egen/internal/model"
)

// Generator 核心接口
// 将生成的代码写入到 writer 里面
// 实际中 writer 可能代表一个文件，也可以是 bytes.Buffer
type Generator interface {
	Generate(m model.Model, writer io.Writer) error
}
