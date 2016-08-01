package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xuebing1110/rtbus/api"
)

func LineInfoHandler(params martini.Params, r render.Render) {
	if BusSess == nil {
		r.JSON(
			502,
			&Response{502, "bjbus sess token error", nil},
		)
		return
	}

	// linenum, _ := url.QueryUnescape(params["linenum"])
	// fmt.Println(params["linenum"])

	stations, err := BusSess.GetLineInfo(params["linenum"], params["direction"])
	if err != nil {
		logger.Error("%v", err)

		r.JSON(
			502,
			&Response{503, err.Error(), nil},
		)
		return
	}

	//不要站牌名称
	s, found := params["simple"]
	if found && s != "0" && s != "false" {
		stations_tmp := make([]*api.BusStation, len(stations))
		for i, station := range stations {
			stations_tmp[i] = &api.BusStation{
				ID:     station.ID,
				Status: station.Status,
			}
		}
		r.JSON(200,
			&Response{
				0,
				"OK",
				stations_tmp,
			},
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
