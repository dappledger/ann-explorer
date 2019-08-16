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
                    <span class="caption-subject font-green sbold">Contract : {{.Contract.Contract}}</span>
                </div>
                <div class="actions">
                </div>
            </div>
            <div class="portlet-body">
                <div class="table-container">
                    <table class="table table-striped table-bordered table-hover">
                        <tbody>
                          <tr><td>Contract Hash:</td><td>{{.Contract.Hash}}</td></tr>
                          <tr><td>Block :</td><td><a href="/view/blocks/hash/{{.Contract.Block}}">{{.Contract.Block}}</td></tr>
                          <tr><td>From:</td><td>{{.Contract.From}}</td></tr>
                          <tr><td>Nonce:</td><td>{{.Contract.Nonce}}</td></tr>
                          <tr><td>Size:</td><td>{{.Contract.Size}}</td></tr>
                          <tr><td>Contract Code:</td><td><textarea style="width:100%;height:200px;" readonly>{{.Contract.Payload}}</textarea></td></tr>
                        </tbody>
                    </table>
                </div>

                 <div class="portlet-title">
                    <div class="caption">
                        <i class="fa fa-book font-green"></i>
                        <span class="caption-subject font-green sbold">Transactions In Contract : </span>
                    </div>
                    <div class="actions">
                    </div>
                </div>
                <div class="portlet-body">
                    <div class="table-container">
                        <table class="table table-striped table-bordered table-hover">
                            <tbody>
                            {{$contractMeta := .ContractMeta}}
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

                                        {{ if $contractMeta.ABI}}
                                            <div><button id="txDecodeBtn_{{$index}}" class="btn btn-default btn-sm" onclick='showDecodedData({{$contractMeta.ABI}},{{$tx.PayloadHex}},{{$index}})'>decode</button></div>
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