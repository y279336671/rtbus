var cache = {
    domain:"http://api.bingbaba.com",
    times:0,
    params:{}
};

function renderLineInfoByParams(city,lineid,dirid,sid){
    cache.params.cityid = city;
    cache.params.lineid = decodeURI(lineid);
    cache.params.dirid = dirid;
    cache.params.sid = sid;

    renderLineInfo();
}

function changefavorite(){
    var img = $(this).find("img");
    if(img.hasClass("cd-fav-img")) {
        img.attr("src","vendor/images/unfavorite.svg");
        img.removeClass("cd-fav-img");
        img.addClass("cd-unfav-img");
    }else {
        img.attr("src","vendor/images/favorite.svg");
        img.removeClass("cd-unfav-img");
        img.addClass("cd-fav-img");
    }
}

function renderLineInfo(){
    cache.times++;
    // console.log(lineid);

    var city = cache.params.cityid;
    var lineid = cache.params.lineid;
    var dirid = cache.params.dirid;
    var sid = cache.params.sid;

    //参数错误
    if(lineid == "" || dirid == ""){
        return
    }

    var requrl;
    if(cache.times === 1){
        requrl = cache.domain+"/rtbus/"+city+"/station/"+lineid+"/"+dirid;
    }else {
        requrl = cache.domain+"/rtbus/"+city+"/bus/"+lineid+"/"+dirid;
    }

    //渲染
    $('#loadingToast').show();
    $.ajax({
        type:"GET",
        url: requrl,
        // url:"http://127.0.0.1:1315/rtbus/bj/info/"+lineid+"/"+dirid,
        contentType:"application/x-www-form-urlencoded; charset=utf-8",
        success:function(data){
            if(cache.times === 1){
                initTimelineContainer(data.data);
            }else {
                updateTimelineContainer(data.data);
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

function initTimelineContainer(businfo){
    //修改标题
    var dirname = businfo[0].sn+"-"+businfo[businfo.length-1].sn;
    var title = cache.params.lineid+"路("+dirname+")";
    $(document).find("title").text(title+"实时公交");
    $("#rtbus_title").text(cache.params.lineid+"路实时公交");
    $("#rtbus_pagedesc").text(dirname);

    $("#cd-timeline").empty();
    for (var i=0;i<businfo.length;i++) {
        station = businfo[i];

        var divid = "station_"+station.order;
        var divh = "<div id=\""+divid;
        var divf= "\" class=\"cd-timeline-block\">\
            <div class=\"cd-timeline-img\">\
                <img class=\"cd-icon\" src=\"vendor/images/cd-icon-location.svg\">\
            </div>\
            <div class=\"cd-timeline-content\">\
                <div class=\"cd-fav-div\">\
                    <img class=\"cd-fav-img\" src=\"vendor/images/unfavorite.svg\">\
                </div>\
                <h2></h2>\
            </div>\
        </div>";
        // <span class=\"cd-date\">未到站</span>

        //新增div
        var div = divh+divf;
        $("#cd-timeline").append(div);
        $("#"+divid).find("h2").html(station.sn);

        //刷新div样式
        refreshStationDiv(divid,station)
        // console.log($("#"+divid));
        // console.log($("#cd-timeline"));
    }
}

function updateTimelineContainer(businfo){
    //初始化
    var children = $("#cd-timeline").children("div");
    for (var i = 0; i < children.length; i++) {
        var child = children[i];
        // $(child).removeClass("cd-mylocation");
        $(child).removeClass("cd-bus");
        $(child).find("span").remove();
        $(child).find(".cd-icon").attr("src","vendor/images/cd-icon-location.svg");
    }


    for (var i=0;i<businfo.length;i++) {
        station = businfo[i];
        var divid = "station_"+station.order;
        
        //刷新div样式
        refreshStationDiv(divid,station);
    }
}

function refreshStationDiv(divid,station) {
    var sid = cache.params.sid;

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
            $("#"+divid).find("cd-icon").attr("src","vendor/images/bus.svg");
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