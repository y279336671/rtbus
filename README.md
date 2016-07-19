# 北京实时公交
```golang
package main

import (
    "github.com/bingbaba/util/logs"
    "github.com/xuebing1110/rtbus/api"
)

func main() {
    //logger
    logs.SetDebug(true)
    logger := logs.GetBlogger()

    //new
    bus, err := api.NewBJBusSess()
    if err != nil {
        logger.Error("%v", err)
        return
    }

    //如直接查看到站情况，此方法的调用可省略，可直接调用 bus.FreshStatus
    err = bus.LoadBusLineConf("300快内")
    if err != nil {
        logger.Error("%v", err)
        return
    }

    //查看<300快内>路公交<大钟寺-大钟寺>方向的公交实时到站情况
    err = bus.FreshStatus("300快内", "大钟寺-大钟寺")
    if err != nil {
        logger.Error("%v", err)
        return
    }

    //Debug Print
    bus.Print()
}
```