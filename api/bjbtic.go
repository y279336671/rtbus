package api

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	BTIC_RAW_KEY_FMT = `aibang%s`
	BTIC_KEY_SECRET  = `bjjw_jtcx`
	btic_headers     = map[string]string{
		"HEADER_KEY_SECRET": "bjjw_jtcx",
		"IMSI":              "360120018215321",
		"UA":                "MATE8",
		"PLATFORM":          "android",
		"CUSTOM":            "aibang",
		"CID":               "67a88ec31de7a589a2344cc5d0469074",
		"IMEI":              "89031020265872",
		"NETWORK":           "gprs",
		"PKG_SOURCE":        "1",
		"CTYPE":             "json",
		"VID":               "5",
		"SOURCE":            "1",
		"UID":               "",
		"SID":               "",
		"PID":               "5",
		"Host":              "transapp.btic.org.cn:8512",
		"User-Agent":        "okhttp/3.3.1",
	}
	HEADER_TOKEN        = "ABTOKEN"
	BTIC_LISTLINES_PATH = `/ssgj/v1.0.0/checkupdate`
	BTIC_BUSLINE_PATH   = `/ssgj/v1.0.0/update`

	BTIC_BUS_PATH = `/ssgj/bus.php`
)

type BticAllLineResp struct {
	ErrCode     string `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	UpdateNum   string `json:"updateNum"`
	LineNum     string `json:"lineNum"`
	DataVersion string `json:"dataversion"`

	Lines struct {
		Line []*BticBasicLine `json:"line"`
	} `json:"lines"`
}

type BticBasicLine struct {
	ID       string `json:"id"`
	LineName string `json:"linename"`
	Classify string `json:"classify"`
	Status   string `json:"status"`
	Version  string `json:"version"`
}

type BticLineResp struct {
	ErrCode string `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	BusLine []*BticLine `json:"busline"`
}

type BticLine struct {
	ID         string  `json:"lineid"`
	ShortName  string  `json:"shotname"`
	LineName   string  `json:"linename"`
	Distince   float64 `json:"distince,string"`
	Ticket     string  `json:"ticket"`
	TotalPrice float64 `json:"totalPrice,string"`
	Time       string  `json:"time"`
	Type       string  `json:"type"`
	Coord      string  `json:"coord"`
	Status     int     `json:"status,string"`
	Version    int     `json:"version,string"`

	Stations struct {
		Station []*BticLineStation `json:"station"`
	} `json:"stations"`
}

type BticLineStation struct {
	Name string `json:"name"`
	No   string `json:"no"`
	Lon  string `json:"lon"`
	Lat  string `json:"lat"`
}

type BticLineRTResp struct {
	Root *BticLineRT `json:"root"`
}

type BticLineRT struct {
	Status  int    `json:"status,string"`
	Message string `json:"message"`
	Encrypt int    `json:"encrypt,string"`
	Num     int    `json:"num,string"`
	LineID  int    `json:"lid,string"`

	Data struct {
		Bus []*BticBus `json:"bus"`
	} `json:"data"`
}

type BticBus struct {
	GT                  string `json:"gt"`
	ID                  string `json:"id"`
	T                   string `json:"t"`
	NextStationName     string `json:"ns"`
	NextStationNum      string `json:"nsn"`
	NextStationDistance string `json:"nsd"`
	NextStationArrTime  string `json:"nst"`
	StationDistance     string `json:"sd"`
	StationArrTime      string `json:"st"`
	Lat                 string `json:"x"`
	Lon                 string `json:"y"`
}

func GetBticAllLine() (bjbls *CityBusLines, err error) {
	bjbls = NewCityBusLines()
	bjbls.CityInfo = &CityInfo{
		Code:   "010",
		ID:     "010",
		Name:   "北京",
		Hot:    1,
		PinYin: "beijing",
		Subway: 0,
	}

	allline := &BticAllLineResp{}
	var params *url.Values = &url.Values{}
	params.Set("version", "0")
	err = bticRequest(BTIC_LISTLINES_PATH, params, allline)
	if err != nil {
		return
	}
	if allline.ErrCode != "200" {
		err = errors.New(allline.ErrMsg)
		return
	}

	var rateLimit chan bool = make(chan bool, 50)
	var wg sync.WaitGroup
	for _, line := range allline.Lines.Line {
		rateLimit <- true

		wg.Add(1)
		go func(id string) {
			defer func() {
				<-rateLimit
				wg.Done()
			}()

			abline, err := getBticLine(id)
			if err != nil {
				LOGGER.Error("init BJ lineid %s failed:%v", id, err)
				return
			}

			bjbls.l.Lock()
			defer bjbls.l.Unlock()
			if bl, found := bjbls.ByLineName[abline.ShortName]; found {
				bdi := NewBusDirInfoByABLine(abline)
				if bdi != nil {
					bl.Put(bdi)
				}
			} else {
				bl := NewBusLineByABLine(abline)
				if bl != nil {
					bjbls.ByLineName[abline.ShortName] = bl
				}
			}
		}(line.ID)
	}

	//wait complete
	wg.Wait()
	close(rateLimit)

	return
}

