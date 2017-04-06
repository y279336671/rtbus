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

	rbuses, err := bp.GetRT("青岛", "643", "0")
	if err != nil {
		t.Fatal(err)
	}

	for _, rbus := range rbuses {
		fmt.Printf("%+v\n", rbus)
	}
}
