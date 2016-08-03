package api

import (
	"errors"
	"fmt"
	"github.com/bingbaba/util/httptool"
	"net/http"
	"time"
)

const (
	URL_CLL_REFER      = "http://web.chelaile.net.cn/ch5/index.html"
	FMT_CLL_URL_PARAMS = "lineId=%s-%s-%s&lineName=%s&direction=%s&lineNo=%s&s=h5&v=3.1.3&userId=1&h5Id=1&sign=1&cityId=%s"
	URL_CLL_BUS_URL    = "http://web.chelaile.net.cn/api/bus/line!lineDetail.action"
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

type CllLineDirBaseInfo struct {
	Data struct {
		Line     *BusDirInfo   `json:"line"`
		Bus      []*RunningBus `json:"buses"`
		Stations []*BusStation `json:"stations"`
	} `json:"data"`
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

func (b *CllBus) GetBusLine(lineid string) (*BusLine, error) {
	_, found := b.BusLines[lineid]
	if !found {
		err := b.initBusline(lineid)
		if err != nil {
			return nil, err
		}
	}

	return b.BusLines[lineid], nil
}

func (b *CllBus) initBusline(lineid string) error {
	bl := NewBusLine(lineid)
	b.BusLines[lineid] = bl

	return b.freshBusline(lineid)
}

func (b *CllBus) freshBusline(lineid string) error {
	dir_arr := []string{"0", "1"}
	for _, dirid := range dir_arr {
		err := b.freshBuslineDir(lineid, dirid)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *CllBus) freshBuslineDir(lineid, dirid string) error {
	httreq, err := b.CityInfo.getHttpRequest(URL_CLL_BUS_URL, lineid, dirid)
	if err != nil {
		return err
	}

	cllresp := &CllLineDirBaseInfo{}
	err = httptool.HttpDoJsonr(httreq, cllresp)
	if err != nil {
		return err
	}

	//初始化
	bl := b.BusLines[lineid]
	busdir, err := bl.GetBusInfo(dirid)

	//第一次加载(bus+station)
	if err != nil {
		busdir = cllresp.Data.Line
		busdir.Name = busdir.StartSn + "-" + busdir.EndSn
		busdir.ID = fmt.Sprintf("%d", busdir.Direction)
		busdir.Direction = 0
		busdir.Stations = cllresp.Data.Stations
		bl.Direction = append(bl.Direction, busdir)
	}

	//更新bus信息
	for _, s := range busdir.Stations {
		s.Buses = make([]*RunningBus, 0)
		for _, rbus := range cllresp.Data.Bus {
			if s.Order == rbus.Order {
				if rbus.Distance == 0 {
					rbus.Status = "1"
				} else {
					rbus.Status = "0.5"
				}
				s.Buses = append(s.Buses, rbus)
			}
		}
	}

	busdir.freshTime = time.Now().Unix()

	return nil
}

func (c *CityInfo) getHttpRequest(req_url, lineid, dirid string) (*http.Request, error) {
	req_url = req_url + "?" + c.getParams(lineid, dirid)
	fmt.Println(req_url)

	httpreq, err := http.NewRequest("GET", req_url, nil)
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
		lineid, dirid, lineid,
		c.CityID,
	)
}
