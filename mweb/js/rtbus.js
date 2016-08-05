function putbusline(){
    $('#loadingToast').show();

    setTimeout(function () {
        $('#loadingToast').hide();
        $('#toast').show();
        setTimeout(function () {
            $('#toast').hide();
        }, 1500);
    }, 1000);
}

function gotoline(){
    //城市
    var city = $("#rtbus_city").val();

    //公交线
    var linenum = $("#busline").val();
    if(linenum === ""){
        $("#busline").focus();
        return
    }

    //方向
    var dirid = $("#rtbus_direction").val();
    if(dirid === ""){
        $("#rtbus_direction").focus();
        return
    }

    //公交站
    var sid = $("#rtbus_station").val();
    if(sid === ""){
        $("#rtbus_station").focus();
        return
    }

    location.hash = "#/line/"+city+"/"+ encodeURI(linenum)+"/"+dirid+"/"+sid;
}

//全局变量
var busline,businfo,busdir;

function getbusline(){
    if (busline != $("#busline").val() && $("#busline").val() != "") {
        busline = $("#busline").val();
        var city = $("#rtbus_city").val();

        $('#loadingToast').show();
        $.ajax({
            type:"GET",
            url:"http://api.bingbaba.com/rtbus/"+city+"/direction/"+busline,
            success:function(data){
                businfo = data.data;
                $("#rtbus_direction").empty();
                for (var i=0;i<businfo.direction.length;i++) {
                    busdir = businfo.direction[i];

                    var option;
                    if(i===0){
                        option = "<option selected='selected' value='"+busdir.id+"'>"+busdir.name+"</option>";
                    }else {
                        option = "<option value='"+busdir.id+"'>"+busdir.name+"</option>"
                    }

                    $("#rtbus_direction").append(option);
                }

                showstation();
                $('#loadingToast').hide();
            },
            error:function(){
                $("#busline").empty();
                $("#busline").attr("placeholder",busline+"路公交暂不支持实时查询，请重新输入");
                $("#rtbus_direction").focus();
                $('#loadingToast').hide();
            }
        })
    }
}

function showstation(){
    citycode =  $("#rtbus_city").val();
    cityname = $("#rtbus_city option").eq($("#rtbus_city").attr("selectedIndex")).text();
    direction = $("#rtbus_direction").val();
    busline = $("#busline").val();

    //未输入公交线路
    if(busline === ""){
        return
    }

    $('#loadingToast').show();
    for (var i=0;i<businfo.direction.length;i++) {
        busdir = businfo.direction[i];

        if(busdir.id == direction){
            $("#rtbus_station").empty();

            for (var i=0;i<busdir.stations.length;i++) {
                var station = busdir.stations[i];
                $("#rtbus_station").append("<option value='"+station.order+"'>"+station.sn+"</option>")
            }
            break
        }
    }
    $('#loadingToast').hide();
}