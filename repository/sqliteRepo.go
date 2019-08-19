package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
)

type sqliteRepo struct {
	Repo
}

func (this *sqliteRepo) Init() {
	ChainID = beego.AppConfig.String("chain_id")
	if ChainID == "" {
		ChainID = GetStatus("").NodeInfo.NetWork
	}
	BLOCK_COLLECT += "_" + ChainID
	TX_COLLECT += "_" + ChainID
	BLOCK_COLLECT = strings.Replace(BLOCK_COLLECT, "-", "_", -1)
	TX_COLLECT = strings.Replace(TX_COLLECT, "-", "_", -1)
	fmt.Println("BLOCK_COLLECT = ", BLOCK_COLLECT)
	fmt.Println("TX_COLLECT = ", TX_COLLECT)

	DB_NAME = "block_browser"
	err := CreateSqlite()
	if err != nil {
		beego.Error(err)
	}
}

func (this *sqliteRepo) Save(br *BlockRepo) (err error) {
	err = SaveBlockBySqlite(br.Blocks)
	if err != nil {
		err = errors.New("insert block failed: " + err.Error())
		return
	}
	if len(br.Txs) > 0 {
		err = SaveTxBySqlite(br.Txs)
		if err != nil {
			err = errors.New("insert tx failed: " + err.Error())
			return
		}
	}
	return
}

func (this *sqliteRepo) SaveContractMeta(meta ContractMeta) (err error) {

	sqlTx, errB := db.Begin()
	defer func() {
		if err != nil {
			sqlTx.Rollback()
		}
	}()
	if errB != nil {
		err = errB
		return
	}
	stmt, errP := sqlTx.Prepare(contractMetaInsertSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(meta.Hash, meta.ABI)
	if err != nil {
		return err
	}
	err = sqlTx.Commit()
	return
}

func (this *sqliteRepo) CollectionItemNum(collect string) (count int, err error) {
	stmt, errP := db.Prepare(fmt.Sprintf(countSQL, collect))
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&count)
	return
}

