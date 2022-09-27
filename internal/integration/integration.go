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

package integration

/*
//go:generate egen dao -type User -import "github.com/gotomicro/egen/internal/integration"
//go:generate egen dao -src . -import "github.com/gotomicro/egen/internal/integration"
//go:generate egen dao -src ./integration.go -import "github.com/gotomicro/egen/internal/integration"
*/
//go:generate egen dao -src ./integration.go -dst ./test_egen_param/test_tx/user_dao_tx.go -type User -import "github.com/gotomicro/egen/internal/integration"

//go:generate egen dao -dst ./test_egen_param/first -type User -import "github.com/gotomicro/egen/internal/integration"
//go:generate egen dao -src . -dst ./test_egen_param/second -import "github.com/gotomicro/egen/internal/integration"
//go:generate egen_param dao -src . -dst ./test_egen_param/third -type User -import "github.com/gotomicro/egen/internal/integration"
type User struct {
	// @PrimaryKey true
	// @ColName id
	ID       uint64
	Username string
	Password string
	Login    string
}
