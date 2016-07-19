package main

import (
	"flag"
	"github.com/bingbaba/util/logs"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xuebing1110/rtbus/api"
)

var (
	logger *logs.Blogger
	DEBUG  = flag.Bool("debug", false, "the debug module")
)

func main() {
	flag.Parse()

	//init logger
	err := logs.Init("log.json")
	if err != nil {
		panic(err)
	}
	logger = logs.GetBlogger()

	//martini
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/rtbus/bj/direction/:linenum", LineNumHandler)

	m.RunOnAddr(":1315")
}

type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}

func LineNumHandler(params martini.Params, r render.Render) {
	bus, err := api.NewBJBusSess()
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
