package routers

import (
	"github.com/astaxie/beego"
	. "github.com/dappledger/AnnChain-browser/controllers"
)

func init() {

	beego.Router("/", &WebController{}, "get:Index")
	webNs := beego.NewNamespace("/view",
		beego.NSRouter("/blocks/latest", &WebController{}, "get:Latest"),
		beego.NSRouter("/blocks/hash/*", &WebController{}, "get:Block"),

		beego.NSRouter("/txs/page", &WebController{}, "get:TxsPage"),
		beego.NSRouter("/txs/latest", &WebController{}, "get:TxsLatest"),

		beego.NSRouter("/contracts/page", &WebController{}, "get:ContractPage"),
		beego.NSRouter("/contracts/latest", &WebController{}, "get:ContractsLatest"),
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
