package api

import (
	"github.com/bingbaba/util/logs"
	"testing"
)

func TestCllBus(t *testing.T) {
	logs.SetDebug(true)

	bus, err := NewCllBus("0532")
	if err != nil {
		t.Fatal(err)
	}

	bl, err := bus.GetBusLine("318")
	if err != nil {
		t.Fatal(err)
	}

	if len(bl.Direction) == 0 {
		t.Fatal("can't get the directions of the 318 road")
	}

}
