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

package main

import (
	"encoding/json"
	"os/user"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/dappledger/ann-explorer/job"
	"github.com/dappledger/ann-explorer/repository"
	_ "github.com/dappledger/ann-explorer/routers"
	"github.com/dappledger/ann-explorer/rpc"
)

type logConfig struct {
	Filename string `json:"filename"`
}

func main() {
	rpc.Init()
	repository.Init()
	go job.SyncTimingTask()
	beego.SetLogFuncCall(true)
	user, _ := user.Current()
	logPath, _ := json.Marshal(logConfig{Filename: user.HomeDir + "/browser.log"})
	logs.SetLogger(logs.AdapterMultiFile, string(logPath))
	beego.Run()
}
