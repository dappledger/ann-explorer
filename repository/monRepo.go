package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type monRepo struct {
	*mgo.Session
}

func (m *monRepo) ensureIndex() (err error) {
	c := m.blockCollection()
	index := mgo.Index{
		Key: []string{"height"},
	}
	err = c.EnsureIndex(index)
	return
}

var MutexCollect = "mutex"

func (monRepo *monRepo) Init() {
	ChainID = beego.AppConfig.String("chain_id")
	if ChainID == "" {
		ChainID = GetStatus("").NodeInfo.NetWork
	}

	collPrefix := beego.AppConfig.String("collection_prefix")
	BLOCK_COLLECT = collPrefix + BLOCK_COLLECT + "_" + ChainID
	TX_COLLECT = collPrefix + TX_COLLECT + "_" + ChainID
	CONTRACT_META_COLLECT = collPrefix + CONTRACT_META_COLLECT + "_" + ChainID
	MutexCollect = collPrefix + MutexCollect + "_" + ChainID

	BLOCK_COLLECT = strings.Replace(BLOCK_COLLECT, "-", "_", -1)
	TX_COLLECT = strings.Replace(TX_COLLECT, "-", "_", -1)
	CONTRACT_META_COLLECT = strings.Replace(CONTRACT_META_COLLECT, "-", "_", -1)
	fmt.Println("BLOCK_COLLECT = ", BLOCK_COLLECT)
	fmt.Println("TX_COLLECT = ", TX_COLLECT)

	DB_NAME = beego.AppConfig.String("mogo_db")
	if DB_NAME == "" {
		DB_NAME = "block_browser"
	}

	if beego.AppConfig.String("mogo_addr") != "" {
		if beego.AppConfig.String("mogo_user") != "" {
			MONGO_URL = "mongodb://" +
				beego.AppConfig.String("mogo_user") + ":" +
				beego.AppConfig.String("mogo_pwd") + "@" +
				beego.AppConfig.String("mogo_addr") + "/" +
				DB_NAME
		} else {
			MONGO_URL = beego.AppConfig.String("mogo_addr")
		}

		session, err := mgo.Dial(MONGO_URL)
		if err != nil {
			beego.Error(err)
		}
		monRepo.Session = session
		err = monRepo.ensureIndex()
		if err != nil {
			beego.Error(err)
		}
	}
}

func (m *monRepo) blockCollection() *mgo.Collection {
	return m.Session.Clone().DB(DB_NAME).C(BLOCK_COLLECT);
}

func (m *monRepo) txCollection() *mgo.Collection {
	return m.Session.Clone().DB(DB_NAME).C(TX_COLLECT);
}

func (m *monRepo) contractMetaCollection() *mgo.Collection {
	return m.Session.Clone().DB(DB_NAME).C(CONTRACT_META_COLLECT);
}

func (m *monRepo) Save(br *BlockRepo) (err error) {

	c := m.blockCollection()
	defer c.Database.Session.Close()
	for _, b := range br.Blocks {
		err = c.Insert(b)
		if err != nil {
			err = errors.New("insert block failed: " + err.Error())
			return
		}
	}

	c2 := m.txCollection()
	defer c2.Database.Session.Close()
	for _, tx := range br.Txs {
		err = c2.Insert(tx)
		if err != nil {
			err = errors.New("insert tx failed: " + err.Error())
			return
		}
	}
	return
}

func (m *monRepo) SaveContractMeta(meta ContractMeta) error {
	c := m.contractMetaCollection()
	defer c.Database.Session.Close()

	return c.Insert(&meta)
}

func (m *monRepo) CollectionItemNum(collect string) (count int, err error) {
	c := m.Session.Clone().DB(DB_NAME).C(collect)
	defer c.Database.Session.Close()
	return c.Count()
}

func (m *monRepo) Height() (maxHeight int, err error) {
	c := m.blockCollection()
	defer c.Database.Session.Close()
	query := c.Find(nil).Sort("-_id").Limit(1)
	result := Block{}
	err = query.One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return 0, nil
		}
		beego.Info("Info : query Block Max Height , err: %v", err)
	}
	maxHeight = result.Height
	return
}

func (m *monRepo) Contract(hash string) (tx Transaction, txs []Transaction, err error) {

	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"contract": hash}).Limit(1)
	err = query.One(&tx)
	if err != nil {
		return
	}
	query2 := c.Find(bson.M{"to": hash})
	err = query2.All(&txs)
	return
}

func (m *monRepo) LatestContracts(limit int, skip int) (txs []Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"contract": bson.M{"$ne": ""}}).Sort("-height", "-nonce").Skip(skip).Limit(limit)
	err = query.All(&txs)
	return
}

func (m *monRepo) Contracts(limit int) (txs []Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"contract": bson.M{"$ne": ""}}).Limit(limit)
	err = query.All(&txs)
	return
}

func (m *monRepo) ContractsCount() (int, error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	return c.Find(bson.M{"contract": bson.M{"$ne": ""}}).Count()
}

func (m *monRepo) TxsCount() (int, error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	return c.Find(bson.M{"contract": ""}).Count()
}

