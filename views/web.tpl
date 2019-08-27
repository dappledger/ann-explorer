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
                    <span class="caption-subject bold uppercase">ChainID: </span>
                    <span class="caption-helper" id="chain-id"></span>
                </div>
                <div class="actions">

                    <a class="btn btn-circle btn-icon-only btn-default fullscreen" href="#"> </a>
                </div>
            </div>
            <div class="portlet-body">
                <div class="table-scrollable table-scrollable-borderless">
                    <table class="table table-hover table-light" >
                        <thead>
                            <tr class="uppercase">
                                <th> Block Hash</th>
                                  <th> Height </th>
                                <th> Validator </th>
                                <th> Txs </th>
                                <th> Time </th>
                                <th> Tps(-25) </th>
                                <th> Interval(-25) </th>
                            </tr>
                        </thead>
                        <tbody id="dashboard-block-table"></tbody>
                    </table>
                </div>
                <div>
                    <nav aria-label="Page navigation example">
                        <ul class="pagination" id="blocks-page">
                        </ul>
                    </nav>
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

    var BlockDashboardInit = function() {


        var renderBlockList = function(blocks) {
            var chainID = ''
            if (blocks.length > 0){                                                                                                                chainID = blocks[0].ChainID
            }

            var blockTrList = []
            for (var i=0; i<blocks.length; i++) {
                var trItem = [
                    '<td title='+blocks[i].Hash +'><a href="/view/blocks/hash/'+blocks[i].Hash+'">'+blocks[i].Hash+'</a></td>',
                    '<td>'+blocks[i].Height+'</td>',
                    '<td title='+blocks[i].ValidatorsHash+'>'+blocks[i].ValidatorsHash+'</td>',
                    '<td>'+blocks[i].NumTxs+'</td>',
                    '<td>'+blocks[i].Time+'</td>',
                    '<td>'+blocks[i].Tps+'</td>',
                    '<td>'+blocks[i].Interval.toFixed(3)+'</td>'
                ].join('');
                blockTrList.push('<tr>'+trItem+'</tr>')
            }
            $('#dashboard-block-table').find('tr').remove();
            $('#dashboard-block-table').append(blockTrList.join(''));
            $('#chain-id').text(chainID);
        }
        var flushBlockList = function() {
            $.ajax({
                url: '/view/blocks?p='+page(),
                type: 'GET',
                contentType: "application/json; charset=utf-8",
                dataType: "json",
                success: function(result) {
                    console.log(result)
                    if (result.success) {
                        renderBlockList(result.data.data);
                        renderPager(result.data.page);
                        var p = page();
                        if(p == "latest"){
                            setTimeout(function() {
                                flushBlockList();
                            }, 5000);
                        }else{

                        }
                    }
                }
            });
        }

        flushBlockList();
    };

    $(function() {
        BlockDashboardInit();
    });
</script>
