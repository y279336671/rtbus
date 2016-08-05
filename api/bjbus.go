package api

import (
	"errors"
	"fmt"
	"github.com/bingbaba/util/httptool"
	"github.com/bingbaba/util/logs"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	REG_BUS_LINE_DIRECTION  = regexp.MustCompile(`<option value="(\d+)">[^\()]+\((\S+?)\)<\/option>`)
	REG_BUS_STATION         = regexp.MustCompile(`<option value="(\d+)">([^>]+)</option>`)
	REG_BUS_STATTION_STATUS = regexp.MustCompile(`(?s)<div\s+id="(\d+)(m?)\"><i\s+class="(\w+)"\s+clstag="`)
)

type BJBusSess struct {
	tokentime int64
	token     []*http.Cookie
	BusLines  map[string]*BusLine
	l         sync.Mutex
}

type StationStatusResp struct {
	HTML string `json:"html"`
	W    int    `json:"w"`
	Seq  string `json:"seq"`
}

func NewBJBusSess() (*BJBusSess, error) {
	b := &BJBusSess{
		BusLines: make(map[string]*BusLine),
	}

	err := b.refreshToken()
	if err != nil {
		return b, errors.New("get token failed:" + err.Error())
	} else {
		return b, nil
	}
}

func (b *BJBusSess) GetBusLine(lineid string) (*BusLine, error) {
	b.l.Lock()
	defer b.l.Unlock()

	_, found := b.BusLines[lineid]
	if !found {
		err := b.initBusline(lineid)
		if err != nil {
			return nil, err
		}
	}

	return b.BusLines[lineid], nil
}

func (b *BJBusSess) initBusline(lineid string) error {
	bl := NewBusLine(lineid)
	b.BusLines[lineid] = bl

	return b.initBusLineDirs(lineid)
}

func (b *BJBusSess) initBusLineDirs(lineid string) error {
	req_url := fmt.Sprintf(URL_BJ_FMT_LINE_DIRECTION, lineid)

	//new http req
	httpreq, err := b.newHttpRequest(req_url)
	if err != nil {
		return err
	}

	//http get
	id_direction_array, err := httptool.HttpDoAllRegexp(httpreq, REG_BUS_LINE_DIRECTION)
	if err != nil {
		return err
	} else if len(id_direction_array) == 0 {
		return errors.New("can't get anything of direction info!")
	}

	busline := b.BusLines[lineid]
	for _, id_direction := range id_direction_array {
		busdir := &BusDirInfo{
			ID:   id_direction[0],
			Name: id_direction[1],
		}
		busline.Direction = append(busline.Direction, busdir)

		//加载公交站信息
		err := b.loadBusStation(lineid, busdir.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

/*******************************************************
func (b *BJBusSess) FreshBusline(lineid string) error {
	busline := b.BusLines[lineid]

	for _, busdir := range busline.Direction {
		err := b.freshBuslineDir(lineid, busdir.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
*******************************************************/

func (b *BJBusSess) freshBuslineDir(lineid, dirid string) error {
	curtime := time.Now().Unix()
	bl := b.BusLines[lineid]
	busdir, _ := bl.getBusDir(dirid)

	//更新该方向bus信息
	req_url := fmt.Sprintf(URL_BJ_FMT_FRESH_STATION_STATUS, lineid, busdir.ID, 1)
	httpreq, err := b.newHttpRequest(req_url)
	if err != nil {
		return err
	}

	//http get
	status_resp := &StationStatusResp{}
	err = httptool.HttpDoJson(httpreq, status_resp)
	if err != nil {
		return errors.New("URL:" + req_url + ",error::::" + err.Error())
	}

	//解析http response
	station_status_array, err2 := httptool.MatchAll(REG_BUS_STATTION_STATUS, []byte(status_resp.HTML))
	if err2 != nil {
		return errors.New("URL:" + req_url + ",error::::" + err2.Error())
	}

	//当前公交状况
	map_cur := make(map[int][]*RunningBus)
	for _, station_status := range station_status_array {
		station_index, err := strconv.Atoi(station_status[0])
		if err == nil {
			if station_status[2] == "buss" { //到站
				buses_tmp, found := map_cur[station_index]
				if !found {
					buses_tmp = make([]*RunningBus, 0)
				}

				buses_tmp = append(buses_tmp,
					&RunningBus{
						Order:  station_index,
						Status: "1",
					})
				map_cur[station_index] = buses_tmp
			} else if station_status[2] == "busc" { //即将到站
				buses_tmp, found := map_cur[station_index]
				if !found {
					buses_tmp = make([]*RunningBus, 0)
				}

				buses_tmp = append(buses_tmp,
					&RunningBus{
						Order:  station_index,
						Status: "0.5",
					})
				map_cur[station_index] = buses_tmp
			}
		}
	}

	//更新存储
	for i := 0; i < len(busdir.Stations); i++ {
		order := busdir.Stations[i].Order
		buses_tmp, found := map_cur[order]
		if !found {
			buses_tmp = make([]*RunningBus, 0)
		}
		busdir.Stations[i].Buses = buses_tmp
	}

	busdir.freshTime = curtime

	return nil
}

func (b *BJBusSess) loadBusStation(linenum, dirid string) error {
	busline := b.BusLines[linenum]

	for _, busdir := range busline.Direction {
		if !busdir.equal(dirid) {
			continue
		}

		req_url := fmt.Sprintf(URL_BJ_FMT_LINE_STATION, linenum, busdir.ID)
		httpreq, err := b.newHttpRequest(req_url)
		if err != nil {
			return err
		}

		//http get
		station_array, err := httptool.HttpDoAllRegexp(httpreq, REG_BUS_STATION)
		if err != nil {
			return err
		} else if len(station_array) == 0 {
			return errors.New("can'get any statition of the line " + linenum)
		}

		for _, station := range station_array {
			order, _ := strconv.Atoi(station[0])
			busstation := &BusStation{
				Order: order,
				Sn:    station[1],
			}

			busdir.Stations = append(busdir.Stations, busstation)
		}

		return nil
	}

	return nil
}

func (b *BJBusSess) newHttpRequest(req_url string) (*http.Request, error) {
	httpreq, err := http.NewRequest("GET", req_url, nil)
	if err != nil {
		return nil, err
	}

	//刷新token
	b.refreshToken()

	for _, c := range b.token {
		httpreq.AddCookie(c)
	}
	httpreq.Header.Add("X-Requested-With", "XMLHttpRequest")

	return httpreq, nil
}

func (b *BJBusSess) Print() {
	logger := logs.GetBlogger()

	for _, busline := range b.BusLines {

		for _, busdir := range busline.Direction {
			last_index := len(busdir.Stations)
			logger.Info("lineNum:%s Direction:%s Station:%s ~ %s\n",
				busline.LineNum,
				busdir.Name,
				busdir.Stations[0].Sn,
				busdir.Stations[last_index-1].Sn,
			)

			for _, station := range busdir.Stations {
				if station.Buses != nil && len(station.Buses) > 0 {
					for _, bus := range station.Buses {
						logger.Info("%d:%s %s\n", station.Order, station.Sn, bus.Status)
					}
				}
			}
		}
	}
}

func (b *BJBusSess) refreshToken() error {
	curtime := time.Now().Unix()
	if time.Duration(curtime-b.tokentime)*time.Second < 30*time.Minute {
		return nil
	}

	var err error
	b.token, err = getToken()
	if err != nil {
		return errors.New("fresh token error:" + err.Error())
	}
	b.tokentime = curtime

	return nil
}
