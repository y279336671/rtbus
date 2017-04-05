package api

import (
	"errors"
	"sync"
	"time"
)

const (
	BUS_ARRIVING_STATUS        = "0.5"
	BUS_ARRIVING_FUTURE_STATUS = "1"
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

type RefreshBuslineDir interface {
	freshBuslineDir(lineid, dirid string) error
}

func NewBusLine(lineid string) *BusLine {
	return &BusLine{
		LineNum:   lineid,
		Direction: make([]*BusDirInfo, 0),
	}
}

func (b *BusLine) GetBusDir(dirid string, r RefreshBuslineDir) (*BusDirInfo, error) {
	busdir, err := b.getBusDir(dirid)
	if err != nil {
		return nil, err
	}

	busdir.l.Lock()
	defer busdir.l.Unlock()

	//无需重新加载 仅更新同步时间即可
	curtime := time.Now().Unix()
	if curtime-busdir.freshTime < 10 {
		for _, s := range busdir.Stations {
			for _, rbus := range s.Buses {
				rbus.SyncTime = rbus.SyncTime + (curtime - busdir.freshTime)
			}
		}
	} else {
		err = r.freshBuslineDir(b.LineNum, dirid)
		if err != nil {
			return nil, err
		}
	}

	return busdir, nil
}

func (b *BusLine) getBusDir(dirid string) (*BusDirInfo, error) {
	for _, busdir := range b.Direction {
		if busdir.equal(dirid) {
			return busdir, nil
		}
	}

	return nil, errors.New("can't found the direction:" + dirid)
}

func (b *BusLine) GetRunningBus(dirid string, r RefreshBuslineDir) ([]*BusStation, error) {
	rbuses := make([]*BusStation, 0)

	busdir, err := b.GetBusDir(dirid, r)
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
