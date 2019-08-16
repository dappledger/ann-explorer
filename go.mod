module github.com/dappledger/AnnChain-browser

go 1.12

require (
	github.com/astaxie/beego v1.11.1
	github.com/dappledger/AnnChain v0.0.0
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190426145343-a29dc8fdc734
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)

replace github.com/dappledger/AnnChain => ../AnnChain
