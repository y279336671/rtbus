package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func BJBusSnHandler(params martini.Params, r render.Render) {
	if BjBusSess == nil {
		r.JSON(
			502,
			&Response{502, "bjbus sess token error", nil},
		)
		return
	}

	stations, err := BjBusSess.GetLineStationInfo(params["linenum"], params["direction"])
	if err != nil {
		logger.Error("%v", err)
		r.JSON(
			502,
			&Response{503, err.Error(), nil},
		)
	} else {
		r.JSON(200,
			&Response{
				0,
				"OK",
				stations,
			},
		)
	}

	return
}

func BJBusSnBusHandler(params martini.Params, r render.Render) {
	if BjBusSess == nil {
		r.JSON(
			502,
			&Response{502, "bjbus sess token error", nil},
		)
		return
	}

	buses, err := BjBusSess.GetLineBusInfo(params["linenum"], params["direction"])
	if err != nil {
		logger.Error("%v", err)
		r.JSON(
			502,
			&Response{503, err.Error(), nil},
		)
	} else {
		r.JSON(200,
			&Response{
				0,
				"OK",
				buses,
			},
		)
	}

	return
}

func BJBusLineHandler(params martini.Params, r render.Render) {
	if BjBusSess == nil {
		r.JSON(
			504,
			&Response{504, "bjbus sess token error", nil},
		)
		return
	}

	err := BjBusSess.LoadBusLineConf(params["linenum"])
	if err != nil {
		logger.Error("%v", err)

		r.JSON(
			502,
			&Response{505, err.Error(), nil},
		)
		return
	}

	busline := BjBusSess.BusLines[params["linenum"]]
	r.JSON(200,
		&Response{
			0,
			"OK",
			busline,
		},
	)
}
