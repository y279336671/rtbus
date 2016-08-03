package api

import (
	"errors"
	"sync"
)

type CllBusPool struct {
	l            sync.RWMutex
	CllCityBuses map[string]*CllBus
}

var (
	MAP_CITY = map[string]CityInfo{
		"0532": CityInfo{
			CityID:  "009",
			Name:    "青岛",
			TelCode: "0532",
		},
	}
)

func NewCllBusPool() *CllBusPool {
	cbp := new(CllBusPool)
	cbp.initCllBusPool()

	return cbp
}

func (p *CllBusPool) initCllBusPool() {
	p.CllCityBuses = make(map[string]*CllBus)
	for citytel, _ := range MAP_CITY {
		cllbus, _ := NewCllBus(citytel)
		p.CllCityBuses[citytel] = cllbus
	}
}

func (p *CllBusPool) GetCllBus(citytel string) (*CllBus, error) {
	if p.CllCityBuses == nil {
		p.initCllBusPool()
	}

	cllbus, found := p.CllCityBuses[citytel]
	if !found {
		return cllbus, errors.New("can't support the city:" + citytel)
	} else {
		return cllbus, nil
	}
}
