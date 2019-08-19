package rpc

/*
import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/astaxie/beego"
)

var RPC_ADDR string
var CLIENT *rpcclient.ClientJSONRPC

func init() {
	nodes := os.Getenv("NODES")
	if nodes == "" {
		RPC_ADDR = beego.AppConfig.String("rpc_addr")
	} else {
		RPC_ADDR = "tcp://" + strings.Split(nodes, ",")[0]
	}

	fmt.Println("RPC_ADDR = ", RPC_ADDR)
	if RPC_ADDR == "" {
		fmt.Println("RPC_ADDR is nil")
		os.Exit(0)
	}
	CLIENT = rpcclient.NewClientJSONRPC(RPC_ADDR)
	wire.RegisterInterface(
		struct{ anntypes.TMResult }{},
		wire.ConcreteType{&anntypes.ResultBlockchainInfo{}, anntypes.ResultTypeBlockchainInfo},
		wire.ConcreteType{&anntypes.ResultBlock{}, anntypes.ResultTypeBlock},
	)
}

type Result struct {
	Code    uint32      `json:"code"`
	Payload interface{} `json:"payload"`
}

func TCall(method string, params []interface{}) *Result {
	tmResult := new(anntypes.TMResult)
	_, err := Call(method, params, tmResult)
	if err != nil {
		return &Result{
			Code:    500,
			Payload: err.Error(),
		}
	}

	return &Result{
		Code:    http.StatusOK,
		Payload: *tmResult,
	}
}

func Call(method string, params []interface{}, result interface{}) (interface{}, error) {
	return CLIENT.Call(method, params, result)
}
*/
