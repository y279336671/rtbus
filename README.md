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

    //init
    bp, err := api.NewBusPool()
    if err != nil {
        logger.Error("%v", err)
    }

    //query lines
    busLines := [][3]string{
        //[3]string{"北京", "675", "0"},
        [3]string{"北京", "675", "通州李庄-左家庄"},
        //[3]string{"青岛", "318", "市政府-虎山军体中心"},
        [3]string{"青岛", "318", "1"},
    }

    for _, line := range busLines {
        logger.Info("Query %s %s %s ...", line[0], line[1], line[2])

        //线路-各公交站
        bss, err := bp.GetStations(line[0], line[1], line[2])
        if err != nil {
            logger.Error("%v", err)
        }
        for _, bs := range bss {
            logger.Info("%+v", bs)
        }

        //线路-到站情况
        rbuses, err := bp.GetRT(line[0], line[1], line[2])
        if err != nil {
            logger.Error("%v", err)
        }
        for _, rbus := range rbuses {
            logger.Info("%+v", rbus)
        }

        logger.Info("Query %s %s %s over!", line[0], line[1], line[2])
    }
}

```
