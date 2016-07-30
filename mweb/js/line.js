function renderLineInfo(){
    var linenum = $.getUrlParam('linenum');
    var uid = $.getUrlParam('uid');
    var dirid = $.getUrlParam('dirid');
    // console.log(linenum);

    $('#loadingToast').show();
    $.ajax({
        type:"GET",
        url:"http://api.bingbaba.com/rtbus/bj/info/"+linenum+"/"+dirid,
        // url:"http://127.0.0.1:1315/rtbus/bj/info/"+linenum+"/"+dirid,
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

                $("#"+divid).find("h2").html(station.name);
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

            $('#loadingToast').hide();
        }
    })

}
