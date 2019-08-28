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

package controllers

import (
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/dappledger/AnnChain/eth/common"
	"github.com/dappledger/ann-explorer/repository"
)

const (
	DisplayNum   = 50
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

func (wc *WebController) Blocks() {
	defer wc.ServeJSON()

	p := getPage(wc)

	count, err := repository.CollectionItemNum(repository.BLOCK_COLLECT)
	if err != nil {
		beego.Error(err)
		return
	}

	from, to := calcFromTo(p, count)

	blocks, err := repository.BlocksFromTo(from, to)
	if err != nil {
		beego.Error(err)
		return
	}

	wc.Data["json"] = &Result{
		Success: true,
		Data: &PaginationResult{
			Data: repository.BlocksToDisplayItems(blocks, DisplayNum),
			Page: Pagination{
				Current: p,
				Total:   totalPage(count),
			},
		},
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

	p := getPage(wc)

	count, err := repository.TxsCount()
	if err != nil {
		return
	}
	skip := (p - 1) * DisplayNum

	data, _ := repository.LatestTxs(DisplayNum, skip)
	wc.Data["json"] = &Result{
		Success: true,
		Data: PaginationResult{
			Data: data,
			Page: Pagination{
				Current: p,
				Total:   totalPage(count),
			},
		},
	}
}

func (wc *WebController) ContractPage() {
	wc.Layout = "layout.html"
	wc.TplName = "contract.tpl"
}

func (wc *WebController) ContractsLatest() {
	defer wc.ServeJSON()

	p := getPage(wc)

	count, err := repository.ContractsCount()
	if err != nil {
		return
	}
	skip := (p - 1) * DisplayNum

	data, _ := repository.LatestContracts(DisplayNum, skip)
	wc.Data["json"] = &Result{
		Success: true,
		Data: PaginationResult{
			Data: data,
			Page: Pagination{
				Current: p,
				Total:   totalPage(count),
			},
		},
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
