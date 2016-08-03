var times = 0;
var domain = "http://api.bingbaba.com";

function renderLineInfo(){
    var linenum = $.getUrlParam('linenum');
    var uid = $.getUrlParam('uid');
    var dirid = $.getUrlParam('dirid');
    var sid = $.getUrlParam('sid');
    var city = $.getUrlParam('city');
    times++;
    // console.log(linenum);

    //参数错误
    if(linenum == "" || dirid == ""){
        return
    }

    var requrl;
    if(times === 1){
        requrl = domain+"/rtbus/"+city+"/station/"+linenum+"/"+dirid;
    }else {
        requrl = domain+"/rtbus/"+city+"/bus/"+linenum+"/"+dirid;
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
    //初始化
    var children = $("#cd-timeline").children("div");
    for (var i = 0; i < children.length; i++) {
        var child = children[i];
        $(child).removeClass("cd-mylocation");
        $(child).removeClass("cd-bus");
        $(child).find("span").remove();
        $(child).find("img").attr("src","vendor/images/cd-icon-location.svg");
    }


    for (var i=0;i<businfo.length;i++) {
        station = businfo[i];
        var divid = "station_"+station.order;
        
        //刷新div样式
        refreshStationDiv(divid,station,sid);
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

            //最近公交
            if(!nearstation || bus.distanceToSc < nearstation.distanceToSc){
                nearstation = bus
            }
        }

        if(sid > 0 && sid == station.order){
            $("#"+divid).addClass("cd-mylocation");
        }else {
            $("#"+divid).addClass("cd-bus");
            $("#"+divid).find("img").attr("src","vendor/images/bus2.png");
        }

        //到站
        var text = "";
        if(arrival) {
            text = "到站";
        }else if(warrival){ //即将到站
            if(nearstation.distanceToSc && nearstation.distanceToSc > 0){
                var intreval = nearstation.syncTime;
                text = "即将到站,还有"+nearstation.distanceToSc+"米";
            }else {
                text = "即将到站...";
            }
        }

        if(nearstation.syncTime && nearstation.syncTime > 0){
            var intreval = nearstation.syncTime;
            text = "["+intreval+"秒前]"+text;
        }else {
            text = text+"...";
        }
        $("#"+divid).find("h2").after("<span class=\"cd-date\">"+text+"</span>");
    }
}