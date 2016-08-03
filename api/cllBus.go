package api

import (
	"errors"
	"fmt"
	"github.com/bingbaba/util/httptool"
	"net/http"
)

const (
	URL_CLL_REFER      = "http://web.chelaile.net.cn/ch5/index.html"
	FMT_CLL_URL_PARAMS = "lineId=%s-%s-%s&lineName=%s&direction=%s&lineNo=%s&s=h5&v=3.1.3&userId=1&h5Id=1&sign=1&cityId=%s"
	URL_CLL_BUS_URL    = "http://web.chelaile.net.cn/api/bus/line!lineDetail.action"
)

var (
	MAP_CITY = map[string]CityInfo{
		"0539": CityInfo{
			CityID:  "009",
			Name:    "青岛",
			TelCode: "0539",
		},
	}
)

type CityInfo struct {
	CityID  string
	Name    string
	TelCode string
}

type CllBus struct {
	BusLines map[string]*BusLine
	CityInfo *CityInfo
}

type CllBusResp struct {
	LineInfo *CllLineResp `json:"line"`
}

type CllLineResp struct {
	Desc      string `json:"desc"`
	Direction string `json:"direction"`
	StartSn   string `json:"startsn,omitempty"`
	EndSn     string `json:"endsn,omitempty"`
	Price     string `json:"price,omitempty"`
	SnNum     int    `json:"stationsNum,omitempty"`
	FirstTime string `json:"firstTime,omitempty"`
	LastTime  string `json:"lastTime,omitempty"`
}

type CllLineDirBaseInfo struct {
	Data struct {
		Line *BusDirInfo `json:"line"`
	} `json:"data"`
	Bus      []*RunningBus `json:"buses"`
	Stations []*BusStation `json:"stations"`
}

func NewCllBus(citytel string) (*CllBus, error) {

	cityinfo, found := MAP_CITY[citytel]
	if !found {
		errors.New("can't support the city:" + citytel)
	}

	return &CllBus{
		BusLines: make(map[string]*BusLine),
		CityInfo: &cityinfo,
	}, nil
}

func (b *CllBus) LoadBusLineConf(lineid string) error {
	_, found := b.BusLines[lineid]
	if found {
		return nil
	}

	directions := make([]*BusDirInfo, 0)
	dir_arr := []string{"0", "1"}
	for _, dirid := range dir_arr {
		httreq, err := b.CityInfo.getHttpRequest(URL_CLL_BUS_URL, lineid, dirid)
		if err != nil {
			return err
		}

		cllresp := &CllLineDirBaseInfo{}
		err = httptool.HttpDoJsonr(httreq, cllresp)
		if err != nil {
			return err
		}

		busdir := cllresp.Data.Line
		busdir.Stations = cllresp.Stations
		directions = append(directions, busdir)
	}

	b.BusLines[lineid] = &BusLine{
		LineNum:   lineid,
		Direction: directions,
	}

	return nil
}

func (c *CityInfo) getHttpRequest(req_url, lineid, dirid string) (*http.Request, error) {
	httpreq, err := http.NewRequest("GET", req_url+c.getParams(lineid, dirid), nil)
	if err != nil {
		return nil, err
	}

	httpreq.Header.Add("Accept", "application/json, text/plain, */*")
	httpreq.Header.Add("Referer", URL_CLL_REFER)

	return httpreq, nil
}

func (c *CityInfo) getParams(lineid, dirid string) string {
	return fmt.Sprintf(
		FMT_CLL_URL_PARAMS,
		c.TelCode, lineid, dirid,
		lineid, c.CityID,
	)
}