func (m *monRepo) TxsQuery(fromTo string) (txs []Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"$or": []bson.M{{"from": bson.M{"$eq": fromTo}}, {"to": bson.M{"$eq": fromTo}}}})
	err = query.All(&txs)
	return
}

func (m *monRepo) LatestTxs(limit int, skip int) (txs []Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"contract": ""}).Sort("-height", "-nonce").Skip(skip).Limit(limit)
	err = query.All(&txs)
	return
}

func (m *monRepo) Txs(limit int) (txs []Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"contract": bson.M{"$eq": ""}}).Limit(limit)
	err = query.All(&txs)
	return
}

func (m *monRepo) OneContract(hash string) (contract Transaction, txs []Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()

	query := c.Find(bson.M{"contract": hash})
	err = query.One(&contract)
	if err != nil {
		return
	}
	query2 := c.Find(bson.M{"to": hash})
	err = query2.All(&txs)
	return
}

func (m *monRepo) OneContractMeta(hash string) (*ContractMeta, error) {

	c3 := m.contractMetaCollection()
	defer c3.Database.Session.Close()
	meta := new(ContractMeta)
	if err := c3.Find(bson.M{"_id": hash}).One(meta); err != nil {
		if err != mgo.ErrNotFound {
			return nil, err
		}
		meta = nil
	}
	return meta, nil
}

func (m *monRepo) OneTransaction(hash string) (tx Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"_id": hash})
	err = query.One(&tx)
	return
}

func (m *monRepo) TxsByBlkHeight(height int) (txs []Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"height": height})
	err = query.All(&txs)
	return
}

func (m *monRepo) TransactionsByBlkhash(hash string) (txs []Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"block": hash})
	err = query.All(&txs)
	return
}

func (m *monRepo) TransactionFromTo(from, to int) (txs []Transaction, err error) {
	c := m.txCollection()
	defer c.Database.Session.Close()
	query := c.Find(nil).Skip(from - 1).Limit(to - from + 1)
	err = query.All(&txs)
	return
}

func (m *monRepo) BlockByHeight(height int) (block Block, err error) {
	c := m.blockCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"_id": height})
	err = query.One(&block)
	return
}

func (m *monRepo) LatestBlocks(limit int) (displayData []DisplayItem, err error) {

	c := m.blockCollection()
	defer c.Database.Session.Close()

	blocks := make([]Block, 0, limit)
	if err = c.Find(nil).Sort("-_id").Limit(limit).All(&blocks); err != nil {
		return
	}

	return BlocksToDisplayItems(blocks, limit), nil
}

func BlocksToDisplayItems(blocks []Block, limit int) (items []DisplayItem) {

	if len(blocks) == 0 {
		return
	}
	dur := blocks[0].Time.Sub(blocks[len(blocks)-1].Time).Seconds()
	totalTxs := 0
	for _, v := range blocks {
		totalTxs += v.NumTxs
	}
	interval := dur / float64(limit)
	tps := int(float64(totalTxs) / dur)
	items = make([]DisplayItem, 0, len(blocks))
	for _, v := range blocks {
		items = append(items, DisplayItem{v, tps, interval})
	}
	return items
}

func (m *monRepo) BlocksFromTo(from, to int) (blocks []Block, err error) {
	c := m.blockCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"_id": bson.M{"$lte": to, "$gte": from}}).Sort("-_id")
	blocks = []Block{}
	err = query.All(&blocks)
	return
}

func (m *monRepo) OneBlock(hash string) (block Block, txs []Transaction, err error) {
	c := m.blockCollection()
	defer c.Database.Session.Close()
	query := c.Find(bson.M{"hash": hash})
	err = query.One(&block)
	if err != nil {
		return
	}
	c2 := m.txCollection()
	defer c2.Database.Session.Close()
	query2 := c2.Find(bson.M{"block": hash})
	err = query2.All(&txs)
	if err != nil {
		return
	}
	return
}

type DBMutex struct {
	Id           string `bson:"_id"`
	Owner        string `bson:"owner"`
	LastLockedAt int64  `bson:"last_locked_at"`
}

func (m *monRepo) TryLock(owner string, longInSec int) (bool, error) {

	_id := "sync_mutex"
	c := m.Session.Clone().DB(DB_NAME).C(MutexCollect)
	defer c.Database.Session.Close()
	now := time.Now().UnixNano()
	lockAt := now + int64(longInSec)*int64(time.Second)
	cond := bson.M{"_id": _id}

	old := DBMutex{}
	if err := c.FindId(_id).One(&old); err != nil {
		if err == mgo.ErrNotFound {
			old.Id = _id
			old.Owner = owner
			old.LastLockedAt = lockAt
			if err := c.Insert(&old); err != nil {
				return false, nil
			}
			return true, nil
		}
		return false, err
	}

	cond["$or"] = []bson.M{{"last_locked_at": bson.M{"$lt": now}}, {"owner": owner}}
	update := bson.M{"$set": bson.M{"owner": owner, "last_locked_at": lockAt}}
	n := DBMutex{}
	chg := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	info, err := c.Find(cond).Apply(chg, &n)
	if err != nil {
		return false, err
	}
	return info.Updated > 0, nil
}
