package api

import (
	"errors"
	"fmt"
	"github.com/bingbaba/util/httptool"
	"time"
)

const (
	URL_CLL_CITYS_FMT = `http://web.chelaile.net.cn/cdatasource/citylist?type=allRealtimeCity&s=h5&v=3.3.9&userId=browser_%d`
)

type CityInfo struct {
	Code   string `json:"-"`
	ID     string `json:"cityId"`
	Name   string `json:"cityName"`
	Hot    int    `json:"hot"`
	PinYin string `json:"pinyin"`
	Subway int    `json:"supportSubway"`
}

type AllCityResp struct {
	Status string `json:"status"`
	Data   struct {
		AllRealtimeCity []*CityInfo `json:"allRealtimeCity"`
	} `json:"data"`
}

func GetCllAllCitys() ([]*CityInfo, error) {
	reqUrl := fmt.Sprintf(URL_CLL_CITYS_FMT, time.Now().UnixNano()/1000000)
	httreq, err := getCllHttpRequest(reqUrl)
	if err != nil {
		return nil, err
	}

	cllresp := &AllCityResp{}
	err = httptool.HttpDoJson(httreq, cllresp)
	if err != nil {
		return nil, err
	}

	if cllresp.Status != "OK" {
		return nil, errors.New(cllresp.Status)
	}

	return cllresp.Data.AllRealtimeCity, nil
}
