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
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Stations   []*BusStation  `json:"stations"`
}

type BusStation struct {
	ID    int           `json:"id"`
	Name  string        `json:"name,omitempty"`
	Buses []*RunningBus `json:"buses,omitempty"`
}

type RunningBus struct {
	StationID int     `json:"sid"`
	Status    string  `json:"status"`
	BusID     string  `json:"busid,omitempty"`
	Lat       float64 `json:"lat,omitempty"`
	Lng       float64 `json:"lng,omitempty"`
	Distance  int     `json:"distance,omitempty"`
}
