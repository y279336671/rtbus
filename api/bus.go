package api

import (
	"errors"
	"fmt"
	"sync"
)

const (
	BUS_ARRIVING_STATUS        = "1"
	BUS_ARRIVING_FUTURE_STATUS = "0.5"
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

	getRT func(bl *BusLine, dirname string) ([]*RunningBus, error)
}

type BusDirInfo struct {
	l         sync.Mutex
	freshTime int64

	ID           string        `json:"id"`
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
	Status   string  `json:"status"`
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
		bp.CityBusLines[cityName] = cllbls
	}

	//BeiJing
	var bjbls *CityBusLines
	bjbls, err = GetAiBangAllLine()
	if err != nil {
		return
	}
	bp.CityBusLines["北京"] = bjbls

	return
}

func (bp *BusPool) GetRT(city, linenum, dirname string) (rbus []*RunningBus, err error) {
	cbl, found := bp.CityBusLines[city]
	if !found {
		err = errors.New(fmt.Sprintf("can't support the city %s", city))
		return
	}

	if cbl.Source == SOURCE_CHELAILE {
		return GetCllLineRT(cbl, linenum, dirname)
	} else {
		bl, found := cbl.ByLineName[linenum]
		if !found {
			err = errors.New(fmt.Sprintf("can't find the line %s in city %s", linenum, city))
			return
		}
		return bl.getRT(bl, dirname)
	}

	return
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
	for bdi_name, bdi := range bl.Directions {
		//fmt.Printf("%+v\n", bdi)
		if dirname == fmt.Sprintf("%d", bdi.Direction) || dirname == bdi_name {
			return bdi, true
		}
	}

	return nil, false
}

func (bdi *BusDirInfo) GetDirName() string {
	return fmt.Sprintf("%s-%s", bdi.Stations[0].Name, bdi.Stations[len(bdi.Stations)-1].Name)
}
