package api

import (
	//"fmt"
	"net/http"
)

func getToken() ([]*http.Cookie, error) {
	resp, err := http.Get(URL_BJ_HOME)
	if err != nil {
		return make([]*http.Cookie, 0), err
	}

	// //test
	// for _, c := range resp.Cookies() {
	// 	fmt.Println(c.Raw)
	// }

	return resp.Cookies(), nil

}
