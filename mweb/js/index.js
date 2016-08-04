$(function () {

    var router = new Router({
        container: '#container',
        enterTimeout: 250,
        leaveTimeout: 250
    });

    // add
    var home = {
        url: '/',
        className: 'home',
        render: function () {
            $('#container').load('html/rtbus.html');
        },
        bind: function () {
            $.getScript("js/rtbus.js",function(){
                $('.container').on('change', '#rtbus_city', showstation);
                $('.container').on('change', '#rtbus_direction', showstation);
                $('.container').on('blur', '#busline', getbusline);
            });
        }
    };

    //line
    var line = {
        url:'/line/:cityid/:lineid/:dirid/:sid',
        className:'line',
        render: function () {
            $('#container').load('html/line.html');

            var cityid = this.params.cityid;
            var lineid = this.params.lineid;
            var dirid = this.params.dirid;
            var sid = this.params.sid;
            $.getScript("js/line.js",function(){
                renderLineInfoByParams(cityid,lineid,dirid,sid);
                
                $('.container').on('click', '.cd-fav-div', changefavorite);
            });
        },
        bind: function () {}
    }

    router.push(home)
        .push(line)
        .setDefault('/')
        .init();


    // .container 设置了 overflow 属性, 导致 Android 手机下输入框获取焦点时, 输入法挡住输入框的 bug
    // 相关 issue: https://github.com/weui/weui/issues/15
    // 解决方法:
    // 0. .container 去掉 overflow 属性, 但此 demo 下会引发别的问题
    // 1. 参考 http://stackoverflow.com/questions/23757345/android-does-not-correctly-scroll-on-input-focus-if-not-body-element
    //    Android 手机下, input 或 textarea 元素聚焦时, 主动滚一把
    if (/Android/gi.test(navigator.userAgent)) {
        window.addEventListener('resize', function () {
            if (document.activeElement.tagName == 'INPUT' || document.activeElement.tagName == 'TEXTAREA') {
                window.setTimeout(function () {
                    document.activeElement.scrollIntoViewIfNeeded();
                }, 0);
            }
        })
    }

    $.getUrlParam = function(name) {  
        var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");  
        var r = window.location.search.substr(1).match(reg);  
        if (r!=null) return decodeURI(r[2]); return null;  
    }
});
