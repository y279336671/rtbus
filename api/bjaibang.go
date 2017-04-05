package api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

var (
	AIBANG_RAW_KEY_FMT = `aibang%s`

	AIBANG_API_URL      = `http://mc.aibang.com/aiguang/bjgj.c`
	AIBANG_REALTIME_URL = `http://bjgj.aibang.com:8899/bus.php`
)

type BJAiBangBusSess struct {
	BusLines map[string]*BusLine
	l        sync.Mutex
}

type AiBangAllLineResp struct {
	XMLName     xml.Name `xml:"root"`
	Status      int      `xml:"status"`
	UpdateNum   int      `xml:"updateNum"`
	LineNum     int      `xml:"lineNum"`
	DataVersion string   `xml:"dataversion"`

	Lines []*AiBangLineId `xml:"lines>line"`
}

type AiBangLineId struct {
	ID      string `xml:"id"`
	Status  int    `xml:"status"`
	Version int    `xml:"version"`
}

type AiBangLineResp struct {
	XMLName    xml.Name    `xml:"root"`
	AiBangLine *AiBangLine `xml:"busline"`
}

type AiBangLine struct {
	ID         string  `xml:"lineid"`
	ShortName  string  `xml:"shotname"`
	LineName   string  `xml:"linename"`
	Distince   float64 `xml:"distince"`
	Ticket     string  `xml:"ticket"`
	TotalPrice float64 `xml:"totalPrice"`
	Time       string  `xml:"time"`
	Type       string  `xml:"type"`
	Coord      string  `xml:"coord"`
	Status     int     `xml:"status"`
	Version    int     `xml:"version"`

	Stations []*AiBangLineStation `xml:"stations>station"`
}

type AiBangLineStation struct {
	Name string `xml:"name"`
	No   string `xml:"no"`
	Lon  string `xml:"lon"`
	Lat  string `xml:"lat"`
}

type AiBangLineRT struct {
	XMLName xml.Name `xml:"root"`
	Status  int      `xml:"status"`
	Message string   `xml:"message"`
	Encrypt int      `xml:"encrypt"`
	Num     int      `xml:"num"`
	LineID  int      `xml:"lid"`

	Data []*AiBangBus `xml:"data>bus"`
}

type AiBangBus struct {
	GT                  string `xml:"gt"`
	ID                  string `xml:"id"`
	T                   string `xml:"t"`
	NextStationName     string `xml:"ns"`
	NextStationNum      string `xml:"nsn"`
	NextStationDistance string `xml:"nsd"`
	NextStationArrTime  string `xml:"nst"`
	StationDistance     string `xml:"sd"`
	StationArrTime      string `xml:"st"`
	Lat                 string `xml:"x"`
	Lon                 string `xml:"y"`
}

func NewBJAiBangBusSess() (*BJAiBangBusSess, error) {
	b := &BJAiBangBusSess{
		BusLines: make(map[string]*BusLine),
	}

	err := b.getAiBangAllLine()
	if err != nil {
		return b, errors.New("get token failed:" + err.Error())
	} else {
		return b, nil
	}
}

func (b *BJAiBangBusSess) getAiBangAllLine() error {
	b.l.Lock()
	defer b.l.Unlock()

	allline := &AiBangAllLineResp{}
	var params *url.Values = &url.Values{}
	params.Set("m", "checkUpdate")
	params.Set("version", "1")
	err := aibangRequest(AIBANG_API_URL, params, allline)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, line := range allline.Lines {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := b.getline(id)
			if err != nil {
				fmt.Println(err)
			}
		}(line.ID)
	}

	wg.Wait()

	/*_, found := b.BusLines[lineid]
	if !found {
		err := b.initBusline(lineid)
		if err != nil {
			return nil, err
		}
	}*/

	return nil
}

func (b *BJAiBangBusSess) Print() {

}

func (b *BJAiBangBusSess) getline(id string) (err error) {
	var params *url.Values = &url.Values{}
	params.Set("m", "update")
	params.Set("id", id)

	lineResp := &AiBangLineResp{}
	err = aibangRequest(AIBANG_API_URL, params, lineResp)
	if err != nil {
		return
	}

	//decrypt
	line := lineResp.AiBangLine
	key := fmt.Sprintf(AIBANG_RAW_KEY_FMT, id)
	line.ShortName, _ = Rc4DecodeString(key, line.ShortName)
	line.LineName, _ = Rc4DecodeString(key, line.LineName)
	line.Coord, _ = Rc4DecodeString(key, line.Coord)

	for _, station := range line.Stations {
		station.Name, _ = Rc4DecodeString(key, station.Name)
		station.No, _ = Rc4DecodeString(key, station.No)
		station.Lon, _ = Rc4DecodeString(key, station.Lon)
		station.Lat, _ = Rc4DecodeString(key, station.Lat)
	}

	line.Print()

	return nil
}

func (b *BJAiBangBusSess) getlineRealTime(id, no string) (err error) {
	var params *url.Values = &url.Values{}
	params.Set("city", "北京")
	params.Set("id", id)
	params.Set("no", no)
	params.Set("type", "2")
	params.Set("encrpt", "1")
	params.Set("versionid", "2")

	linert := &AiBangLineRT{}
	err = aibangRequest(AIBANG_REALTIME_URL, params, linert)
	if err != nil {
		return
	}

	for _, bus := range linert.Data {
		//fmt.Printf("bus:%+v\n", bus)
		rawkey := fmt.Sprintf(AIBANG_RAW_KEY_FMT, bus.GT)
		bus.Lat, _ = Rc4DecodeString(rawkey, bus.Lat)
		bus.Lon, _ = Rc4DecodeString(rawkey, bus.Lon)
		bus.NextStationName, _ = Rc4DecodeString(rawkey, bus.NextStationName)
		bus.NextStationNum, _ = Rc4DecodeString(rawkey, bus.NextStationNum)
		//bus.NextStationDistance, _ = Rc4DecodeString(rawkey, bus.NextStationDistance)
		//bus.NextStationArrTime, _ = Rc4DecodeString(rawkey, bus.NextStationArrTime)
		bus.StationDistance, _ = Rc4DecodeString(rawkey, bus.StationDistance)
		bus.StationArrTime, _ = Rc4DecodeString(rawkey, bus.StationArrTime)

		fmt.Printf("%+v\n", bus)
	}
	return
}

func (bl *AiBangLine) Print() {
	fmt.Printf("id:%s, shortname:%s, linename: %s, distance:%f, ticket:%s, totalPrice:%f, time:%s, type:%s\n",
		bl.ID, bl.ShortName, bl.LineName,
		bl.Distince, bl.Ticket, bl.TotalPrice,
		bl.Time, bl.Type)
	//fmt.Printf("coord:%s", bl.Coord)

	for i, station := range bl.Stations {
		fmt.Printf("id:%s[%d] stationName:%s, stationNo:%s, Lon:%s, Lat:%s\n",
			bl.ID, i, station.Name,
			station.No, station.Lon, station.Lat,
		)
	}
}

func aibangRequest(reqUrl string, params *url.Values, v interface{}) (err error) {
	req, err := http.NewRequest(http.MethodGet, reqUrl+"?"+params.Encode(), nil)
	if err != nil {
		return
	}
	req.Header.Add("cid", "1024")

	var resp *http.Response
	resp, err = HttpDo(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	//read all
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return xml.Unmarshal(data, v)
}

/*func (b *BJBusSess) GetBusLine(lineid string) (*BusLine, error) {
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
*/
