# AnnChain BlockChain-Explorer

安链区块链浏览器APP ， 可以查询Block区块信息 、 Transaction交易信息、 Contract合约信息.

## _API_

+ ### 查询指定账户的所有交易
    + _Request URI : /v1/txs/query/:fromTo_
    + _Response :_
    
         ``` 
         {
             "success": true,
             "data": [
                 {
                     "Payload": "", //
                     "Hash": "",    //
                     "From": "",    //
                     "To": "",      //
                     "Receipt": "", //
                     "Amount": 0,   //
                     "Nonce": 1,    //
                     "Gas": null,   //
                     "Size": 296,   //
                     "Block": "",   //
                     "Contract": "" //
                 }
             ]
         }
         ```
    + _Sample :_
    
         ```
         curl --url http://localhost:9090/v1/txs/query/0xb85600baec7119bd2d277d3cb57ec5cccc0750d6

         ---------------------------------------
     
         {
              "success": true,
              "data": [
                  {
                      "Payload": "6UK1FgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAkNOAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAVDaGluYQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
                      "Hash": "0x9c8c1c8e7b4bcb56b524fdc89670c2591d54fe343c858c60d8983deb4befec20",
                      "From": "0x8e26df12a2aeeaad4dea21704ae6bca0cea08ab3",
                      "To": "0xb85600baec7119bd2d277d3cb57ec5cccc0750d6",
                      "Receipt": "",
                      "Amount": 0,
                      "Nonce": 1,
                      "Gas": null,
                      "Size": 296,
                      "Block": "0x6cd76b38859cd9ef458480ec4457afe016ef5b2c",
                      "Contract": ""
                  }
              ]
          }
         ```

### 设置合约元数据
```
POST /contract/meta HTTP/1.1
Content-Type: application/json; charset=utf-8

{"hash":"0xee1bb760aa9d3393edd42a8d0b8e44e2e1cdfb89","abi":"contracts abi"}
```