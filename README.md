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

    //北京实时公交
    bjbus, err := api.NewBJBusSess()
    if err != nil {
        logger.Error("%v", err)
        return
    }

    //查看到站情况
    _, err = bjbus.GetBusLine("300快内")
    if err != nil {
        logger.Error("%v", err)
        return
    }

    //Debug Print
    bjbus.Print()

    //青岛实时公交
    cllbus, err := api.NewCllBus("0532")
    if err != nil {
        logger.Error("%v", err)
        return
    }

    _, err = cllbus.GetBusLine("318")
    if err != nil {
        logger.Error("%v", err)
        return
    }

}
```

# 其他城市实时公交（青岛）
```golang
    //青岛实时公交
    cllbus, err := api.NewCllBus("0532")
    if err != nil {
        logger.Error("%v", err)
        return
    }

    _, err = cllbus.GetBusLine("318")
    if err != nil {
        logger.Error("%v", err)
        return
    }

```