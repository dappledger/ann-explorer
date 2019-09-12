package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/dappledger/ann-explorer/rpc"
)

var (
	BLOCK_COLLECT         = "block"
	TX_COLLECT            = "transaction"
	CONTRACT_META_COLLECT = "contract_meta"
	MONGO_URL             string
	DB_NAME               string
	ChainID               string
)

type Repo interface {
	Init()
	Save(br *BlockRepo) (err error)
	SaveContractMeta(meta ContractMeta) (err error)
	LatestBlocks(limit int) (displayData []DisplayItem, err error)
	BlocksFromTo(from, to int) (blocks []Block, err error)
	CollectionItemNum(collect string) (count int, err error)
	Contract(hash string) (tx Transaction, txs []Transaction, err error)
	Height() (maxHeight int, err error)
	Contracts(limit int) (txs []Transaction, err error)
	LatestContracts(limit int) (txs []Transaction, err error)
	TxsQuery(fromTo string) (txs []Transaction, err error)
	Txs(limit int) (txs []Transaction, err error)
	LatestTxs(limit int) (txs []Transaction, err error)
	OneContract(hash string) (contract Transaction, txs []Transaction, err error)
	OneContractMeta(hash string) (*ContractMeta, error)
	OneTransaction(hash string) (tx Transaction, err error)
	TransactionsByBlkhash(hash string) (txs []Transaction, err error)
	TxsByBlkHeight(height int) (txs []Transaction, err error)
	TransactionFromTo(from, to int) (txs []Transaction, err error)
	BlockByHeight(height int) (block Block, err error)
	OneBlock(hash string) (block Block, txs []Transaction, err error)
}

type HTTPResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      string           `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   string           `json:"error"`
}

type Status struct {
	NodeInfo          *NodeInfo `json:"node_info"`
	LatestBlockHeight int       `json:"latest_block_height"`
}

type NodeInfo struct {
	NetWork string `json:"network"`
}

func GetHTTPResp(url string) (bytez []byte, err error) {

	resp, errR := http.Get(url)
	if errR != nil {
		err = errR
		return
	}
	defer resp.Body.Close()
	bytez, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var hr HTTPResponse
	err = json.Unmarshal(bytez, &hr)
	if err != nil {
		return
	}
	if hr.Result == nil {
		err = errors.New(fmt.Sprintf("json.Unmarshal (%s)HTTPResponse wrong ,maybe you need config 'chain_id'", url))
		return
	}
	bytez, err = hr.Result.MarshalJSON()
	if err != nil {
		return
	}
	return
}

func GetStatus(chainID string) (status Status) {

	url := fmt.Sprintf("%s/status?chainid=\"%s\"", rpc.HTTP_ADDR, chainID)
	bytez, err := GetHTTPResp(url)
	if err != nil {
		log.Fatalf("GetHTTPResp failed: %s", err.Error())
	}

	err = json.Unmarshal(bytez, &status)
	if err != nil {
		log.Fatalf("json.Unmarshal(Status) failed: %s", err.Error())
	}
	return
}

var deleteRepo Repo

func Init() {

	if beego.AppConfig.String("mogo_addr") != "" {
		deleteRepo = &monRepo{}
	} else {
		deleteRepo = &sqliteRepo{}
	}
	deleteRepo.Init()
	return
}

type BlockRepo struct {
	Blocks []Block
	Txs    []Transaction
}

type DisplayItem struct {
	Block
	Tps      int
	Interval float64
}

type Block struct {
	Hash            string
	ParentHash      string
	ChainID         string
	Height          int `bson:"_id"`
	Time            time.Time
	NumTxs          int
	LastCommitHash  string
	DataHash        string
	ValidatorsHash  string
	AppHash         string
	ProposerAddress string
}

type ContractMeta struct {
	Hash string `json:"hash" bson:"_id"`
	ABI  string `json:"abi" bson:"abi"`
}

type Transaction struct {
	Payload      []byte `json:"-"`
	PayloadHex   string
	ContractMeta ContractMeta `bson:"-"`
	Hash         string       `bson:"_id"`
	From         string
	To           string
	Receipt      string
	Amount       string
	Nonce        uint64
	Gas          string
	Size         int64
	Block        string
	Contract     string
	Time         time.Time
	Height       int
}

func Save(br *BlockRepo) (err error) {
	err = deleteRepo.Save(br)
	return
}
func LatestBlocks(limit int) (displayData []DisplayItem, err error) {
	displayData, err = deleteRepo.LatestBlocks(limit)
	return
}
func BlocksFromTo(from, to int) (blocks []Block, err error) {
	blocks, err = deleteRepo.BlocksFromTo(from, to)
	return
}
func CollectionItemNum(collect string) (count int, err error) {
	count, err = deleteRepo.CollectionItemNum(collect)
	return
}
func Contract(hash string) (contract Transaction, txs []Transaction, err error) {
	contract, txs, err = deleteRepo.Contract(hash)
	return
}
func Height() (maxHeight int, err error) {
	maxHeight, err = deleteRepo.Height()
	return
}

func LatestContracts(limit int) (txs []Transaction, err error) {
	txs, err = deleteRepo.LatestContracts(limit)
	return
}
func TxsQuery(fromTo string) (txs []Transaction, err error) {

	txs, err = deleteRepo.TxsQuery(fromTo)

	return
}

func LatestTxs(limit int) (txs []Transaction, err error) {
	txs, err = deleteRepo.LatestTxs(limit)
	return
}

func Txs(limit int) (txs []Transaction, err error) {
	txs, err = deleteRepo.Txs(limit)
	return
}
func OneContract(hash string) (contract Transaction, txs []Transaction, err error) {
	contract, txs, err = deleteRepo.OneContract(hash)
	return
}

func OneContractMeta(hash string) (meta *ContractMeta, err error) {
	return deleteRepo.OneContractMeta(hash)
}

func OneTransaction(hash string) (tx Transaction, err error) {
	tx, err = deleteRepo.OneTransaction(hash)
	return
}
func TxsByBlkHeight(height int) (txs []Transaction, err error) {
	txs, err = deleteRepo.TxsByBlkHeight(height)
	return
}
func TransactionsByBlkhash(hash string) (txs []Transaction, err error) {
	txs, err = deleteRepo.TransactionsByBlkhash(hash)
	return
}
func TransactionFromTo(from, to int) (txs []Transaction, err error) {
	txs, err = deleteRepo.TransactionFromTo(from, to)
	return
}
func BlockByHeight(height int) (block Block, err error) {
	block, err = deleteRepo.BlockByHeight(height)
	return
}
func OneBlock(hash string) (block Block, txs []Transaction, err error) {
	block, txs, err = deleteRepo.OneBlock(hash)
	return
}
func SaveContractMeta(meta ContractMeta) error {
	return deleteRepo.SaveContractMeta(meta)
}

func TryLock(owner string, long int) (bool, error) {
	m, ok := deleteRepo.(*monRepo)
	if ok {
		return m.TryLock(owner, long)
	}
	return true, nil
}
