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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileVisitor_Get(t *testing.T) {
	testCases := []struct {
		src  string
		want File
	}{
		{
			src: `
package model

type (
	Into struct {
		// @ColName countryside
		CountrySide string
		// @PrimaryKey true
		Suburb  string
	}
	// @TableName Order
	Order struct{
	// @PrimaryKey true
	// @ColName user_id
		UserId uint32
	}
)
`,
			want: File{
				TypeNodes: []TypeNode{
					{
						GoName: "Into",
						Fields: []Field{
							{
								Annotations: Annotations{
									Ans: []Annotation{
										{
											Key:   "ColName",
											Value: "countryside",
										},
									},
								},
								GoType: "string",
								GoName: "CountrySide",
							},
							{
								Annotations: Annotations{
									Ans: []Annotation{
										{
											Key:   "PrimaryKey",
											Value: "true",
										},
									},
								},
								GoType: "string",
								GoName: "Suburb",
							},
						},
					},
					{
						Annotations: Annotations{
							Ans: []Annotation{
								{
									Key:   "TableName",
									Value: "Order",
								},
							},
						},
						GoName: "Order",
						Fields: []Field{
							{
								Annotations: Annotations{
									Ans: []Annotation{
										{
											Key:   "PrimaryKey",
											Value: "true",
										},
										{
											Key:   "ColName",
											Value: "user_id",
										},
									},
								},
								GoName: "UserId",
								GoType: "uint32",
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		file := LookUp("", tc.src)
		assertAnnotations(t, tc.want.Annotations, file.Annotations)
		if len(tc.want.TypeNodes) != len(file.TypeNodes) {
			t.Fatal()
		}
		for i, typ := range file.TypeNodes {
			wantType := tc.want.TypeNodes[i]
			assertAnnotations(t, wantType.Annotations, typ.Annotations)
			if len(wantType.Fields) != len(typ.Fields) {
				t.Fatal()
			}
			for j, fd := range typ.Fields {
				wantFd := wantType.Fields[j]
				assertAnnotations(t, wantFd.Annotations, fd.Annotations)
			}
		}
	}
}

func assertAnnotations(t *testing.T, wantAns Annotations, dst Annotations) {
	want := wantAns.Ans
	if len(want) != len(dst.Ans) {
		
		t.Fatal()
	}
	for i, an := range want {
		val := dst.Ans[i]
		assert.Equal(t, an.Value, val.Value)
	}
}
