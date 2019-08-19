package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/dappledger/AnnChain/eth/common"
	"github.com/dappledger/ann-explorer/repository"
)

type ApiController struct {
	beego.Controller
}

type BlockListResp struct {
	Total  int
	Blocks []repository.Block
}

type TxListResp struct {
	Total int
	Txs   []repository.Transaction
}

func (p *ApiController) QueryBlock() {
	param := p.Ctx.Input.Query("param")
	defer p.ServeJSON()

	var err error
	var bk repository.Block
	if strings.HasPrefix(param, "0x") {
		bk, _, err = repository.OneBlock(param)
		if err != nil {
			p.Data["json"] = &Result{
				Success: false,
				Data:    "get block from repo failed",
			}
			return
		}
	} else {
		height, err := strconv.Atoi(param)
		if err != nil {
			p.Data["json"] = &Result{
				Success: false,
				Data:    "input format error.",
			}
			return
		}
		bk, _ = repository.BlockByHeight(height)
	}

	if bk.AppHash == "" {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "can not find block:" + param,
		}
		return
	}

	p.Data["json"] = &Result{
		Success: true,
		Data:    bk,
	}
}

func (p *ApiController) QueryBlockList() {
	from := p.Ctx.Input.Query("from")
	to := p.Ctx.Input.Query("to")
	sort := p.Ctx.Input.Query("sort")
	defer p.ServeJSON()

	htFrom, err := strconv.Atoi(from)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "input format error.",
		}
		return
	}
	htTo, err := strconv.Atoi(to)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "input format error.",
		}
		return
	}
	if htTo-htFrom <= 0 || htFrom < 1 {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "to should be bigger than from",
		}
		return
	}

	var resp BlockListResp
	resp.Total, err = repository.Height()
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get height from repo failed",
		}
		return
	}
	if sort == "desc" {
		htFrom, htTo = reverse(resp.Total, htFrom, htTo)
	}
	resp.Blocks, err = repository.BlocksFromTo(htFrom, htTo)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get blocks from repo failed",
		}
		return
	}

	p.Data["json"] = &Result{
		Success: true,
		Data:    resp,
	}
}

func reverse(total, from, to int) (nFrom, nTo int) {
	nFrom = total - to + 1
	nTo = total - from + 1
	return
}

func (p *ApiController) QueryTx() {
	hash := p.Ctx.Input.Query("hash")
	defer p.ServeJSON()

	tx, err := repository.OneTransaction(hash)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get tx from repo failed",
		}
		return
	}
	tx.PayloadHex = common.ToHex(tx.Payload)
	p.Data["json"] = &Result{
		Success: true,
		Data:    tx,
	}
}

func (p *ApiController) QueryTxsList() {
	from := p.Ctx.Input.Query("from")
	to := p.Ctx.Input.Query("to")
	defer p.ServeJSON()

	htFrom, err := strconv.Atoi(from)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "input format error.",
		}
		return
	}
	htTo, err := strconv.Atoi(to)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "input format error.",
		}
		return
	}
	if htTo-htFrom <= 0 || htFrom < 1 {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "to should be bigger than from",
		}
		return
	}

	var resp TxListResp
	resp.Total, err = repository.CollectionItemNum(repository.TX_COLLECT)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get collection num error.",
		}
		return
	}
	resp.Txs, err = repository.TransactionFromTo(htFrom, htTo)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get txs from repo failed",
		}
		return
	}
	for i := 0; i < len(resp.Txs); i++ {
		resp.Txs[i].PayloadHex = common.ToHex(resp.Txs[i].Payload)
	}

	p.Data["json"] = &Result{
		Success: true,
		Data:    resp,
	}
}

func (p *ApiController) QueryTxsByBlkHeight() {

	height := p.Ctx.Input.Query("height")
	defer p.ServeJSON()

	var (
		err  error
		resp TxListResp
		hi   int
	)
	hi, err = strconv.Atoi(height)
	if err != nil || hi <= 0 {
		p.respErr(fmt.Sprintf("invalid input: %s", height))
		return
	}
	resp.Txs, err = repository.TxsByBlkHeight(hi)
	if err != nil {
		p.respErr(fmt.Sprintf("get txs failed: %v", err))
		return
	}
	resp.Total = len(resp.Txs)
	for i := 0; i < resp.Total; i++ {
		resp.Txs[i].PayloadHex = common.ToHex(resp.Txs[i].Payload)
	}

	p.Data["json"] = &Result{
		Success: true,
		Data:    resp,
	}
}

func (p *ApiController) respErr(data string) {

	p.Data["json"] = &Result{
		Success: false,
		Data:    data,
	}
}

func (p *ApiController) QueryTxsListByBlk() {
	hash := p.Ctx.Input.Query("hash")
	defer p.ServeJSON()

	var err error
	var resp TxListResp
	resp.Txs, err = repository.TransactionsByBlkhash(hash)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get txs by block hash from repo failed",
		}
		return
	}
	resp.Total = len(resp.Txs)
	for i := 0; i < len(resp.Txs); i++ {
		resp.Txs[i].PayloadHex = common.ToHex(resp.Txs[i].Payload)
	}

	p.Data["json"] = &Result{
		Success: true,
		Data:    resp,
	}
}

func (p *ApiController) SaveContractMeta() {
	defer p.ServeJSON()

	m := repository.ContractMeta{}
	r := Result{}
	p.Data["json"] = &r

	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&m); err != nil {
		r.Data = err.Error()
		return
	}

	if m.ABI == "" || m.Hash == "" {
		r.Success = false
		r.Data = "invalid params"
		return
	}

	if err := repository.SaveContractMeta(m); err != nil {
		r.Data = err.Error()
		return
	}
	r.Success = true
	r.Data = &m
}
