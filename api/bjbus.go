package api

import (
	"errors"
	"fmt"
	"github.com/bingbaba/util/httptool"
	"github.com/bingbaba/util/logs"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	//URL_BJ_HOME               = "http://www.bjbus.com/home/fun_rtbus.php"
	URL_BJ_HOME                     = "http://www.bjbus.com/home/index.php"
	URL_BJ_FMT_LINE_DIRECTION       = "http://www.bjbus.com/home/ajax_search_bus_stop_token.php?act=getLineDirOption&selBLine=%s"
	URL_BJ_FMT_LINE_STATION         = "http://www.bjbus.com/home/ajax_search_bus_stop_token.php?act=getDirStationOption&selBLine=%s&selBDir=%s"
	URL_BJ_FMT_FRESH_STATION_STATUS = "http://www.bjbus.com/home/ajax_search_bus_stop_token.php?act=busTime&selBLine=%s&selBDir=%s&selBStop=%d"

	REG_BUS_LINE_DIRECTION  = regexp.MustCompile(`<option value="(\d+)">[^\()]+\((\S+?)\)<\/option>`)
	REG_BUS_STATION         = regexp.MustCompile(`<option value="(\d+)">([^>]+)</option>`)
	REG_BUS_STATTION_STATUS = regexp.MustCompile(`(?s)<div\s+id="(\d+)(m?)\"><i\s+class="(\w+)"\s+clstag="`)
	REG_BUS_LINE_INFO       = regexp.MustCompile(`<article>\S+&nbsp;([0-9:]+)-([0-9:]+)&nbsp;(\S+)&nbsp;`)

	BJRtBusLines     *CityBusLines
	DefaultBJBusSess *BJBusSess

	BJBUS_HTTP_RETRY = 2
)

func init() {
	BJRtBusLines = NewCityBusLines()
	BJRtBusLines.CityInfo = &CityInfo{
		Code:   "010",
		ID:     "010",
		Name:   "北京",
		Hot:    1,
		PinYin: "beijing",
		Subway: 0,
	}

	DefaultBJBusSess = new(BJBusSess)
}

type BJBusSess struct {
	tokentime int64
	token     []*http.Cookie
	l         sync.Mutex
}

type StationStatusResp struct {
	HTML string `json:"html"`
	W    int    `json:"w"`
	Seq  string `json:"seq"`
}

func GetBJBusLine(lineno string) (bl *BusLine, err error) {
	inited := BJRtBusLines.hasInit(lineno)
	if !inited {
		bl, err = DefaultBJBusSess.NewBusLine(lineno)
		if err != nil {
			return
		} else {
			BJRtBusLines.ByLineName[lineno] = bl
		}
	} else {
		return BJRtBusLines.ByLineName[lineno], nil
	}

	return
}

func GetBJBusRT(lineno, dirid string) (rbus []*RunningBus, err error) {
	var bl *BusLine
	bl, err = GetBJBusLine(lineno)
	if err != nil {
		return
	}

	bdi, found := bl.GetBusDirInfo(dirid)
	if !found {
		err = errors.New(fmt.Sprintf("can't find the direction %s in BJ line %s from bjrtus", dirid, bl.LineNum))
		return
	}

	return DefaultBJBusSess.getBusRt(bdi, lineno, false)
}

func (b *BJBusSess) NewBusLine(lineno string) (bl *BusLine, err error) {
	req_url := fmt.Sprintf(URL_BJ_FMT_LINE_DIRECTION, lineno)

	//new http req
	var httpreq *http.Request
	httpreq, err = b.newHttpRequest(req_url)
	if err != nil {
		return
	}

	//http get
	var id_direction_array [][]string
	for i := 0; i <= BJBUS_HTTP_RETRY; i++ {
		id_direction_array, err = httptool.HttpDoAllRegexp(httpreq, REG_BUS_LINE_DIRECTION)
		if isHttpTimeOut(err) {
			if i == BJBUS_HTTP_RETRY {
				return
			} else {
				continue
			}
		}

		if err != nil {
			return
		} else if len(id_direction_array) == 0 {
			err = errors.New("can't get anything of direction info!")
			return
		} else {
			break
		}
	}

	bl = &BusLine{
		LineNum:    lineno,
		LineName:   lineno,
		Directions: make(map[string]*BusDirInfo),
	}
	var dirLength = len(id_direction_array)
	for dindex, id_direction := range id_direction_array {
		var bdi *BusDirInfo
		bdi, err = b.newBusDirInfo(lineno, strconv.Itoa(dindex), id_direction[0])
		if err != nil {
			return
		}

		for i := 0; i < dirLength; i++ {
			if i == dindex {
				continue
			}
			bdi.OtherDirIDs = append(bdi.OtherDirIDs, strconv.Itoa(i))
		}

		//公交线路信息
		_, err = b.getBusRt(bdi, lineno, true)
		if err != nil {
			return
		}

		bl.Directions[bdi.GetDirName()] = bdi
	}

	return
}

