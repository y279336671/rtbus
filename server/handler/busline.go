package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func LineStationHandler(params martini.Params, r render.Render) {
	if BusSess == nil {
		r.JSON(
			502,
			&Response{502, "bjbus sess token error", nil},
		)
		return
	}

	stations, err := BusSess.GetLineStationInfo(params["linenum"], params["direction"])
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

func LineBusHandler(params martini.Params, r render.Render) {
	if BusSess == nil {
		r.JSON(
			502,
			&Response{502, "bjbus sess token error", nil},
		)
		return
	}

	buses, err := BusSess.GetLineBusInfo(params["linenum"], params["direction"])
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
