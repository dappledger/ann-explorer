<style>
    td {
          white-space:nowrap;
          overflow:hidden;
          text-overflow: ellipsis;
    }
    .pure-table {
        width: 100%;
    }
    .txwrapper{
        min-height:200px;
    }
</style>

<script>

function showDecodedData(abi, data, id){
    console.log(abi,data)
    $("#txpayloadHex_" + id).toggle();
    $("#txpayloadHtml_" + id).toggle();
    var tab = decodeFunctionArgsToTable(JSON.parse(abi), data);
    $("#txpayloadHtml_"+id).html(tab);
    if($("#txDecodeBtn_"+id).text()=="decode"){
       $("#txDecodeBtn_"+id).text('back');
    }else{
       $("#txDecodeBtn_"+id).text('decode');
    }
}

</script>

<div class="row">
    <div class="col-md-12">
        <!-- Begin: life time stats -->
        <div class="portlet light portlet-fit portlet-datatable ">
            <div class="portlet-title">
                <div class="caption">
                    <i class="fa fa-book font-green"></i>
                    <span class="caption-subject font-green sbold">Hash : {{.Block.Hash}}</span>
                </div>
                <div class="actions">
                </div>
            </div>
            <div class="portlet-body">
                <div class="table-container">
                    <table class="table table-striped table-bordered table-hover">
                        <tbody>
                          <tr><td>Chain ID:</td><td>{{.Block.ChainID}}</td></tr>
                          <tr><td>Height:</td><td>{{.Block.Height}}</td></tr>
                          <tr><td>NumTxs:</td><td>{{.Block.NumTxs}}</td></tr>
                          <tr><td>ProposerAddress:</td><td>{{.Block.ProposerAddress}}</td></tr>
                          <tr><td>AppHash:</td><td>{{.Block.AppHash}}</td></tr>
                          <tr><td>Time:</td><td>{{.Block.Time}}</td></tr>
                        </tbody>
                    </table>
                </div>

                 <div class="portlet-title">
                    <div class="caption">
                        <i class="fa fa-book font-green"></i>
                        <span class="caption-subject font-green sbold">Transactions : </span>
                    </div>
                    <div class="actions">
                    </div>
                </div>
                <div class="portlet-body">
                    <div class="table-container">
                        <table class="table table-striped table-bordered table-hover">
                            <tbody>
                            {{range $index,$tx := .Transactions}}
                                  <tr><td>Hash:</td><td>{{$tx.Hash}}</td></tr>
                                  <tr><td>From:</td><td>{{$tx.From}}</td></tr>
                                  <tr><td>To:</td><td>{{$tx.To}}</td></tr>
                                  <tr>
                                    <td>PayLoad:</td>
                                    <td>
                                        <div class="txwrapper">
                                            <textarea id="txpayloadHex_{{$index}}" style="width:100%;height:200px;" readonly>{{$tx.PayloadHex}}</textarea>
                                            <div id="txpayloadHtml_{{$index}}" style="display:none"></div>
                                        </div>
                                        {{ if $tx.ContractMeta.ABI}}
                                            <div><button id="txDecodeBtn_{{$index}}" class="btn btn-default btn-sm" onclick='showDecodedData({{$tx.ContractMeta.ABI}},{{$tx.PayloadHex}},{{$index}})'>decode</button></div>
                                        {{end}}
                                    </td>
                                  </tr>
                             {{end}}
                            </tbody>
                        </table>
                    </div>
                </div>

            </div>
        </div>
        <!-- End: life time stats -->
    </div>
</div>