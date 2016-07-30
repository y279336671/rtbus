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

type BusLine struct {
	LineNum   string        `json:"linenum"`
	Direction []*BusDirInfo `json:"direction"`
}

type BusDirInfo struct {
	l          sync.Mutex
	freshTime  int64
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Stations   []*BusStation  `json:"stations"`
	Name2Index map[string]int `json:"-"`
}

type BusStation struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
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

	err := b.FreshToken()
	if err != nil {
		return b, errors.New("get token failed:" + err.Error())
	} else {
		return b, nil
	}
}

func (b *BJBusSess) FreshToken() error {
	b.l.Lock()
	defer b.l.Unlock()

	curtime := time.Now().Unix()
	if time.Duration(curtime-b.tokentime)*time.Second < 30*time.Minute {
		return nil
	}

	var err error
	b.token, err = getToken()
	if err != nil {
		return err
	}
	b.tokentime = curtime

	return nil
}

func (b *BJBusSess) LoadBusLineConf(linenum string) error {
	b.l.Lock()
	defer b.l.Unlock()

	err := b.loadBusLineDirection(linenum)
	if err != nil {
		return errors.New(fmt.Sprintf("get %s line direction info failed:%s", linenum, err.Error()))
	}

	busline := b.BusLines[linenum]
	for _, busdir := range busline.Direction {
		err_tmp := b.loadBusStation(linenum, busdir.ID)
		if err_tmp != nil {
			return errors.New(fmt.Sprintf("get %s(%s) failed:%v", linenum, busdir.ID, err_tmp))
		}
	}

	return nil
}

func (b *BJBusSess) FreshStatus(linenum, direction string) error {
	err := b.FreshToken()
	if err != nil {
		return err
	}

	busdir, err_tmp := b.getBusDir(linenum, direction)
	if err_tmp != nil {
		return err_tmp
	}

	err_tmp = b.freshStatus(linenum, busdir)
	if err_tmp != nil {
		return err_tmp
	}

	return nil
}

func (b *BJBusSess) FreshStatusByStation(linenum, direction, station string) error {
	err := b.FreshToken()
	if err != nil {
		return err
	}

	busdir, err_tmp := b.getBusDir(linenum, direction)
	if err_tmp != nil {
		return err_tmp
	}

	err_tmp = b.freshStatus(linenum, busdir)
	if err_tmp != nil {
		return err_tmp
	}

	return nil
}

func (b *BJBusSess) freshStatus(linenum string, busdir *BusDirInfo) error {
	busdir.l.Lock()
	defer busdir.l.Unlock()

	if time.Now().Unix()-busdir.freshTime < 10 {
		return nil
	}

	req_url := fmt.Sprintf(URL_BJ_FMT_FRESH_STATION_STATUS, linenum, busdir.ID, 1)

	//new http req
	httpreq, err := b.newHttpRequest(req_url)
	if err != nil {
		return err
	}

	//http get
	status_resp := &StationStatusResp{}
	err = httptool.HttpDoJson(httpreq, status_resp)
	if err != nil {
		return err
	}

	station_status_array, err2 := httptool.MatchAll(REG_BUS_STATTION_STATUS, []byte(status_resp.HTML))
	if err2 != nil {
		return err2
	}

	map_cur := make(map[int]string)
	for _, station_status := range station_status_array {
		station_index, err := strconv.Atoi(station_status[0])
		if err == nil {
			if station_status[2] == "buss" {
				map_cur[station_index-1] = "1" //到站
			} else if station_status[2] == "busc" {
				map_cur[station_index-1] = "0.5" //即将到站
			}

		}
	}

	for i := 0; i < len(busdir.Stations); i++ {
		status, found := map_cur[i]
		if found {
			busdir.Stations[i].Status = status
		} else {
			busdir.Stations[i].Status = ""
		}
	}

	busdir.freshTime = time.Now().Unix()

	return nil
}

func (b *BJBusSess) loadBusLineDirection(linenum string) error {
	busline := &BusLine{LineNum: linenum}

	req_url := fmt.Sprintf(URL_BJ_FMT_LINE_DIRECTION, linenum)

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

	for _, id_direction := range id_direction_array {
		BusDirInfo := &BusDirInfo{
			ID:   id_direction[0],
			Name: id_direction[1],
		}

		busline.Direction = append(busline.Direction, BusDirInfo)
	}
	b.BusLines[linenum] = busline

	return nil

}

func (b *BJBusSess) loadBusStation(linenum, dirid string) error {
	busline := b.BusLines[linenum]

	for _, busdir := range busline.Direction {
		if len(busdir.Stations) > 0 {
			continue
		}

		if busdir.Name2Index == nil {
			busdir.Name2Index = make(map[string]int)
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
			busstation := &BusStation{
				ID:   station[0],
				Name: station[1],
			}

			busdir.Name2Index[station[1]] = len(busdir.Stations)
			busdir.Stations = append(busdir.Stations, busstation)
		}
	}

	return nil
}

func (b *BJBusSess) newHttpRequest(req_url string) (*http.Request, error) {
	httpreq, err := http.NewRequest("GET", req_url, nil)
	if err != nil {
		return nil, err
	}

	//6小时刷新一次token
	if time.Now().Unix()-b.tokentime >= 3600*6 {
		b.FreshToken()
	}

	for _, c := range b.token {
		httpreq.AddCookie(c)
	}
	httpreq.Header.Add("X-Requested-With", "XMLHttpRequest")

	return httpreq, nil
}

func (b *BJBusSess) GetLineInfo(linenum, direction string) ([]*BusStation, error) {
	err := b.FreshStatus(linenum, direction)
	if err != nil {
		return nil, err
	}

	busdir, err2 := b.getBusDir(linenum, direction)
	if err2 != nil {
		return nil, err2
	}

	return busdir.Stations, nil

}

func (b *BJBusSess) getBusDir(linenum, direction string) (*BusDirInfo, error) {
	busline, found := b.BusLines[linenum]
	if !found {
		err := b.LoadBusLineConf(linenum)
		if err != nil {
			return nil, err
		}
		busline = b.BusLines[linenum]
	}

	for _, busdir := range busline.Direction {
		if busdir.Name != direction && busdir.ID != direction {
			continue
		}

		err_tmp := b.freshStatus(linenum, busdir)
		if err_tmp != nil {
			return nil, err_tmp
		}

		return busdir, nil
	}

	return nil, errors.New("not found")
}

func (b *BJBusSess) Print() {
	logger := logs.GetBlogger()

	for _, busline := range b.BusLines {

		for _, busdir := range busline.Direction {
			last_index := len(busdir.Stations)
			logger.Info("lineNum:%s Direction:%s Station:%s ~ %s\n",
				busline.LineNum,
				busdir.Name,
				busdir.Stations[0].Name,
				busdir.Stations[last_index-1].Name,
			)

			for _, station := range busdir.Stations {
				if station.Status != "" {
					logger.Info("%s:%s %s\n", station.ID, station.Name, station.Status)
				}
			}
		}
	}
}
