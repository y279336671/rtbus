var times = 0;
var domain = "http://api.bingbaba.com";

function renderLineInfo(){
    var linenum = $.getUrlParam('linenum');
    var uid = $.getUrlParam('uid');
    var dirid = $.getUrlParam('dirid');
    var sid = $.getUrlParam('sid');
    times++;
    // console.log(linenum);

    //参数错误
    if(linenum == "" || dirid == ""){
        return
    }

    var requrl;
    if(times === 1){
        requrl = domain+"/rtbus/bj/station/"+linenum+"/"+dirid;
    }else {
        requrl = domain+"/rtbus/bj/bus/"+linenum+"/"+dirid;
    }

    //渲染
    $('#loadingToast').show();
    $.ajax({
        type:"GET",
        url: requrl,
        // url:"http://127.0.0.1:1315/rtbus/bj/info/"+linenum+"/"+dirid,
        contentType:"application/x-www-form-urlencoded; charset=utf-8",
        success:function(data){
            if(times === 1){
                initTimelineContainer(data.data,sid);
            }else {
                updateTimelineContainer(data.data,sid);
            }

            //锚点跳到响应位置
            if(sid > 0){
                var t = $("#container").attr("scrollTop");
                if(t <= 1 && sid > 1){
                    t = $("#cd-timeline").find("#station_"+(sid-1)).offset().top;
                }
                // console.log(t);
                $("#container").scrollTop(t);
            }
            $('#loadingToast').hide();

            //继续刷新
            setTimeout(renderLineInfo,10100);
        },
        error: function(){
            $('#loadingToast').hide();
            setTimeout(renderLineInfo,10100);
        }
    })
}

function initTimelineContainer(businfo,sid){
    $("#cd-timeline").empty();
    for (var i=0;i<businfo.length;i++) {
        station = businfo[i];

        var divid = "station_"+station.order;
        var divh = "<div id=\""+divid;
        var divf= "\" class=\"cd-timeline-block\">\
            <div class=\"cd-timeline-img\">\
                <img src=\"vendor/images/cd-icon-location.svg\">\
            </div>\
            <div class=\"cd-timeline-content\">\
                <h2></h2>\
            </div>\
        </div>";
        // <span class=\"cd-date\">未到站</span>

        //新增div
        var div = divh+divf;
        $("#cd-timeline").append(div);
        $("#"+divid).find("h2").html(station.sn);

        //刷新div样式
        refreshStationDiv(divid,station,sid)
        // console.log($("#"+divid));
        // console.log($("#cd-timeline"));
    }
}

function updateTimelineContainer(businfo,sid){
    for (var i=0;i<businfo.length;i++) {
        station = businfo[i];

        console.log(station);

        //初始化
        var divid = "station_"+station.order;
        $("#"+divid).removeClass("cd-mylocation");
        $("#"+divid).removeClass("cd-bus");
        $("#"+divid).find("span").remove();
        $("#"+divid).find("img").attr("src","vendor/images/cd-icon-location.svg");

        //刷新div样式
        refreshStationDiv(divid,station,sid)
    }
}

function refreshStationDiv(divid,station,sid) {
    //未到站
    if(!station.buses){
        if(sid > 0 && sid == station.order){
            $("#station_"+sid).addClass("cd-mylocation");
        }
    }else { //到站 OR 即将到站
        var nearstation;
        var arrival,warrival = false;
        for (var i = 0; i < station.buses.length; i++) {
            var bus = station.buses[i];

            //到站
            if(bus.status == "1") {
                arrival = true;
            }else if(bus.status == "0.5"){ //即将到站
                warrival = true;
            }
        }

        console.log(sid +"<=>"+station.order);
        if(sid > 0 && sid == station.order){
            console.log("true...")
            $("#"+divid).addClass("cd-mylocation");
        }else {
            $("#"+divid).addClass("cd-bus");
            $("#"+divid).find("img").attr("src","vendor/images/bus2.png");
        }


        //到站
        if(arrival) {
            $("#"+divid).find("h2").after("<span class=\"cd-date\">到站</span>");
        }else if(warrival){ //即将到站
            $("#"+divid).find("h2").after("<span class=\"cd-date\">即将到站...</span>");
        }
    }
}