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

	tMySQL := template.New("mysql")
	tMySQL = template.Must(tMySQL.Parse(`
type {{.GoName}}DAO struct{
	DB *sql.DB
}

func (dao *{{.GoName}}DAO) Insert(vals ...*{{.GoName}})(int64,error) {
	var agrs = make([]interface,len(vals)*({{len .Fields}}))
	var str = ""
	for k,v := range vals {
		if k != 0 {
			str += ","
		}
		str += "({{.InsertWithReplaceParameter}})"
		args = append(args,{{.QuotedExecArgsWithAll}})
	}
	sqlSen := "INSERT INTO {{.QuotedTableName}}({{.QuotedAllCol}}) VALUES" + str
	res,err := dao.DB.Exec(sqlSen,args)
	if err != nil {
		return 0,err
	}
	return res.RowsAffected()
}`))
	return tMySQL.Execute(writer, m)
}
