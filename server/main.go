package main

import (
	"flag"
	"github.com/bingbaba/util/logs"
	"github.com/xuebing1110/rtbus/server/handler"
)

var (
	DEBUG = flag.Bool("debug", false, "the debug module")
)

func main() {
	flag.Parse()

	//init logger
	logs.DEBUG = *DEBUG
	err := logs.Init("log.json")
	if err != nil {
		panic(err)
	}

	handler.Run()
}
