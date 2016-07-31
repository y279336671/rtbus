function renderLineInfo(){
    var linenum = $.getUrlParam('linenum');
    var uid = $.getUrlParam('uid');
    var dirid = $.getUrlParam('dirid');
    var sid = $.getUrlParam('sid');
    // console.log(linenum);

    //参数错误
    if(linenum == "" || dirid == ""){
        return
    }

    //渲染
    $('#loadingToast').show();
    $.ajax({
        type:"GET",
        url:"http://api.bingbaba.com/rtbus/bj/info/"+linenum+"/"+dirid,
        // url:"http://127.0.0.1:1315/rtbus/bj/info/"+linenum+"/"+dirid,
        contentType:"application/x-www-form-urlencoded; charset=utf-8",
        success:function(data){
            businfo = data.data;

            $("#cd-timeline").empty();
            for (var i=0;i<businfo.length;i++) {
                station = businfo[i]
                
                var divid = "station_"+station.id;
                var div = "<div id=\""+divid+
"\" class=\"cd-timeline-block\">\
    <div class=\"cd-timeline-img\">\
        <img src=\"vendor/images/cd-icon-location.svg\" alt=\"Picture\">\
    </div>\
    <div class=\"cd-timeline-content\">\
        <h2></h2>\
        <span class=\"cd-date\">未到站</span>\
    </div>\
</div>";
                $("#cd-timeline").append(div);

                if(sid > 0 && sid === station.id){
                    $("#"+divid).find("h2").html(station.name+"(本站)");
                    $("#station_"+sid).addClass("cd-mylocation");
                }else{
                    $("#"+divid).find("h2").html(station.name);
                }
                // console.log($("#"+divid));

                if (station.status != ""){
                    $("#"+divid).addClass("cd-bus");
                    $("#"+divid).find("img").attr("src","vendor/images/bus2.png");

                    if(station.status == "0.5"){
                        $("#"+divid).find("span").html("即将到站");
                    }else if(station.status == "1"){
                        $("#"+divid).find("span").html("到站");
                    }
                }else {
                    // $("#"+divid).find(".cd-timeline-content").hide();
                }
                // console.log($("#"+divid));
                // console.log($("#cd-timeline"));
            }

            //锚点跳到响应位置
            if(sid > 0){
                var t = $("#container").attr("scrollTop");
                if(t <= 1){
                    t = $("#cd-timeline").find("#station_"+(sid-1)).offset().top;
                }
                // console.log(t);
                $("#container").scrollTop(t);
            }
            $('#loadingToast').hide();


            //继续刷新
            setTimeout(renderLineInfo,10100);
        }
    })
}
