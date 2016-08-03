function putbusline(){
    $('#loadingToast').show();

    setTimeout(function () {
        $('#loadingToast').hide();
        $('#toast').show();
        setTimeout(function () {
            $('#toast').hide();
        }, 1500);
    }, 1000);
    console.log($("#rtbus_direction").val());
}

function gotoline(){
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
    location.href = "index.html?linenum="+linenum+"&dirid="+dirid+"&sid="+sid+"#/line";
}

//全局变量
var busline,businfo,busdir;

function getbusline(){
    if (busline != $("#busline").val() && $("#busline").val() != "") {
        busline = $("#busline").val();

        $('#loadingToast').show();
        $.ajax({
            type:"GET",
            url:"http://api.bingbaba.com/rtbus/bj/direction/"+busline,
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

                    console.log(busdir)
                    $("#rtbus_direction").append(option);
                }

                showstation();
                $('#loadingToast').hide();
            }
        })
    }
}

function showstation(){
    direction = $("#rtbus_direction").val()
    console.log(direction)

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