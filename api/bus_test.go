package api

import (
	"fmt"
	"testing"
)

func TestBjAiBangBus(t *testing.T) {
	bp, err := NewBusPool()
	if err != nil {
		t.Fatal(err)
	}

	rbuses, err := bp.GetRT("北京", "675", "0")
	if err != nil {
		t.Fatal(err)
	}

	for _, rbus := range rbuses {
		fmt.Printf("%+v\n", rbus)
	}
}

func TestQingDaoCllBus(t *testing.T) {
	bp, err := NewBusPool()
	if err != nil {
		t.Fatal(err)
	}

	rbuses, err := bp.GetRT("青岛", "318", "0")
	if err != nil {
		t.Fatal(err)
	}

	for _, rbus := range rbuses {
		fmt.Printf("%+v\n", rbus)
	}
}

/*
func TestGetBusStation(t *testing.T) {
	bp, err := NewBusPool()
	if err != nil {
		t.Fatal(err)
	}

	bss, err := bp.GetStations("青岛", "318", "0")
	if err != nil {
		t.Fatal(err)
	}

	for _, bs := range bss {
		fmt.Printf("%+v\n", bs)
	}

	bss, err = bp.GetStations("北京", "675", "0")
	if err != nil {
		t.Fatal(err)
	}

	for _, bs := range bss {
		fmt.Printf("%+v\n", bs)
	}
}
*/
