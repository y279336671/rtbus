package api

import (
	"errors"
	"fmt"
	"github.com/bingbaba/util/httptool"
	"net/http"
	"sync"
	"time"
)

const (
	URL_CLL_REFER      = "http://web.chelaile.net.cn/ch5/index.html"
	FMT_CLL_URL_PARAMS = "lineId=%s&lineName=%s&lineNo=%s&s=h5&v=3.3.9&userId=browser_%d&h5Id=browser_%d&sign=1&cityId=%s"
	URL_CLL_BUS_URL    = "http://web.chelaile.net.cn/api/bus/line!lineDetail.action"
	FMT_CLL_URL_SEARCH = `http://web.chelaile.net.cn/api/basesearch/client/clientSearch.action?key=%s&count=3&s=h5&v=3.3.9&userId=browser_%d&h5Id=browser_%d&sign=1&cityId=%s`
)

type CllBus struct {
	l        sync.Mutex
	BusLines map[string]*BusLine
	CityInfo *CityInfo
}

type CllLineDirResp struct {
	Data *CllLineDirData `json:"data"`
}

type CllLineDirData struct {
	Line       *BusDirInfo   `json:"line"`
	Bus        []*RunningBus `json:"buses"`
	Stations   []*BusStation `json:"stations"`
	Otherlines []struct {
		LineId string `json:"lineid"`
	} `json:"otherlines"`
}

func GetCllAllBusLine() (cbls []*CityBusLines, err error) {
	cbls = make([]*CityBusLines, 0)

	var citys []*CityInfo
	citys, err = GetCllAllCitys()
	if err != nil {
		return
	}

	for _, city := range citys {
		cbls = append(cbls, &CityBusLines{
			Source:     SOURCE_CHELAILE,
			CityInfo:   city,
			ByLineName: make(map[string]*BusLine),
		})
	}

	return
}

func GetCllLineRT(cbl *CityBusLines, lineno, dirid string) (rbus []*RunningBus, err error) {
	//busline
	var bl *BusLine
	var hasInit bool
	bl, hasInit, err = getCllBusLine(cbl, lineno, dirid)
	if err != nil {
		return
	}

	//BusDirInfo
	bdi, found := bl.GetBusDirInfo(dirid)
	if !found {
		err = errors.New(fmt.Sprintf("can't find the direction %s in %s line %s", dirid, cbl.CityInfo.Name, bl.LineNum))
		return
	}

	//hasInit
	if !hasInit {
		rbus = bdi.RunningBuses
	} else { //fresh
		var cdd *CllLineDirData
		cdd, err = getCllLineDirData(cbl.CityInfo.ID, bdi.ID, lineno)
		if err != nil {
			return
		}

		rbus = cdd.Bus
		bdi.RunningBuses = cdd.Bus
	}

	return
}

func getCllBusLine(cbl *CityBusLines, lineno, dirid string) (bl *BusLine, hasInit bool, err error) {
	cbl.l.Lock()
	defer cbl.l.Unlock()

	//init
	if bl, hasInit = cbl.ByLineName[lineno]; !hasInit {
		bl, err = getCllLine(cbl, lineno)
		if err != nil {
			return
		}
		cbl.ByLineName[lineno] = bl
	}
	return
}

func getCllLineDirData(cityid, lineid, lineno string) (cdd *CllLineDirData, err error) {
	curtime := time.Now().UnixNano() / 1000000
	reqUrl := URL_CLL_BUS_URL +
		"?" +
		fmt.Sprintf(
			FMT_CLL_URL_PARAMS,
			lineid, lineno, lineno,
			curtime, curtime,
			cityid,
		)
	//fmt.Println(reqUrl) //debug

	var httpreq *http.Request
	httpreq, err = getCllHttpRequest(reqUrl)
	if err != nil {
		return
	}

	var cllresp *CllLineDirResp = &CllLineDirResp{}
	err = httptool.HttpDoJsonr(httpreq, cllresp)
	if err != nil {
		return
	}

	cdd = cllresp.Data
	return
}

type CllLineSearchResp struct {
	ErrMsg   string `json:"errmsg"`
	SVersion string `json:"sversion"`
	Data     struct {
		Lines []struct {
			EndSn  string `json:"endSn"`
			LineId string `json:"lineId"`
			LineNo string `json:"lineNo"`
			Name   string `json:"name"`
		} `json:"lines"`
	} `json:"data"`
}

func getCllLine(cbl *CityBusLines, lineno string) (bl *BusLine, err error) {
	city := cbl.CityInfo
	curtime := time.Now().UnixNano() / 1000000
	reqUrl := fmt.Sprintf(
		FMT_CLL_URL_SEARCH,
		lineno,
		curtime, curtime,
		city.ID,
	)
	//fmt.Println(reqUrl)

	var httpreq *http.Request
	httpreq, err = getCllHttpRequest(reqUrl)
	if err != nil {
		return
	}

	cllresp := &CllLineSearchResp{}
	err = httptool.HttpDoJsonr(httpreq, cllresp)
	if err != nil {
		return
	}

	if len(cllresp.Data.Lines) == 0 || cllresp.ErrMsg != "" {
		err = errors.New(fmt.Sprintf("search %s line failed:%s", lineno, cllresp.ErrMsg))
		return
	}

	var cdd *CllLineDirData
	lineid := cllresp.Data.Lines[0].LineId
	cdd, err = getCllLineDirData(city.ID, lineid, lineno)
	if err != nil {
		return
	}

	//BusDirInfo
	bdi := cdd.getBusDirInfo()

	//BusLine
	bl = &BusLine{
		LineNum:  lineno,
		LineName: lineno,
		Directions: map[string]*BusDirInfo{
			bdi.GetDirName(): bdi,
		},
	}

	//other line
	//fmt.Printf("%+v\n", cdd.Otherlines)
	if len(cdd.Otherlines) > 0 {
		lineid2 := cdd.Otherlines[0].LineId
		var cdd2 *CllLineDirData
		cdd2, err = getCllLineDirData(city.ID, lineid2, lineno)
		if err != nil {
			return
		}

		bdi2 := cdd2.getBusDirInfo()
		bl.Directions[bdi2.GetDirName()] = bdi2
	}

	return
}

func (cdd *CllLineDirData) getBusDirInfo() (bdi *BusDirInfo) {
	bdi = cdd.Line
	bdi.Stations = cdd.Stations
	bdi.RunningBuses = cdd.Bus

	curtime := time.Now().Unix()
	for _, rb := range bdi.RunningBuses {
		if rb.Distance == 0 {
			rb.Status = BUS_ARRIVING_STATUS
		} else {
			rb.Status = BUS_ARRIVING_FUTURE_STATUS
		}
		rb.Name = bdi.Stations[rb.No-1].Name
		rb.SyncTime = curtime - rb.SyncTime
	}

	return
}

func getCllHttpRequest(req_url string) (httpreq *http.Request, err error) {
	httpreq, err = http.NewRequest("GET", req_url, nil)
	if err != nil {
		return
	}

	httpreq.Header.Add("Accept", "application/json, text/plain, */*")
	httpreq.Header.Add("Referer", URL_CLL_REFER)
	httpreq.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.98 Mobile Safari/537.36`)
	return
}
