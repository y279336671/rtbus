package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func BJBusLineHandler(params martini.Params, r render.Render) {
	if BjBusSess == nil {
		r.JSON(
			502,
			&Response{502, "bjbus sess token error", nil},
		)
		return
	}

	lineid := params["linenum"]

	busline, err := BjBusSess.GetBusLine(lineid)
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
			busline,
		},
	)

	return
}

func BJBusSnHandler(params martini.Params, r render.Render) {
	if BjBusSess == nil {
		r.JSON(
			502,
			&Response{502, "bjbus sess token error", nil},
		)
		return
	}

	lineid := params["linenum"]
	dirid := params["direction"]

	//路线
	busline, err2 := BjBusSess.GetBusLine(lineid)
	if err2 != nil {
		r.JSON(
			502,
			&Response{502, err2.Error(), nil},
		)
		return
	}

	//方向
	busdir, err3 := busline.GetBusDir(dirid, BjBusSess)
	if err3 != nil {
		r.JSON(
			502,
			&Response{502, err3.Error(), nil},
		)
		return
	}

	r.JSON(200,
		&Response{
			0,
			"OK",
			busdir.Stations,
		},
	)

	return
}

func BJRunningBusHandler(params martini.Params, r render.Render) {
	if BjBusSess == nil {
		r.JSON(
			502,
			&Response{502, "bjbus sess token error", nil},
		)
		return
	}

	lineid := params["linenum"]
	dirid := params["direction"]

	//路线
	busline, err2 := BjBusSess.GetBusLine(lineid)
	if err2 != nil {
		r.JSON(
			502,
			&Response{502, err2.Error(), nil},
		)
		return
	}

	//方向
	rbus, err3 := busline.GetRunningBus(dirid, BjBusSess)
	if err3 != nil {
		r.JSON(
			502,
			&Response{502, err3.Error(), nil},
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
	return
}
