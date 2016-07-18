package main

import (
	"flag"
	"github.com/astaxie/beego/logs"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xuebing1110/rtbus"
	"io/ioutil"
)

var (
	logger *logs.BeeLogger
	DEBUG  = flag.Bool("debug", false, "the debug module")
)

func main() {
	flag.Parse()
	err := InitLog()
	if err != nil {
		panic(err)
	}

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/bjbus/direction/:linenum", LineNumHandler)

	m.RunOnAddr(":1315")
}

func InitLog() error {
	log_conf := "./log.json"
	bytes_conf, err := ioutil.ReadFile(log_conf)
	if err != nil {
		return err
	}

	logger = logs.NewLogger(1)
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(2)
	logger.SetLogger("file", string(bytes_conf))

	if *DEBUG {
		logger.SetLevel(logs.LevelDebug)
	} else {
		logger.SetLevel(logs.LevelInfo)
		logger.DelLogger("console")
	}
	logger.Info("startint...")

	return nil
}

type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}

func LineNumHandler(params martini.Params, r render.Render) {
	bus, err := rtbus.NewBJBusSess()
	if err != nil {
		logger.Error("%v", err)

		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	err = bus.LoadBusLineConf(params["linenum"])
	if err != nil {
		logger.Error("%v", err)

		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	busline := bus.BusLines[params["linenum"]]
	r.JSON(200,
		&Response{
			0,
			"OK",
			busline,
		},
	)
}
