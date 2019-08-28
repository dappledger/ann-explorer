// Copyright © 2017 ZhongAn Technology
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"fmt"
	"github.com/astaxie/beego"
)

var HTTP_ADDR string

func Init() {

	HTTP_ADDR = fmt.Sprintf("http://%s", beego.AppConfig.String("api_addr"))
	if HTTP_ADDR == "" {
		beego.Error("HTTP_ADDR is nil")
	}
}

type Result struct {
	Code    uint32      `json:"code"`
	Payload interface{} `json:"payload"`
}
