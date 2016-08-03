package api

import (
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
	Direction  string         `json:"direction"`
	Name       string         `json:"name"`
	StartSn    string         `json:"startsn,omitempty"`
	EndSn      string         `json:"endsn,omitempty"`
	Price      string         `json:"price,omitempty"`
	SnNum      int            `json:"stationsNum,omitempty"`
	FirstTime  string         `json:"firstTime,omitempty"`
	LastTime   string         `json:"lastTime,omitempty"`
	Stations   []*BusStation  `json:"stations"`
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
	Distance int     `json:"distance,omitempty"`
	SyncTime int     `json:"syncTime,omitempty"`
}
