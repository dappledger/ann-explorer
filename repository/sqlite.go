package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createBlockCollectSQL = `CREATE TABLE IF NOT EXISTS block_t
    (
		height				INTEGER NOT NULL,
		hash				VARCHAR(64) PRIMARY KEY,
		parent_hash			VARCHAR(64) ,
		chain_id			VARCHAR(64) NOT NULL,
		time				DATETIME NOT NULL,
		num_txs				INTEGER NOT NULL,
		last_commit_hash    VARCHAR(64) NOT NULL,
		data_hash			VARCHAR(64) NOT NULL,
		validators_hash		VARCHAR(64) NOT NULL,
		app_hash			VARCHAR(64) NOT NULL,
		proposer_address  VARCHAR(64) NOT NULL
	);`

	createTxCollectSQL = `CREATE TABLE IF NOT EXISTS transaction_t
    (
		hash			VARCHAR(64) PRIMARY KEY,
		payload			VARCHAR(64) NOT NULL,
		payload_hex		VARCHAR(64) NOT NULL,
		from_addr		VARCHAR(64) NOT NULL,
		to_addr			VARCHAR(64) NOT NULL,
		receipt			VARCHAR(64) NOT NULL,
		amount			VARCHAR(64) NOT NULL,
		nonce			INTEGER NOT NULL,
		gas				VARCHAR(64) NOT NULL,
		size			INTEGER NOT NULL,
		block			VARCHAR(64) NOT NULL,
		contract		VARCHAR(64),
		time			DATETIME NOT NULL,
		height			INTEGER NOT NULL
    );`

	createContractMetaCollectSQL = `CREATE TABLE IF NOT EXISTS contract_meta_t
    (
		hash			VARCHAR(64) PRIMARY KEY,
		abi				TEXT
    );`

	blockSortSQL      = `SELECT * FROM block_t ORDER BY height DESC limit ?`
	blockRangeSQL     = `SELECT * FROM block_t WHERE height >= ? AND height <= ?`
	blockHashSQL      = `SELECT * FROM block_t WHERE hash = ?`
	blockHeightSQL    = `SELECT * FROM block_t WHERE height = ?`
	blockMaxHeightSQL = `SELECT MAX(height) FROM block_t`
	blockInsertSQL    = `INSERT INTO block_t VALUES(?,?,?,?,?,?,?,?,?,?,?)`

	txsByBlockHashSQL      = `SELECT * FROM transaction_t WHERE block = ?`
	txsByBlockHeightSQL    = `SELECT * FROM transaction_t WHERE height = ?`
	txsByToSQL             = `SELECT * FROM transaction_t WHERE to_addr = ?`
	txsByFromOrToSQL       = `SELECT * FROM transaction_t WHERE to_addr = ? OR from_addr = ?`
	txsNoContractSQL       = `SELECT * FROM transaction_t WHERE contract = "" limit ?`
	txsNoContractLatestSQL = `SELECT * FROM transaction_t WHERE contract = "" ORDER BY height DESC, nonce DESC limit ? `
	txsContractSQL         = `SELECT * FROM transaction_t WHERE contract != "" limit ?`
	txsContractLatestSQL   = `SELECT * FROM transaction_t WHERE contract != "" ORDER BY height DESC, nonce DESC limit ?`
	txHashSQL              = `SELECT * FROM transaction_t WHERE hash = ?`
	txContractSQL          = `SELECT * FROM transaction_t WHERE contract = ?`
	txInsertSQL            = `INSERT INTO transaction_t VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	contractMetaSQL        = `SELECT * FROM contract_meta_t WHERE hash = ?`
	contractMetaInsertSQL  = `INSERT INTO contract_meta_t VALUES(?,?)`

	countSQL = "SELECT COUNT(hash) FROM %s"
)

var db *sql.DB

func CreateSqlite() (err error) {
	db, err = sql.Open("sqlite3", fmt.Sprintf("./%s.db", DB_NAME))
	if err != nil {
		return
	}
	_, err = db.Exec(createBlockCollectSQL)
	if err != nil {
		return
	}
	_, err = db.Exec(createTxCollectSQL)
	if err != nil {
		return
	}
	_, err = db.Exec(createContractMetaCollectSQL)
	return
}

func LatestBlocksBySqlite(limit int) (blocks []Block, err error) {

	rows, errQ := db.Query(blockSortSQL, limit)
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

func SaveBlockBySqlite(blocks []Block) (err error) {

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
	stmt, errP := sqlTx.Prepare(blockInsertSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	for _, block := range blocks {
		_, err = stmt.Exec(
			block.Height,
			block.Hash,
			block.ParentHash,
			block.ChainID,
			block.Time,
			block.NumTxs,
			block.LastCommitHash,
			block.DataHash,
			block.ValidatorsHash,
			block.AppHash,
			block.ProposerAddress,
		)
		if err != nil {
			return
		}
	}
	err = sqlTx.Commit()
	return
}

func SaveTxBySqlite(txs []Transaction) (err error) {

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
	stmt, errP := sqlTx.Prepare(txInsertSQL)
	if errP != nil {
		err = errP
		return
	}
	defer stmt.Close()
	for _, tx := range txs {
		_, err = stmt.Exec(
			tx.Hash,
			tx.Payload,
			tx.PayloadHex,
			tx.From,
			tx.To,
			tx.Receipt,
			tx.Amount,
			tx.Nonce,
			tx.Gas,
			tx.Size,
			tx.Block,
			tx.Contract,
			tx.Time,
			tx.Height,
		)
		if err != nil {
			return
		}
	}
	err = sqlTx.Commit()
	return
}