func NewBusLineByABLine(line *BticLine) (bl *BusLine) {
	bdi := NewBusDirInfoByABLine(line)
	if bdi == nil {
		return
	}
	bdi.Direction = 0
	bdi.OtherDirIDs = []string{"1"}

	return &BusLine{
		LineNum:  line.ShortName,
		LineName: line.ShortName,
		Directions: map[string]*BusDirInfo{
			bdi.GetDirName(): bdi,
		},
	}
}

func NewBusDirInfoByABLine(line *BticLine) (bdi *BusDirInfo) {
	sNum := len(line.Stations.Station)
	if sNum == 0 {
		LOGGER.Warn("can't find the any station of line %s from aibang", line.LineName)
		return
	}

	firstS := line.Stations.Station[0]
	lastS := line.Stations.Station[sNum-1]

	var price string
	if line.TotalPrice == 0 {
		price = line.Ticket
	} else {
		price = fmt.Sprintf("%.0f", line.TotalPrice)
	}

	var firsrt_time, last_time string
	start_end_time := strings.SplitN(line.Time, "-", 2)
	firsrt_time = start_end_time[0]
	if len(start_end_time) > 1 {
		last_time = start_end_time[1]
	}

	bdi = &BusDirInfo{
		ID:          line.ID,
		Direction:   1,
		OtherDirIDs: []string{"0"},
		StartSn:     firstS.Name,
		EndSn:       lastS.Name,
		Price:       price,
		SnNum:       sNum,
		FirstTime:   firsrt_time,
		LastTime:    last_time,
		Stations:    convertABLineStation(line.Stations.Station),
	}

	bdi.Name = line.ShortName
	return
}

func convertABLineStation(abss []*BticLineStation) []*BusStation {
	var err error

	bss := make([]*BusStation, len(abss))
	for i, abs := range abss {
		bs := &BusStation{
			Name: abs.Name,
		}

		//StationNo
		bs.No, err = strconv.Atoi(abs.No)
		if err != nil {
			bs.No = i + 1
		}

		//lat lon
		bs.Lat, _ = strconv.ParseFloat(abs.Lat, 10)
		bs.Lon, _ = strconv.ParseFloat(abs.Lon, 10)

		bss[i] = bs
	}

	return bss
}

func getBticLine(id string) (line *BticLine, err error) {
	var params *url.Values = &url.Values{}
	params.Set("id", id)

	lineResp := &BticLineResp{}
	err = bticRequest(BTIC_BUSLINE_PATH, params, lineResp)
	if err != nil {
		return
	}
	if lineResp.BusLine == nil || len(lineResp.BusLine) == 0 {
		err = errors.New("busline info is null")
		return
	}

	//decrypt
	line = lineResp.BusLine[0]
	key := fmt.Sprintf(BTIC_RAW_KEY_FMT, id)
	line.ShortName, _ = Rc4DecodeString(key, line.ShortName)
	line.LineName, _ = Rc4DecodeString(key, line.LineName)
	line.Coord, _ = Rc4DecodeString(key, line.Coord)

	for _, station := range line.Stations.Station {
		station.Name, _ = Rc4DecodeString(key, station.Name)
		station.No, _ = Rc4DecodeString(key, station.No)
		station.Lon, _ = Rc4DecodeString(key, station.Lon)
		station.Lat, _ = Rc4DecodeString(key, station.Lat)
	}

	//fmt.Printf("%+v\n", line.Stations)

	return
}

