package main

import (
	"encoding/json"
	"os/user"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/dappledger/AnnChain-browser/job"
	"github.com/dappledger/AnnChain-browser/repository"
	_ "github.com/dappledger/AnnChain-browser/routers"
	"github.com/dappledger/AnnChain-browser/rpc"
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
