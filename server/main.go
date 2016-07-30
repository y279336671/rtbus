package main

import (
	"flag"
	"github.com/bingbaba/util/logs"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xuebing1110/rtbus/api"
)

var (
	logger  *logs.Blogger
	DEBUG   = flag.Bool("debug", false, "the debug module")
	BusSess *api.BJBusSess
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
	m.Get("/rtbus/bj/info/:linenum/:direction", LineInfoHandler)

	m.RunOnAddr(":1315")
}

func init() {
	var err error
	BusSess, err = api.NewBJBusSess()
	if err != nil {
		panic(err)
	}
}

type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}

func LineInfoHandler(params martini.Params, r render.Render) {
	if BusSess == nil {
		r.JSON(
			502,
			&Response{502, "bjbus sess token error", nil},
		)
		return
	}

	stations, err := BusSess.GetLineInfo(params["linenum"], params["direction"])
	if err != nil {
		logger.Error("%v", err)

		r.JSON(
			502,
			&Response{503, err.Error(), nil},
		)
		return
	}

	r.JSON(200,
		&Response{
			0,
			"OK",
			stations,
		},
	)

}

func LineNumHandler(params martini.Params, r render.Render) {
	if BusSess == nil {
		r.JSON(
			504,
			&Response{504, "bjbus sess token error", nil},
		)
		return
	}

	err := BusSess.LoadBusLineConf(params["linenum"])
	if err != nil {
		logger.Error("%v", err)

		r.JSON(
			502,
			&Response{505, err.Error(), nil},
		)
		return
	}

	busline := BusSess.BusLines[params["linenum"]]
	r.JSON(200,
		&Response{
			0,
			"OK",
			busline,
		},
	)
}
