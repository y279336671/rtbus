package api

import (
	"errors"
	"fmt"
	"github.com/xuebing1110/location"
	"sync"
)

const (
	BUS_ARRIVING_STATUS        = 1
	BUS_ARRIVING_FUTURE_STATUS = 0.5
	SOURCE_CHELAILE            = "cll"
)

type BusPool struct {
	CityBusLines map[string]*CityBusLines
}

type CityBusLines struct {
	l          sync.Mutex
	Source     string
	CityInfo   *CityInfo
	ByLineName map[string]*BusLine
}

type BusLine struct {
	l          sync.Mutex
	LineNum    string                 `json:"linenum"`
	LineName   string                 `json:"lineName"`
	Directions map[string]*BusDirInfo `json:"direction"`
}

type BusDirInfo struct {
	l         sync.Mutex
	freshTime int64

	ID           string        `json:"id"`
	OtherDirIDs  []string      `json:"otherDirIds"`
	Direction    int           `json:"direction,omitempty"`
	Name         string        `json:"name"`
	StartSn      string        `json:"startsn,omitempty"`
	EndSn        string        `json:"endsn,omitempty"`
	Price        string        `json:"price,omitempty"`
	SnNum        int           `json:"stationsNum,omitempty"`
	FirstTime    string        `json:"firstTime,omitempty"`
	LastTime     string        `json:"lastTime,omitempty"`
	Stations     []*BusStation `json:"stations"`
	RunningBuses []*RunningBus `json:"buses,omitempty"`
}

type BusStation struct {
	No   int     `json:"order"`
	Name string  `json:"sn,omitempty"`
	Lat  float64 `json:"lat,omitempty"`
	Lon  float64 `json:"lon,omitempty"`
}

type RunningBus struct {
	No       int     `json:"order"`
	Name     string  `json:"-"`
	Status   float64 `json:"status"`
	BusID    string  `json:"busid,omitempty"`
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	Distance int     `json:"distanceToSc,omitempty"`
	SyncTime int64   `json:"syncTime,omitempty"`
}

func NewBusPool() (bp *BusPool, err error) {
	bp = &BusPool{
		CityBusLines: make(map[string]*CityBusLines),
	}

	//CheLaiLe
	var cll_cbls []*CityBusLines
	cll_cbls, err = GetCllAllBusLine()
	if err != nil {
		return
	}
	for _, cllbls := range cll_cbls {
		cityName := cllbls.CityInfo.Name
		citycode := location.GetCitycode(location.MustParseCity(cityName))
		if citycode == "" {
			bp.CityBusLines[cityName] = cllbls
		} else {
			bp.CityBusLines[citycode] = cllbls
		}
	}

	//BeiJing
	var bjbls *CityBusLines
	bjbls, err = GetAiBangAllLine()
	if err != nil {
		return
	}
	citycode := location.GetCitycode(location.MustParseCity("北京"))
	bp.CityBusLines[citycode] = bjbls

	return
}

func NewBusPoolAsync() (bp *BusPool) {
	bp = &BusPool{
		CityBusLines: make(map[string]*CityBusLines),
	}

	//CheLaiLe
	cll_cbls, err := GetCllAllBusLine()
	if err != nil {
		LOGGER.Error("%v", err)
		return
	} else {
		for _, cllbls := range cll_cbls {
			cityName := cllbls.CityInfo.Name
			citycode := location.GetCitycode(location.MustParseCity(cityName))
			if citycode == "" {
				bp.CityBusLines[cityName] = cllbls
			} else {
				bp.CityBusLines[citycode] = cllbls
			}
		}
	}

	//BeiJing
	go func() {
		bjbls, err := GetAiBangAllLine()
		if err != nil {
			LOGGER.Error("%v", err)
			return
		}
		citycode := location.GetCitycode(location.MustParseCity("北京"))
		bp.CityBusLines[citycode] = bjbls
	}()

	return
}

func (bp *BusPool) GetBusLineInfo(city, lineno string) (bl *BusLine, err error) {
	//check wether support the city
	cbl, found := bp.CityBusLines[city]
	if !found {
		city = location.GetCitycode(location.MustParseCity(city))
		cbl, found = bp.CityBusLines[city]
	}
	if !found {
		err = errors.New(fmt.Sprintf("can't support the city %s", city))
		return
	}

	if cbl.Source == SOURCE_CHELAILE {
		return getCllBusLine(cbl, lineno)
	} else {
		return cbl.getBusLine(lineno)
	}
}

func (bp *BusPool) GetBusLineDirInfo(city, lineno, dirname string) (bdi *BusDirInfo, err error) {
	//get bus line
	var bl *BusLine
	bl, err = bp.GetBusLineInfo(city, lineno)
	if err != nil {
		return
	}

	//check wether find the direction
	var found bool
	bdi, found = bl.GetBusDirInfo(dirname)
	if !found {
		err = errors.New(fmt.Sprintf("can't find %s(%s) in city %s", lineno, dirname, city))
		return
	}

	//return
	return bdi, nil
}

func (bp *BusPool) GetRT(city, lineno, dirname string) (rbus []*RunningBus, err error) {
	cbl, found := bp.CityBusLines[city]
	if !found {
		city = location.GetCitycode(location.MustParseCity(city))
		cbl, found = bp.CityBusLines[city]
	}
	if !found {
		err = errors.New(fmt.Sprintf("can't support the city %s", city))
		return
	}

	if cbl.Source == SOURCE_CHELAILE {
		return GetCllLineRT(cbl, lineno, dirname)
	} else {
		var bl *BusLine
		bl, err = cbl.getBusLine(lineno)
		if err != nil {
			return
		}
		return GetAiBangLineRT(bl, dirname)
	}

	return
}

func (cbl *CityBusLines) hasInit(lineno string) bool {
	_, found := cbl.ByLineName[lineno]
	return found
}

func (cbl *CityBusLines) getBusLine(lineno string) (*BusLine, error) {
	bl, found := cbl.ByLineName[lineno]
	if !found {
		return nil, errors.New(fmt.Sprintf("can't get the %s line %s info!", cbl.CityInfo.Name, lineno))
	} else {
		return bl, nil
	}

}

func NewCityBusLines() *CityBusLines {
	return &CityBusLines{
		ByLineName: make(map[string]*BusLine),
	}
}

func (bl *BusLine) Put(bdi *BusDirInfo) {
	if bl == nil {
		return
	}

	bl.l.Lock()
	defer bl.l.Unlock()
	bl.Directions[bdi.Name] = bdi
}

func (bl *BusLine) GetBusDirInfo(dirname string) (*BusDirInfo, bool) {
	for dirkey, bdi := range bl.Directions {
		//fmt.Printf("%+v\n", bdi)
		if dirname == fmt.Sprintf("%d", bdi.Direction) ||
			dirname == bdi.GetDirName() ||
			dirkey == dirname {
			return bdi, true
		}
	}

	return nil, false
}

func (bdi *BusDirInfo) GetDirName() string {
	return fmt.Sprintf("%s-%s", bdi.StartSn, bdi.EndSn)
}
