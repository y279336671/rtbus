package api

import (
	"github.com/bingbaba/util/logs"
	"testing"
)

func TestBjBus(t *testing.T) {
	logs.SetDebug(true)

	bus, err := NewBJBusSess()
	if err != nil {
		t.Fatal(err)
	}

	err = bus.LoadBusLineConf("675")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1; i++ {
		busline := bus.BusLines["675"]

		for _, dir := range busline.Direction {
			err = bus.FreshStatus(busline.LineNum, dir.ID)
			if err != nil {
				t.Error(err)
			}
		}

		bus.Print()
	}

}
