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

package routers

import (
	"github.com/astaxie/beego"
	. "github.com/dappledger/ann-explorer/controllers"
)

func init() {

	beego.Router("/", &WebController{}, "get:Index")
	webNs := beego.NewNamespace("/view",
		beego.NSRouter("/blocks", &WebController{}, "get:Blocks"),
		beego.NSRouter("/blocks/hash/*", &WebController{}, "get:Block"),

		beego.NSRouter("/txs/page", &WebController{}, "get:TxsPage"),
		beego.NSRouter("/txs", &WebController{}, "get:TxsLatest"),

		beego.NSRouter("/contracts/page", &WebController{}, "get:ContractPage"),
		beego.NSRouter("/contracts", &WebController{}, "get:ContractsLatest"),
		beego.NSRouter("/contracts/hash/:hash", &WebController{}, "get:Contract"),

		beego.NSRouter("/search/:hash", &WebController{}, "get:Search"),
	)

	apiNs := beego.NewNamespace("/v1",
		beego.NSNamespace("/account"), //beego.NSRouter("/list", &AccountController{}, "get:List"),

		beego.NSNamespace("/txs",
			beego.NSRouter("/query/:fromTo", &TxsController{}, "get:Query")),

		beego.NSNamespace("/info"), //beego.NSRouter("/status", &InfoController{}, "get:Status"),

		beego.NSNamespace("/net"), //beego.NSRouter("/info", &NetController{}, "get:Info"),

		beego.NSNamespace("/block"), //			beego.NSRouter("/blockchain/:maxHeight", &BlockController{}, "get:BlockChain"),

	)

	query := beego.NewNamespace("/query",
		beego.NSRouter("/block", &ApiController{}, "get:QueryBlock"),
		beego.NSRouter("/blocklist", &ApiController{}, "get:QueryBlockList"),
		beego.NSRouter("/tx", &ApiController{}, "get:QueryTx"),
		beego.NSRouter("/txlist", &ApiController{}, "get:QueryTxsList"),
		beego.NSRouter("/txs", &ApiController{}, "get:QueryTxsByBlkHeight"),
	)
	contract := beego.NewNamespace("/contract",
		beego.NSRouter("/meta", &ApiController{}, "post:SaveContractMeta"),
	)
	beego.AddNamespace(webNs)
	beego.AddNamespace(apiNs)
	beego.AddNamespace(query)
	beego.AddNamespace(contract)

	beego.SetStaticPath("/assets", "static")

}
