package job

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/dappledger/AnnChain/eth/common"
	"github.com/dappledger/AnnChain/eth/crypto"
	"github.com/dappledger/AnnChain/eth/rlp"
	"github.com/dappledger/AnnChain/gemmill/types"
	"github.com/dappledger/ann-explorer/repository"
	"github.com/dappledger/ann-explorer/rpc"
)

const (
	IntervalInSec = 3
	Step          = 500
)

type JsonTime time.Time

func (j *JsonTime) UnmarshalJSON(data []byte) error {

	dataStr := string(data)
	nano, err := strconv.ParseInt(dataStr, 10, 64)
	if err != nil {
		// try parse UTC format time
		if len(dataStr) > 10 && strings.HasPrefix(dataStr, "\"") && strings.HasSuffix(dataStr, "\"") {
			dataStr = dataStr[1 : len(dataStr)-1]
			it, err := time.Parse(time.RFC3339, dataStr)
			if err != nil {
				return err
			}
			*j = JsonTime(it)
			return nil
		}
		return err
	}
	*j = JsonTime(time.Unix(0, nano))
	return nil
}

func calcSyncFromHeight() (int, error) {

	h, err := repository.Height()
	if err != nil {
		return 0, err
	}

	fromStr := beego.AppConfig.String("sync_from_height")
	if strings.ToUpper(fromStr) == "LATEST" {
		latestHeight := repository.GetStatus(repository.ChainID).LatestBlockHeight
		return latestHeight, nil
	} else {
		fromHeight, err := strconv.Atoi(fromStr)
		if err != nil {
			return h, nil
		}
		if fromHeight < h {
			return h, nil
		}
		return fromHeight, nil
	}
}

func SyncTimingTask() {

	owner := fmt.Sprintf("%d-%d", os.Getpid(), time.Now().UnixNano())
	for {

		func() {
			if ok, _ := repository.TryLock(owner, IntervalInSec); !ok {
				return
			}
			h, err := calcSyncFromHeight()
			if err != nil {
				beego.Error("sync failed, calculate from_height err", err)
				time.Sleep(time.Second * 3)
				return
			}
			beego.Info(time.Now())
			blockChain(h)
		}()

		time.Sleep(IntervalInSec * time.Second)
	}
}

type Metas struct {
	BlockMetas []BlockMeta `json:"block_metas"`
}

type BlockMeta struct {
	Hash   string  `json:"hash"`   // The block hash
	Header *Header `json:"header"` // The block's Header
}
type Header struct {
	ChainID         string   `json:"chain_id"`
	Height          int      `json:"height"`
	Time            JsonTime `json:"time"`
	NumTxs          int      `json:"num_txs"`          // XXX: Can we get rid of this?
	LastCommitHash  string   `json:"last_commit_hash"` // commit from validators from the last block
	DataHash        string   `json:"data_hash"`        // transactions
	ValidatorsHash  string   `json:"validators_hash"`  // validators for the current block
	AppHash         string   `json:"app_hash"`         // state after txs from the previous block
	ReceiptsHash    string   `json:"recepits_hash"`    // recepits_hash from previous block
	LastBlockID     BlockID  `json:"last_block_id"`
	ProposerAddress string   `json:"proposer_address"`
}

type BlockID struct {
	Hash string `json:"hash"`
}

func decodeTxString(tx string) []byte {
	return common.FromHex(tx)
}

func toLower0xHex(str string) string {
	str = strings.ToLower(str)
	if !strings.HasPrefix(str, "0x") {
		return "0x" + str
	}
	return str
}

func blockChain(h int) int {
	beego.Info("current block height : ", h)
	url := fmt.Sprintf("%s/blockchain?minHeight=%d&maxHeight=%d&chainid=\"%s\"", rpc.HTTP_ADDR, h+1, h+Step, repository.ChainID)

	bytez, err := repository.GetHTTPResp(url)
	if err != nil {
		beego.Info(err)
		return 0
	}
	var metas Metas
	err = json.Unmarshal(bytez, &metas)
	if err != nil {
		beego.Error(err.Error())
		return 0
	}
	//save block
	br := &repository.BlockRepo{
		Blocks: []repository.Block{},
		Txs:    []repository.Transaction{},
	}
	for _, o := range metas.BlockMetas {
		rb := repository.Block{
			Hash:            toLower0xHex(o.Hash),
			ParentHash:      toLower0xHex(o.Header.LastBlockID.Hash),
			ChainID:         o.Header.ChainID,
			Height:          o.Header.Height,
			Time:            time.Time(o.Header.Time),
			NumTxs:          o.Header.NumTxs,
			LastCommitHash:  toLower0xHex(o.Header.LastCommitHash),
			DataHash:        toLower0xHex(o.Header.DataHash),
			ValidatorsHash:  toLower0xHex(o.Header.ValidatorsHash),
			AppHash:         toLower0xHex(o.Header.AppHash),
			ProposerAddress: toLower0xHex(o.Header.ProposerAddress),
		}
		br.Blocks = append(br.Blocks, rb)
		if o.Header.NumTxs > 0 {
			resultBlock, err := GetBlock(o.Header.Height)
			if err != nil {
				beego.Error("GetBlock(height:%d) Error :%v\n", o.Header.Height, err)
				return 0
			} else {
				for _, v := range resultBlock.Block.Data.Txs {
					tx := new(Transaction)
					data := decodeTxString(v)
					err := rlp.DecodeBytes(data, tx)
					if err != nil {
						beego.Error("Decode Transaction tx bytes: [%v], error : %v\n", v, err)
						return 0
					}

					var signer Signer = new(HomesteadSigner)

					from, err := Sender(signer, tx)
					if err != nil {
						beego.Error("Error : Get Transaction From error: %v\n", err)
						return 0
					}

					rtx := repository.Transaction{
						Payload: tx.Data(),
						Hash:    common.ToHex(types.Tx(data).Hash()),
						From:    from.Hex(),
						Amount:  new(big.Int).Set(tx.Value()).String(),
						Nonce:   tx.Nonce(),
						Size:    int64(tx.Size()),
						Block:   rb.Hash,
						Time:    rb.Time,
						Height:  o.Header.Height,
					}

					to := tx.To()
					if to == nil {
						// contract create
						contractAddr := crypto.CreateAddress(from, tx.Nonce())
						rtx.Contract = contractAddr.Hex()
					} else {
						rtx.To = to.Hex()
					}
					br.Txs = append(br.Txs, rtx)
				}
			}
		}
	}

	err = repository.Save(br)
	if err != nil {
		beego.Error("[save failed] %s\n", err.Error())
		return 0
	}

	beego.Info("  save blocks len:", len(br.Blocks))
	beego.Info("  save txs len:", len(br.Txs))
	return len(br.Blocks)
}

type Block struct {
	*Header `json:"header"`
	*Data   `json:"data"`
}

type ResultBlock struct {
	BlockMeta *BlockMeta `json:"block_meta"`
	Block     *Block     `json:"block"`
}

type Data struct {
	Txs []string `json:"txs"`

	// Volatile
	hash []byte
}

func GetBlock(height int) (result ResultBlock, err error) {
	url := fmt.Sprintf("%s/block?height=%d&chainid=\"%s\"", rpc.HTTP_ADDR, height, repository.ChainID)
	bytez, errB := repository.GetHTTPResp(url)
	if errB != nil {
		err = errB
		return
	}
	err = json.Unmarshal(bytez, &result)
	return
}
