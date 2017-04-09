package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xuebing1110/location/amap"
	"github.com/xuebing1110/rtbus/api"
	"strings"
	"sync"
)

type BusLineDirOverview struct {
	LineNo    string            `json:"lineno"`
	Direction int               `json:"linedir"`
	Another   string            `json:"another"`
	StartSN   string            `json:"startsn"`
	EndSn     string            `json:"endsn"`
	Buses     []*api.RunningBus `json:"buses"`
	IsSupport bool              `json:"issupport"`
}

var (
	amapClient *amap.Client
	AMAP_KEY   = `b3abf03fa1e83992727f0625a918fe73`
)

func init() {
	amapClient = amap.NewClient(AMAP_KEY)
	amapClient.HttpClient = api.DEFAULT_HTTP_CLIENT
}

func BusLineSuggest(params martini.Params, r render.Render) {
	lat := params["lat"]
	lon := params["lon"]

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

	bldos := make([]*BusLineDirOverview, 0)
	for _, tip := range resp.Tips {
		sn := strings.TrimRight(tip.Name, "(公交站)")
		linenames := strings.Split(tip.Address, ";")

		var lock sync.Mutex
		var wg sync.WaitGroup
		for _, linename := range linenames {
			wg.Add(1)
			go func(linename string) {
				defer wg.Done()
				bldo := GetBusLineDirOverview(req.City, linename, sn)

				lock.Lock()
				defer lock.Unlock()
				bldos = append(bldos, bldo)
			}(linename)
		}

		wg.Wait()
	}

	r.JSON(200,
		&Response{
			0,
			"OK",
			bldos,
		},
	)
}

func GetBusLineDirOverview(city, linename, station string) (bldo *BusLineDirOverview) {
	lineno := strings.TrimRight(linename, "线")
	lineno = strings.TrimRight(lineno, "路")
	lineno = strings.Replace(lineno, "路内环", "内", 1)
	lineno = strings.Replace(lineno, "路外环", "外", 1)
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
	bldo.StartSN = bdi.StartSn
	bldo.EndSn = bdi.EndSn
	if len(bdi.OtherDirIDs) > 0 {
		bldo.Another = bdi.OtherDirIDs[0]
	}

	//get station index
	var stationIndex int
	for _, s := range bdi.Stations {
		if s.Name == station {
			stationIndex = s.No
			break
		}
	}

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

	return bldo
}