func (b *BJBusSess) newBusDirInfo(lineno, dindex, dirid string) (bdi *BusDirInfo, err error) {
	req_url := fmt.Sprintf(URL_BJ_FMT_LINE_STATION, lineno, dirid)
	var httpreq *http.Request
	httpreq, err = b.newHttpRequest(req_url)
	if err != nil {
		return
	}

	//http get
	var station_array [][]string
	for i := 0; i <= BJBUS_HTTP_RETRY; i++ {
		station_array, err = httptool.HttpDoAllRegexp(httpreq, REG_BUS_STATION)
		if isHttpTimeOut(err) {
			if i == BJBUS_HTTP_RETRY {
				return
			} else {
				continue
			}
		}

		if err != nil {
			return
		} else if len(station_array) == 0 {
			err = errors.New("can'get any statition of the line " + lineno)
			return
		} else {
			break
		}
	}

	bdi = &BusDirInfo{
		did:         dirid,
		ID:          dindex,
		Name:        lineno,
		OtherDirIDs: make([]string, 0),
	}

	for _, station := range station_array {
		order, _ := strconv.Atoi(station[0])
		busstation := &BusStation{
			No:   order,
			Name: station[1],
		}

		bdi.Stations = append(bdi.Stations, busstation)
	}

	bdi.StartSn = bdi.Stations[0].Name
	bdi.EndSn = bdi.Stations[len(bdi.Stations)-1].Name

	return
}

func (b *BJBusSess) getBusRt(bdi *BusDirInfo, lineno string, need_init bool) (rbs []*RunningBus, err error) {
	curtime := time.Now().Unix()
	if curtime-bdi.freshTime <= 5 {
		return bdi.RunningBuses, nil
	}

	var httpreq *http.Request
	req_url := fmt.Sprintf(URL_BJ_FMT_FRESH_STATION_STATUS, lineno, bdi.did, 1)
	httpreq, err = b.newHttpRequest(req_url)
	if err != nil {
		return
	}

	//http get
	status_resp := &StationStatusResp{}

	for i := 0; i <= BJBUS_HTTP_RETRY; i++ {
		err = httptool.HttpDoJson(httpreq, status_resp)

		if isHttpTimeOut(err) {
			if i == BJBUS_HTTP_RETRY {
				return
			} else {
				continue
			}
		}

		if err != nil {
			err = errors.New("URL:" + req_url + ",error::::" + err.Error())
			return
		} else {
			break
		}
	}

	//第一次需查询公交线路基础信息： 票价、发车时间等等
	if need_init {
		var lineinfo []string
		lineinfo, err = httptool.Match(REG_BUS_LINE_INFO, []byte(status_resp.HTML))
		if err != nil {
			return
		}

		if len(lineinfo) == 0 {
			err = errors.New("can't get line base info from " + status_resp.HTML)
			return
		}
		bdi.FirstTime = lineinfo[0]
		bdi.LastTime = lineinfo[1]
		bdi.Price = lineinfo[2]
	}

	//解析http response
	station_status_array, err2 := httptool.MatchAll(REG_BUS_STATTION_STATUS, []byte(status_resp.HTML))
	if err2 != nil {
		err = errors.New("URL:" + req_url + ",error::::" + err2.Error())
		return
	}

	//当前公交状况
	rbs = make([]*RunningBus, 0)
	for _, station_status := range station_status_array {
		station_index, err := strconv.Atoi(station_status[0])
		if err == nil {
			if station_status[2] == "buss" { //到站
				rbs = append(rbs,
					&RunningBus{
						No:       station_index,
						Status:   BUS_ARRIVING_STATUS,
						SyncTime: curtime,
					},
				)
			} else if station_status[2] == "busc" { //即将到站
				rbs = append(rbs,
					&RunningBus{
						No:       station_index,
						Status:   BUS_ARRIVING_FUTURE_STATUS,
						SyncTime: curtime,
					},
				)
			}
		}
	}

	//更新存储
	bdi.RunningBuses = rbs
	bdi.freshTime = curtime

	return
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

func (b *BJBusSess) refreshToken() error {
	curtime := time.Now().Unix()
	if time.Duration(curtime-b.tokentime)*time.Second < 10*time.Minute {
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

func getToken() ([]*http.Cookie, error) {
	logs.GetBlogger().Info("refresh token...")

	resp, err := http.Get(URL_BJ_HOME)
	if err != nil {
		return make([]*http.Cookie, 0), err
	}

	// //test
	// for _, c := range resp.Cookies() {
	//  fmt.Println(c.Raw)
	// }

	return resp.Cookies(), nil
}

func isHttpTimeOut(err error) bool {
	if err == nil {
		return false
	} else if strings.Index(err.Error(), "timeout") >= 0 {
		return true
	} else {
		return false
	}
}