func GetBticLineRT(bl *BusLine, dirname string) (rbus []*RunningBus, err error) {
	bdi, found := bl.GetBusDirInfo(dirname)
	if !found {
		err = errors.New(fmt.Sprintf("can't find the direction %s in BJ line %s", dirname, bl.LineNum))
		return
	}

	// BTIC_BUS_URL_PARAM_FMT = `versionid=5&encrypt=1&datatype=json&no=1&type=0&id=%s&city=%E5%8C%97%E4%BA%AC`
	curtime := time.Now().Unix()
	var params *url.Values = &url.Values{}
	params.Set("city", "北京")
	params.Set("id", bdi.ID)
	params.Set("no", "1")
	params.Set("type", "0")
	params.Set("datatype", "json")
	params.Set("encrypt", "1")
	params.Set("versionid", "5")

	linert_resp := &BticLineRTResp{}
	err = bticRequest(BTIC_BUS_PATH, params, linert_resp)
	if err != nil {
		return
	}
	linert := linert_resp.Root
	if linert == nil {
		LOGGER.Error("can't found root from %s", ToJsonString(linert_resp))
		err = errors.New("can't found root")
		return
	} else if linert.Status != 200 {
		LOGGER.Error("%s", ToJsonString(linert))
		err = errors.New(linert.Message)
		return
	} else {

	}

	rbus = make([]*RunningBus, len(linert.Data.Bus))
	for i, bus := range linert.Data.Bus {
		//fmt.Printf("bus:%+v\n", bus)
		rawkey := fmt.Sprintf(BTIC_RAW_KEY_FMT, bus.GT)
		bus.Lat, _ = Rc4DecodeString(rawkey, bus.Lat)
		bus.Lon, _ = Rc4DecodeString(rawkey, bus.Lon)
		bus.NextStationName, _ = Rc4DecodeString(rawkey, bus.NextStationName)
		bus.NextStationNum, _ = Rc4DecodeString(rawkey, bus.NextStationNum)
		//bus.NextStationDistance, _ = Rc4DecodeString(rawkey, bus.NextStationDistance)
		//bus.NextStationArrTime, _ = Rc4DecodeString(rawkey, bus.NextStationArrTime)
		bus.StationDistance, _ = Rc4DecodeString(rawkey, bus.StationDistance)
		bus.StationArrTime, _ = Rc4DecodeString(rawkey, bus.StationArrTime)

		// LOGGER.Warn(ToJsonString(bus))

		rbus[i] = &RunningBus{Name: bus.NextStationName, BusID: bus.ID, SyncTime: curtime}
		rbus[i].No, _ = strconv.Atoi(bus.NextStationNum)
		rbus[i].Lat, _ = strconv.ParseFloat(bus.Lat, 10)
		rbus[i].Lng, _ = strconv.ParseFloat(bus.Lon, 10)

		if bus.NextStationDistance == "0" || bus.NextStationDistance == "-1" {
			rbus[i].Status = BUS_ARRIVING_STATUS
			rbus[i].Distance = 0
		} else {
			rbus[i].Status = BUS_ARRIVING_FUTURE_STATUS
			rbus[i].Distance, _ = strconv.Atoi(bus.NextStationDistance)
		}
	}

	sort.Sort(SortRunningBus(rbus))
	return
}

type SortRunningBus []*RunningBus

func (rb SortRunningBus) Len() int      { return len(rb) }
func (rb SortRunningBus) Swap(i, j int) { rb[i], rb[j] = rb[j], rb[i] }
func (rb SortRunningBus) Less(i, j int) bool {
	if rb[i].No == rb[j].No {
		return rb[i].Status < rb[j].Status
	} else {
		return rb[i].No < rb[j].No
	}
}

func bticRequest(path string, params *url.Values, v interface{}) (err error) {
	req_url := "http://" + btic_headers["Host"] + path + "?" + params.Encode()
	// LOGGER.Info("url: %s", req_url)
	req, err := http.NewRequest(http.MethodGet, req_url, nil)
	if err != nil {
		return
	}

	// header
	for key, value := range btic_headers {
		req.Header.Set(key, value)
	}
	cur_time := fmt.Sprintf("%d", time.Now().Unix())
	req.Header.Set("TIME", cur_time)

	// token
	token := generateToken(cur_time, path)
	req.Header.Set("ABTOKEN", token)

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

	return json.Unmarshal(data, v)
}

func generateToken(cur_time, path string) string {
	body := fmt.Sprintf("%s%s%s", btic_headers["HEADER_KEY_SECRET"]+btic_headers["PLATFORM"]+btic_headers["CID"], cur_time, path)

	// LOGGER.Info("content: %s", body)
	sha1_data := fmt.Sprintf("%x", sha1.Sum([]byte(body)))
	token := fmt.Sprintf("%x", md5.Sum([]byte(sha1_data)))
	// LOGGER.Info("token: %s", token)
	return token
}
