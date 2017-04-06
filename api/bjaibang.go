package api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	AIBANG_RAW_KEY_FMT = `aibang%s`

	AIBANG_API_URL      = `http://mc.aibang.com/aiguang/bjgj.c`
	AIBANG_REALTIME_URL = `http://bjgj.aibang.com:8899/bus.php`
)

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

func GetAiBangAllLine() (bjbls *CityBusLines, err error) {
	bjbls = NewCityBusLines()

	allline := &AiBangAllLineResp{}
	var params *url.Values = &url.Values{}
	params.Set("m", "checkUpdate")
	params.Set("version", "1")
	err = aibangRequest(AIBANG_API_URL, params, allline)
	if err != nil {
		return
	}

	var rateLimit chan bool = make(chan bool, 50)
	var wg sync.WaitGroup
	for _, line := range allline.Lines {
		rateLimit <- true

		wg.Add(1)
		go func(id string) {
			defer func() {
				<-rateLimit
				wg.Done()
			}()

			abline, err := getAiBangLine(id)
			if err != nil {
				fmt.Println(err)
			}

			bjbls.l.Lock()
			defer bjbls.l.Unlock()
			if bl, found := bjbls.ByLineName[abline.ShortName]; found {
				bl.Put(NewBusDirInfoByABLine(abline))
			} else {
				bjbls.ByLineName[abline.ShortName] = NewBusLineByABLine(abline)
			}
		}(line.ID)
	}

	//wait complete
	wg.Wait()
	close(rateLimit)

	return
}

func NewBusLineByABLine(line *AiBangLine) (bl *BusLine) {
	bdi := NewBusDirInfoByABLine(line)
	if bdi == nil {
		return
	}
	bdi.Direction = 1

	return &BusLine{
		LineNum:  line.ShortName,
		LineName: line.ShortName,
		Directions: map[string]*BusDirInfo{
			bdi.Name: bdi,
		},
		getRT: GetAiBangLineRT,
	}
}

func NewBusDirInfoByABLine(line *AiBangLine) (bdi *BusDirInfo) {
	sNum := len(line.Stations)
	if sNum == 0 {
		return
	}

	firstS := line.Stations[0]
	lastS := line.Stations[sNum-1]

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
		ID:        line.ID,
		Direction: 0,
		StartSn:   firstS.No,
		EndSn:     lastS.No,
		Price:     price,
		SnNum:     sNum,
		FirstTime: firsrt_time,
		LastTime:  last_time,
		Stations:  convertABLineStation(line.Stations),
	}

	bdi.Name = bdi.GetDirName()
	return
}

func convertABLineStation(abss []*AiBangLineStation) []*BusStation {
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

func getAiBangLine(id string) (line *AiBangLine, err error) {
	var params *url.Values = &url.Values{}
	params.Set("m", "update")
	params.Set("id", id)

	lineResp := &AiBangLineResp{}
	err = aibangRequest(AIBANG_API_URL, params, lineResp)
	if err != nil {
		return
	}

	//decrypt
	line = lineResp.AiBangLine
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

	//fmt.Printf("%+v\n", line.Stations)

	return
}

func GetAiBangLineRT(bl *BusLine, dirname string) (rbus []*RunningBus, err error) {
	bdi, found := bl.GetBusDirInfo(dirname)
	if !found {
		err = errors.New(fmt.Sprintf("can't find the direction %s in BJ line %s", dirname, bl.LineNum))
		return
	}

	curtime := time.Now().Unix()
	var params *url.Values = &url.Values{}
	params.Set("city", "北京")
	params.Set("id", bdi.ID)
	params.Set("no", "1")
	params.Set("type", "2")
	params.Set("encrpt", "1")
	params.Set("versionid", "2")

	linert := &AiBangLineRT{}
	err = aibangRequest(AIBANG_REALTIME_URL, params, linert)
	if err != nil {
		return
	}

	rbus = make([]*RunningBus, len(linert.Data))
	for i, bus := range linert.Data {
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

		//fmt.Printf("%+v\n", bus)

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
	return
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
