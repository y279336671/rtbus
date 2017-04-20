package api

import (
	"fmt"
	"testing"
)

func TestGetBJBuses(t *testing.T) {
	bl, err := GetBJBusLine("675")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", ToJsonString(bl))

	bdi, err := GetBJBusRT("675", "0")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s\n", ToJsonString(bdi))
}