func (this *sqliteRepo) OneTransaction(hash string) (tx Transaction, err error) {
	stmt, errP := db.Prepare(txHashSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(hash).Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
	return
}

func (this *sqliteRepo) TransactionsByBlkhash(hash string) (txs []Transaction, err error) {
	txs, err = TransactionsByBlkhashBySqlite(hash)
	return
}

func (this *sqliteRepo) OneContract(hash string) (tx Transaction, txs []Transaction, err error) {
	stmt, errP := db.Prepare(txContractSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(hash).Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
	if err != nil {
		return
	}

	txs, err = txsByTo(hash)

	return
}

func (this *sqliteRepo) OneContractMeta(hash string) (meta *ContractMeta, err error) {
	stmt, errP := db.Prepare(contractMetaSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	meta = new(ContractMeta)

	err = stmt.QueryRow(hash).Scan(&meta.Hash, &meta.ABI)
	if err != nil {
		return nil, err
	}
	return
}

func (this *sqliteRepo) LatestTxs(limit int) (txs []Transaction, err error) {
	rows, errQ := db.Query(txsNoContractLatestSQL, limit)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()

	return
}

func (this *sqliteRepo) Txs(limit int) (txs []Transaction, err error) {
	rows, errQ := db.Query(txsNoContractSQL, limit)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()

	return
}

func (this *sqliteRepo) Height() (maxHeight int, err error) {

	stmt, errP := db.Prepare(blockMaxHeightSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	if err := stmt.QueryRow().Scan(&maxHeight); err != nil {
		return 0, nil
	}
	return maxHeight, nil
}

func (this *sqliteRepo) TxsQuery(fromTo string) (txs []Transaction, err error) {
	rows, errQ := db.Query(txsByFromOrToSQL, fromTo)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()
	return
}

func (this *sqliteRepo) Contracts(limit int) (txs []Transaction, err error) {
	rows, errQ := db.Query(txsContractSQL, limit)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()

	return
}

func (this *sqliteRepo) LatestContracts(limit int) (txs []Transaction, err error) {
	rows, errQ := db.Query(txsContractLatestSQL, limit)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()

	return
}

func (this *sqliteRepo) BlockByHeight(height int) (block Block, err error) {
	block, err = BlockByHeightBySqlite(height)
	return
}

func (this *sqliteRepo) OneBlock(hash string) (block Block, txs []Transaction, err error) {
	block, txs, err = OneBlockBySqlite(hash)
	return
}

func (this *sqliteRepo) BlocksFromTo(from, to int) (blocks []Block, err error) {
	blocks, err = BlocksFromToBySqlite(from, to)
	return
}

func (this *sqliteRepo) LatestBlocks(limit int) (displayData []DisplayItem, err error) {
	blocks, err := LatestBlocksBySqlite(limit)
	if err != nil {
		return
	}
	dur := blocks[0].Time.Sub(blocks[limit-1].Time).Seconds()
	var totalTxs int
	for _, v := range blocks {
		totalTxs += v.NumTxs
	}
	interval := dur / float64(limit)
	tps := int(float64(totalTxs) / dur)
	for _, v := range blocks {
		displayData = append(displayData, DisplayItem{v, tps, interval})
	}
	return
}

func (this *sqliteRepo) Contract(hash string) (tx Transaction, txs []Transaction, err error) {
	tx, err = OneTransactionBySqlite(hash)
	if err != nil {
		return
	}
	txs, err = txsByTo(hash)
	return
}

func (this *sqliteRepo) TransactionFromTo(from, to int) (txs []Transaction, err error) {
	txs, err = TransactionFromToBySqlite(from, to)
	return
}

func BlocksFromToBySqlite(from, to int) (blocks []Block, err error) {

	rows, errQ := db.Query(blockRangeSQL, from, to)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var block Block
		err = rows.Scan(&block.Height, &block.Hash, &block.ParentHash, &block.ChainID, &block.Time, &block.NumTxs, &block.LastCommitHash, &block.DataHash, &block.ValidatorsHash, &block.AppHash, &block.ProposerAddress)
		if err != nil {
			return
		}
		blocks = append(blocks, block)
	}
	err = rows.Err()
	return
}
func OneBlockBySqlite(hash string) (block Block, txs []Transaction, err error) {

	stmt, errP := db.Prepare(blockHashSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(hash).Scan(&block.Height, &block.Hash, &block.ParentHash, &block.ChainID, &block.Time, &block.NumTxs, &block.LastCommitHash, &block.DataHash, &block.ValidatorsHash, &block.AppHash, &block.ProposerAddress)
	if err != nil {
		return
	}
	txs, err = TransactionsByBlkhashBySqlite(hash)
	return
}

func BlockByHeightBySqlite(height int) (block Block, err error) {

	stmt, errP := db.Prepare(blockHeightSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(height).Scan(&block.Height, &block.Hash, &block.ParentHash, &block.ChainID, &block.Time, &block.NumTxs, &block.LastCommitHash, &block.DataHash, &block.ValidatorsHash, &block.AppHash, &block.ProposerAddress)
	return
}

func TransactionFromToBySqlite(from, to int) (txs []Transaction, err error) {
	return
}

func (this *sqliteRepo) TxsByBlkHeight(height int) (txs []Transaction, err error) {
	rows, errQ := db.Query(txsByBlockHeightSQL, height)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()
	return
}

func TransactionsByBlkhashBySqlite(hash string) (txs []Transaction, err error) {

	rows, errQ := db.Query(txsByBlockHashSQL, hash)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()
	return
}

func txsByTo(hash string) (txs []Transaction, err error) {

	rows, errQ := db.Query(txsByToSQL, hash)
	if errQ != nil {
		err = errQ
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	err = rows.Err()
	return
}

func OneTransactionBySqlite(hash string) (tx Transaction, err error) {

	stmt, errP := db.Prepare(txHashSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(hash).Scan(&tx.Hash, &tx.Payload, &tx.PayloadHex, &tx.From, &tx.To, &tx.Receipt, &tx.Amount, &tx.Nonce, &tx.Gas, &tx.Size, &tx.Block, &tx.Contract, &tx.Time, &tx.Height)
	return
}
