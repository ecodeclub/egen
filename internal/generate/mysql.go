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
	"github.com/gotomicro/egen/internal/model"
	"io"
	"text/template"
)

type MySQLGenerator struct {
}

func (*MySQLGenerator) Generate(m *model.Model, writer io.Writer) error {
	var err error
	files := []string{"insert.gohtml", "select.gohtml", "update.gohtml", "delete.gohtml"}
	tMySQL := template.Must(template.ParseGlob("mysql_template/*.gohtml"))
	for _, v := range files {
		t := tMySQL.Lookup(v)
		err = t.Execute(writer, m)
		if err != nil {
			return err
		}
	}
	return nil
}
