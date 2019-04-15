
var ReadyDashboard = function() {
    return {
        init: function(dataHashs,dataDays) {
            /* With CountTo, Check out examples and documentation at https://github.com/mhuggins/jquery-countTo */
            refreshData();
            var widgetChartLineOptions = {
                type: 'line',
                width: '200px',
                height: '109px',
                tooltipOffsetX: -25,
                tooltipOffsetY: 20,
                lineColor: '#9bdfe9',
                fillColor: '#9bdfe9',
                spotColor: '#555555',
                minSpotColor: '#555555',
                maxSpotColor: '#555555',
                highlightSpotColor: '#555555',
                highlightLineColor: '#555555',
                spotRadius: 3,
                tooltipPrefix: '',
                tooltipSuffix: ' Sales',
                tooltipFormat: '{{prefix}}{{y}}{{suffix}}'
            };
            $('#widget-dashchart-sales').sparkline('html', widgetChartLineOptions);

        
            // Get the element where we will attach the chart
            var chartClassicDash    = $('#chart-classic-dash');
            // Classic Chart
            $.plot(chartClassicDash,
                [
                    {
                       
                        data: dataHashs,
                        lines: {show: true, fill: true, fillColor: {colors: [{opacity: .6}, {opacity: .6}]}},
                        points: {show: true, radius: 5}
                    }
                ],
                {
                    colors: ['#5ccdde', '#454e59', '#ffffff'],
                    legend: {show: true, position: 'nw', backgroundOpacity: 0},
                    grid: {borderWidth: 0, hoverable: true, clickable: true},
                    yaxis: {show: true, tickColor: '#f5f5f5', ticks: 3},
                    xaxis: {ticks: dataDays, tickColor: '#f9f9f9'}
                }
            );

            var previousPoint = null, ttlabel = null;
            chartClassicDash.bind('plothover', function(event, pos, item) {
                if (item) {
                    if (previousPoint !== item.dataIndex) {
                        previousPoint = item.dataIndex;
                        $('#chart-tooltip').remove();
                        var y = item.datapoint[1];
                        ttlabel = '<strong>' + y + '</strong>';            
                        $('<div id="chart-tooltip" class="chart-tooltip">' + ttlabel + '</div>')
                            .css({top: item.pageY - 45, left: item.pageX + 5}).appendTo("body").show();
                    }
                }
                else {
                    $('#chart-tooltip').remove();
                    previousPoint = null;
                }
            });
        }
    };
}();
var refreshData = function(){
    $('[data-toggle="counter"]').each(function(){
        $(this).countTo({
            from:parseInt($(this).attr("data-from")),
            to:parseInt($(this).attr("data-to")),
            speed: 1000,
            refreshInterval: 25
        });
    });
}

var App = function() {
    var interval;
    var Init = function() {
        $(".refresh").on("click",null,function(){
            AutoRefresh();
        });
        interval = setInterval(AutoRefresh,4000);
    };
    var AutoRefresh = function () {
        $.ajax({
            type: 'GET',
            url: "/api/getNews",
            dataType: "json",
            success: function(result){
                updateResult(result);
                refreshData();
            }
        });
    }
    var updateResult = function (result) {
        $("#total_hashs").attr("data-to")
        moveValue("total_hashs",result["total_hashs"])
        moveValue("yesterday_hashs",result["yesterday_hashs"])
        moveValue("today_hashs",result["today_hashs"])
        moveValue("spiders",result["spiders"])
        $("#new_hashs").empty()
        result["new_hash"].forEach(function (v,k) {
            var item = '<a href="javascript:void(0)" class="widget-content themed-background-muted text-right clearfix">'+
                '<span class="widget-heading h4 text-muted">'+ v['Hash']+'</span>'+
            '<span class="pull-right text-muted">'+v['CreateTime'] +'</span></a>'
            $("#new_hashs").append(item)
        })
    }
    var moveValue = function (id,value) {
        var node = $("#"+id);
        var from =  node.attr("data-to");
        node.attr("data-from",from).attr("data-to",value)
    }
    return {
        init: function() {
            Init(); 
        }
    };
}();

/* Initialize App when page loads */
$(function(){ App.init();});