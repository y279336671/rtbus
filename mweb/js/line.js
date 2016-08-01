var times = 0;
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
        requrl = "http://api.bingbaba.com/rtbus/bj/info/"+linenum+"/"+dirid;
    }else {
        requrl = "http://api.bingbaba.com/rtbus/bj/info/"+linenum+"/"+dirid+"?simple=1";
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

function updateTimelineContainer(businfo,sid){
    for (var i=0;i<businfo.length;i++) {
        station = businfo[i];

        //初始化
        var divid = "station_"+station.id;
        $("#"+divid).removeClass("cd-mylocation");
        $("#"+divid).removeClass("cd-bus");
        $("#"+divid).find("span").remove();
        $("#"+divid).find("img").attr("src","vendor/images/cd-icon-location.svg");

        //未到站 站点
        if(station.status <= 0){
            if(sid > 0 && sid === station.id){
                $("#station_"+sid).addClass("cd-mylocation");
            }
        }else { //到站 OR 即将到站
            if(sid > 0 && sid === station.id){
                $("#"+divid).addClass("cd-mylocation");
            }else {
                $("#"+divid).addClass("cd-bus");
                $("#"+divid).find("img").attr("src","vendor/images/bus2.png");
            }

            //到站
            if(station.status == "1") {
                $("#"+divid).find("h2").after("<span class=\"cd-date\">到站</span>");
            }else if(station.status == "0.5"){ //即将到站
                $("#"+divid).find("h2").after("<span class=\"cd-date\">即将到站...</span>");
            }
        }
    }
}

function initTimelineContainer(businfo,sid){
    $("#cd-timeline").empty();
    for (var i=0;i<businfo.length;i++) {
        station = businfo[i];
        
        var divid = "station_"+station.id;
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


        var div = divh+divf;
        $("#cd-timeline").append(div);
        $("#"+divid).find("h2").html(station.name);

        //未到站 站点
        if(station.status <= 0){
            if(sid > 0 && sid === station.id){
                $("#station_"+sid).addClass("cd-mylocation");
            }
        }else { //到站 OR 即将到站
            if(sid > 0 && sid === station.id){
                $("#"+divid).addClass("cd-mylocation");
            }else {
                $("#"+divid).addClass("cd-bus");
                $("#"+divid).find("img").attr("src","vendor/images/bus2.png");
            }

            //到站
            if(station.status == "1") {
                $("#"+divid).find("h2").after("<span class=\"cd-date\">到站</span>");
            }else if(station.status == "0.5"){ //即将到站
                $("#"+divid).find("h2").after("<span class=\"cd-date\">即将到站...</span>");
            }
        }

        // console.log($("#"+divid));
        // console.log($("#cd-timeline"));
    }
}