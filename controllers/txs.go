package controllers

import (
	"github.com/astaxie/beego"
	"github.com/dappledger/ann-explorer/repository"
)

type TxsController struct {
	beego.Controller
}

func (tc *TxsController) Query() {
	fromTo := tc.Ctx.Input.Param(":fromTo")
	defer tc.ServeJSON()
	data, _ := repository.TxsQuery(fromTo)
	tc.Data["json"] = &Result{
		Success: true,
		Data:    data,
	}
}
