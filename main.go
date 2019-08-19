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
