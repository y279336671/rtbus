package handler

import (
	// "fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func BusLineHandler(params martini.Params, r render.Render) {
	cityid := params["city"]
	lineno := params["linenum"]

	busline, err := BusTool.GetBusLineInfo(cityid, lineno)
	if err != nil {
		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	r.JSON(200,
		&Response{
			0,
			"OK",
			busline,
		},
	)
}

func BusDirHandler(params martini.Params, r render.Render) {
	cityid := params["city"]
	lineno := params["linenum"]
	dirid := params["direction"]

	//方向
	bdi, err := BusTool.GetBusLineDirInfo(cityid, lineno, dirid)
	if err != nil {
		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	r.JSON(200,
		&Response{
			0,
			"OK",
			bdi,
		},
	)
}

func RunningBusHandler(params martini.Params, r render.Render) {
	cityid := params["city"]
	lineno := params["linenum"]
	dirid := params["direction"]

	rbus, err := BusTool.GetRT(cityid, lineno, dirid)
	if err != nil {
		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	r.JSON(200,
		&Response{
			0,
			"OK",
			rbus,
		},
	)
}
