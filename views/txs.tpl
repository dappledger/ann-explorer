<style>
    td {
          white-space:nowrap;
          overflow:hidden;
          text-overflow: ellipsis;
    }
</style>

<div class="row">
    <div class="col-md-12">
        <div class="portlet light ">
            <div class="portlet-title">
                <div class="caption font-green">
                    <span class="caption-subject bold uppercase">TRANSACTIONS</span>
                    <span class="caption-helper">Recent transactions...</span>
                </div>
                <div class="actions">
                    <a class="btn btn-circle btn-icon-only btn-default fullscreen" href="#"> </a>
                </div>
            </div>
            <div class="portlet-body">
                <div class="table-scrollable table-scrollable-borderless">
                    <table class="table table-hover table-light" style="table-layout: fixed;">
                        <thead>
                            <tr class="uppercase">
                                <th> Block Hash </th>
                                <th> Tx Hash</th>
                                <th> From </th>
                                <th> To </th>
                                <th> Nonce </th>
                            </tr>
                        </thead>
                        <tbody id="dashboard-block-table"></tbody>
                    </table>
                </div>

            </div>
        </div>
    </div>
</div>

<script>
    $.extend( $.fn.dataTable.defaults, {
        searching: false,
        ordering:  false
    });

    var TxDashboardInit = function() {
        var renderTxList = function(txs) {
            var txsTrList = []
            for (var i=0; i<txs.length; i++) {
                var trItem = [
                    '<td><a href="/view/blocks/hash/'+txs[i].Block+'">'+txs[i].Block+'</a></td>',
                    '<td>'+txs[i].Hash+'</td>',
                    '<td>'+txs[i].From+'</td>',
                    '<td>'+txs[i].Contract+'</td>',
                    '<td>'+txs[i].Nonce+'</td>'
                ].join('');
                txsTrList.push('<tr>'+trItem+'</tr>')
            }
            $('#dashboard-block-table').find('tr').remove();
            $('#dashboard-block-table').append(txsTrList.join(''));
        }
        var flushTxList = function() {
            $.ajax({
                url: '/view/txs/latest',
                type: 'GET',
                contentType: "application/json; charset=utf-8",
                dataType: "json",
                success: function(result) {
                    console.log(result)
                    if (result.success) {
                        renderTxList(result.data)
                        setTimeout(function() {
                            flushTxList();
                        }, 5000);
                    }
                }
            });
        }

        flushTxList();
    };

    $(function() {
        TxDashboardInit();
    });
</script>