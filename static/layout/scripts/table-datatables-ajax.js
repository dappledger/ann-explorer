$.extend( $.fn.dataTable.defaults, {
    searching: false,
    ordering:  false
});

var BlockDashboardInit = function() {
    var renderBlockList = function(blocks) {
        var blockTrList = []
        for (var i=0; i<blocks.length; i++) {
            var trItem = [
                '<td><a href="/blocks/'+blocks[i].Number+'">'+blocks[i].Hash+'</a></td>',
                '<td><a href="/blocks/'+blocks[i].Number+'">'+blocks[i].Number+'</a></td>',
                '<td>'+blocks[i].ChainID+'</td>',
                '<td class="uppercase">'+blocks[i].ValidatorsHash+'</td>',
                '<td>'+blocks[i].NumTxs+'</td>',
                '<td>'+blocks[i].Time+'</td>'
            ].join('');
            blockTrList.push('<tr>'+trItem+'</tr>')
        }
        $('#dashboard-block-table').find('tr').remove();
        $('#dashboard-block-table').append(blockTrList.join(''));
    }
    var flushBlockList = function() {
        $.ajax({
            url: '/view/block/latest',
            type: 'GET',
            contentType: "application/json; charset=utf-8",
            dataType: "json",
            success: function(result) {
                if (result.Success) {
                    renderBlockList(result.Data)
                    setTimeout(function() {
                        flushBlockList();
                    }, 5000);
                }
            },
            error: function() {

            }
        });
    }

    flushBlockList();
};

jQuery(document).ready(function() {
    BlockListTableInit();
    BlockDashboardInit();
});