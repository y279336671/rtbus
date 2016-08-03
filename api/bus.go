package api

import (
	"errors"
	"sync"
)

type BusLine struct {
	LineNum   string        `json:"linenum"`
	Direction []*BusDirInfo `json:"direction"`
}

type BusDirInfo struct {
	l          sync.Mutex
	freshTime  int64
	Name2Index map[string]int `json:"-"`

	ID        string        `json:"id"`
	Direction int           `json:"direction,omitempty"`
	Name      string        `json:"name"`
	StartSn   string        `json:"startsn,omitempty"`
	EndSn     string        `json:"endsn,omitempty"`
	Price     string        `json:"price,omitempty"`
	SnNum     int           `json:"stationsNum,omitempty"`
	FirstTime string        `json:"firstTime,omitempty"`
	LastTime  string        `json:"lastTime,omitempty"`
	Stations  []*BusStation `json:"stations"`
}

type BusStation struct {
	Order int           `json:"order"`
	Sn    string        `json:"sn,omitempty"`
	Buses []*RunningBus `json:"buses,omitempty"`
}

type RunningBus struct {
	Order    int     `json:"order"`
	Status   string  `json:"status"`
	BusID    string  `json:"busid,omitempty"`
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	Distance int     `json:"distanceToSc,omitempty"`
	SyncTime int64   `json:"syncTime,omitempty"`
}

func NewBusLine(lineid string) *BusLine {
	return &BusLine{
		LineNum:   lineid,
		Direction: make([]*BusDirInfo, 0),
	}
}

func (b *BusLine) GetBusInfo(dirid string) (*BusDirInfo, error) {
	for _, busdir := range b.Direction {
		if busdir.equal(dirid) {
			return busdir, nil
		}
	}

	return nil, errors.New("can't found the direction:" + dirid)
}

func (b *BusLine) GetRunningBus(dirid string) ([]*BusStation, error) {
	rbuses := make([]*BusStation, 0)

	busdir, err := b.GetBusInfo(dirid)
	if err != nil {
		return rbuses, err
	}

	for _, station := range busdir.Stations {
		if len(station.Buses) > 0 {
			rbuses = append(rbuses, station)
		}
	}

	return rbuses, nil
}

func (s *BusDirInfo) getSnDesc() string {
	return s.StartSn + "-" + s.EndSn
}

func (d *BusDirInfo) equal(dirid string) bool {
	if d.ID == dirid ||
		d.Name == dirid ||
		dirid == d.getSnDesc() {
		return true
	} else {
		return false
	}
}
