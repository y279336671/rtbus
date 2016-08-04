package handler

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func CllBusLineHandler(params martini.Params, r render.Render) {
	citytel := params["city"]
	lineid := params["linenum"]

	fmt.Println(citytel)

	cllbus, err := DCllBusPool.GetCllBus(citytel)
	if err != nil {
		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	busline, err2 := cllbus.GetBusLine(lineid)
	if err2 != nil {
		r.JSON(
			502,
			&Response{502, err2.Error(), nil},
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

func CllBusSnHandler(params martini.Params, r render.Render) {
	citytel := params["city"]
	lineid := params["linenum"]
	dirid := params["direction"]

	//城市
	cllbus, err := DCllBusPool.GetCllBus(citytel)
	if err != nil {
		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	//路线
	busline, err2 := cllbus.GetBusLine(lineid)
	if err2 != nil {
		r.JSON(
			502,
			&Response{502, err2.Error(), nil},
		)
		return
	}

	//方向
	busdir, err3 := busline.GetBusDir(dirid, cllbus)
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
}

func CllRunningBusHandler(params martini.Params, r render.Render) {
	citytel := params["city"]
	lineid := params["linenum"]
	dirid := params["direction"]

	//城市
	cllbus, err := DCllBusPool.GetCllBus(citytel)
	if err != nil {
		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	//路线
	busline, err2 := cllbus.GetBusLine(lineid)
	if err2 != nil {
		r.JSON(
			502,
			&Response{502, err2.Error(), nil},
		)
		return
	}

	//方向
	rbus, err3 := busline.GetRunningBus(dirid, cllbus)
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
}
