<div class="row">
    <div class="col-md-12">
        <!-- Begin: life time stats -->
        <div class="portlet light portlet-fit portlet-datatable ">
            <div class="portlet-title">
                <div class="caption">
                    <i class="fa fa-book font-green"></i>
                    <span class="caption-subject font-green sbold">Block Height: {{.Block.Number}} ({{.Block.Time}})</span>
                </div>
                <div class="actions">
                </div>
            </div>
            <div class="portlet-body">
                <div class="table-container">
                    <table class="table table-striped table-bordered table-hover">
                        <tbody>
                          <tr><td>Chain ID:</td><td>{{.Block.ChainID}}</td></tr>
                          <tr><td>Height:</td><td>{{.Block.Number}}</td></tr>
                          <tr><td>Hash:</td><td class="uppercase">{{.Block.Hash}}</td></tr>
                          <tr><td>Time:</td><td>{{.Block.Time}}</td></tr>
                          <tr><td>NumTxs:</td><td>{{.Block.NumTxs}}</td></tr>
                          <tr><td>LastBlockHash:</td><td class="uppercase">{{.Block.LastBlockHash}}</td></tr>
                          <tr><td>LastCommitHash:</td><td class="uppercase">{{.Block.LastCommitHash}}</td></tr>
                          <tr><td>ProposerAddress:</td><td class="uppercase">{{.Block.ProposerAddress}}</td></tr>
                          <tr><td>AppHash:</td><td class="uppercase">{{.Block.AppHash}}</td></tr>
                          <tr><td>Txs:</td><td class="uppercase">{{.Block.Txs}}</td></tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
        <!-- End: life time stats -->
    </div>
</div>