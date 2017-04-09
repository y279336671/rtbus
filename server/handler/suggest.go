package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xuebing1110/location/amap"
	"github.com/xuebing1110/rtbus/api"
	"net/http"
	"strings"
	"sync"
)

type NearestBusStation struct {
	City        string                `json:"city"`
	CityName    string                `json:"cityname"`
	StationName string                `json:"sn"`
	LineNos     []string              `json:"linenos"`
	Lines       []*BusLineDirOverview `json:"lines"`
}

type BusLineDirOverview struct {
	LineNo     string            `json:"lineno"`
	Direction  int               `json:"linedir"`
	AnotherDir string            `json:"another_dir"`
	StartSn    string            `json:"startsn,omitempty"`
	EndSn      string            `json:"endsn,omitempty"`
	Price      string            `json:"price,omitempty"`
	SnNum      int               `json:"stationsNum,omitempty"`
	FirstTime  string            `json:"firstTime,omitempty"`
	LastTime   string            `json:"lastTime,omitempty"`
	Buses      []*api.RunningBus `json:"buses"`
	IsSupport  bool              `json:"issupport"`
}

var (
	amapClient *amap.Client
	AMAP_KEY   = `b3abf03fa1e83992727f0625a918fe73`
)

func init() {
	amapClient = amap.NewClient(AMAP_KEY)
	amapClient.HttpClient = api.DEFAULT_HTTP_CLIENT
}

func BusLineOverview(params martini.Params, r render.Render) {
	city := params["city"]
	linenos := params["linenos"]
	sn := params["station"]

	var lock sync.Mutex
	var wg sync.WaitGroup
	var bldos = make([]*BusLineDirOverview, 0)
	var lineno_array = strings.Split(linenos, ",")
	for _, lineno := range lineno_array {
		wg.Add(1)
		go func(lineno string) {
			defer wg.Done()
			bldo := GetBusLineDirOverview(city, lineno, sn, true)

			lock.Lock()
			defer lock.Unlock()
			bldos = append(bldos, bldo)
		}(lineno)
	}
	wg.Wait()

	r.JSON(200,
		&Response{
			0,
			"OK",
			bldos,
		},
	)

}

func BusLineSuggest(params martini.Params, r render.Render, httpreq *http.Request) {
	lat := params["lat"]
	lon := params["lon"]

	httpreq.ParseForm()
	lazy := httpreq.Form.Get("lazy")

	//nearest bus line
	req := amap.NewInputtipsRequest(amapClient, "公交车站").
		SetCityLimit(true).
		SetLatLon(lat, lon)
	resp, err := req.Do()
	if err != nil {
		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	nbss := make([]*NearestBusStation, 0)
	for sni, tip := range resp.Tips {
		sn := strings.TrimRight(tip.Name, "(公交站)")
		nbs := &NearestBusStation{
			City:        req.City,
			CityName:    req.CityName,
			StationName: sn,
		}

		//lazy load
		var loadBus bool = true
		if lazy != "" && sni > 0 {
			loadBus = false
		}

		var linenos = make([]string, 0)
		var bldos = make([]*BusLineDirOverview, 0)
		var lock sync.Mutex
		var wg sync.WaitGroup
		linenames := strings.Split(tip.Address, ";")
		for _, linename := range linenames {
			//lineno
			lineno := strings.TrimRight(linename, "线")
			lineno = strings.TrimRight(lineno, "路")
			lineno = strings.Replace(lineno, "路内环", "内", 1)
			lineno = strings.Replace(lineno, "路外环", "外", 1)
			linenos = append(linenos, lineno)

			wg.Add(1)
			go func(lineno string) {
				defer wg.Done()
				bldo := GetBusLineDirOverview(req.City, lineno, sn, loadBus)

				lock.Lock()
				defer lock.Unlock()
				bldos = append(bldos, bldo)
			}(lineno)
		}

		wg.Wait()
		nbs.Lines = bldos
		nbs.LineNos = linenos
		nbss = append(nbss, nbs)
	}

	r.JSON(200,
		&Response{
			0,
			"OK",
			nbss,
		},
	)
}

func GetBusLineDirOverview(city, lineno, station string, loadBus bool) (bldo *BusLineDirOverview) {
	bldo = &BusLineDirOverview{
		LineNo:    lineno,
		IsSupport: false,
	}

	dirid := "0"
	bdi, err := BusTool.GetBusLineDirInfo(city, lineno, dirid)
	if err != nil {
		logger.Warn("%v", err)
		return
	}

	bldo.IsSupport = true
	bldo.Direction = bdi.Direction
	bldo.StartSn = bdi.StartSn
	bldo.EndSn = bdi.EndSn
	bldo.Price = bdi.Price
	bldo.SnNum = bdi.SnNum
	bldo.FirstTime = bdi.FirstTime
	bldo.LastTime = bdi.LastTime

	if len(bdi.OtherDirIDs) > 0 {
		bldo.AnotherDir = bdi.OtherDirIDs[0]
	}

	//get station index
	var stationIndex int
	for _, s := range bdi.Stations {
		if s.Name == station {
			stationIndex = s.No
			break
		}
	}

	//get running buses
	if loadBus {
		rbuses, err := BusTool.GetRT(city, lineno, dirid)
		if err != nil {
			logger.Warn("%v", err)
		} else {
			bldo.Buses = make([]*api.RunningBus, 0)
			for _, rbus := range rbuses {
				if rbus.No <= stationIndex {
					bldo.Buses = append(bldo.Buses, rbus)
				} else {
					break
				}
			}
		}
	}

	return bldo
}
