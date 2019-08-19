package controllers

import (
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/dappledger/AnnChain/eth/common"
	"github.com/dappledger/ann-explorer/repository"
)

const (
	DisplayNum   = 25
	BlockHashLen = 40
	TxHashLen    = 64
)

type WebController struct {
	beego.Controller
}

type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func (wc *WebController) Index() {
	wc.Layout = "layout.html"
	wc.TplName = "web.tpl"
}

func (wc *WebController) Latest() {
	defer wc.ServeJSON()
	data, err := repository.LatestBlocks(DisplayNum)
	if err != nil {
		beego.Error(err)
	}
	wc.Data["json"] = &Result{
		Success: true,
		Data:    data,
	}
}

func (wc *WebController) Block() {
	hash := wc.GetString(":splat")
	block, txs, err := repository.OneBlock(hash)
	for i, _ := range txs {
		txs[i].PayloadHex = common.ToHex(txs[i].Payload)
		meta, _ := repository.OneContractMeta(txs[i].To)
		if meta != nil {
			txs[i].ContractMeta = *meta
		}
	}

	if err != nil {
		wc.TplName = "error.tpl"
	} else {
		wc.Layout = "layout.html"
		wc.TplName = "block.tpl"
		wc.Data["Block"] = block
		wc.Data["Transactions"] = txs
	}

}

func (wc *WebController) TxsPage() {
	wc.Layout = "layout.html"
	wc.TplName = "txs.tpl"
}

func (wc *WebController) TxsLatest() {
	defer wc.ServeJSON()
	data, _ := repository.LatestTxs(DisplayNum)
	wc.Data["json"] = &Result{
		Success: true,
		Data:    data,
	}
}

func (wc *WebController) ContractPage() {
	wc.Layout = "layout.html"
	wc.TplName = "contract.tpl"
}

func (wc *WebController) ContractsLatest() {
	defer wc.ServeJSON()
	data, _ := repository.LatestContracts(DisplayNum)
	wc.Data["json"] = &Result{
		Success: true,
		Data:    data,
	}
}

func (wc *WebController) Contract() {
	hash := wc.GetString(":hash")
	wc.Layout = "layout.html"
	wc.TplName = "contract_view.tpl"
	contract, txs, _ := repository.OneContract(hash)
	for i, _ := range txs {
		txs[i].PayloadHex = common.ToHex(txs[i].Payload)
	}
	wc.Data["Contract"] = contract

	meta, _ := repository.OneContractMeta(hash)
	if meta != nil {
		wc.Data["ContractMeta"] = meta
	} else {
		wc.Data["ContractMeta"] = repository.ContractMeta{}
	}
	wc.Data["Transactions"] = txs
}

func (wc *WebController) Search() {
	hash := wc.GetString(":hash")
	beego.Info("--------------input hash: ", hash)

	reg := regexp.MustCompile(`^[1-9][0-9]*$`)

	page := 0

	if reg.MatchString(hash) {
		beego.Info("hash is block height ")
		height, err := strconv.Atoi(hash)
		if err != nil {
			beego.Info("err is : %v", err)
		}
		block, err := repository.BlockByHeight(height)
		if err == nil {
			page = 1
			hash = block.Hash
		}
	} else {

		if len(hash) == 64 {
			hash = "0x" + hash
		}

		if _, _, err := repository.OneBlock(hash); err == nil {
			page = 1
		} else if _, _, err := repository.OneContract(hash); err == nil {
			page = 2
		} else if tx, err := repository.OneTransaction(hash); err == nil {
			page = 1
			hash = tx.Block
		}
	}

	switch page {
	case 0:
		wc.Redirect("/", 302)
		break
	case 1:
		wc.Redirect("/view/blocks/hash/"+hash, 302)
		break
	case 2:
		wc.Redirect("/view/contracts/hash/"+hash, 302)
		break
	}

}
