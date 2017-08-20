package handler

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xuebing1110/location/amap"
	"net/http"
)

func GetCityInfoByLocation(params martini.Params, r render.Render, httpreq *http.Request) {
	lat := params["lat"]
	lon := params["lon"]

	var city, cityname string

	req := amap.NewPoiSearchRequest(amapClient, "").
		SetAroundSearch(fmt.Sprintf("%s,%s", lon, lat)).
		SetPageSize(1)
	resp, err := req.Do()
	if err != nil {
		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
	}
	if len(resp.Pois) > 0 {
		city = resp.Pois[0].CityCode
		cityname = resp.Pois[0].CityName
	}

	r.JSON(200,
		&Response{
			0,
			"OK",
			map[string]string{
				"city":     city,
				"cityname": cityname,
			},
		},
	)
}
