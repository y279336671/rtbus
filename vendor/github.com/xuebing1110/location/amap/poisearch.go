package amap

import (
	"fmt"
	"github.com/xuebing1110/location"
	"strings"
)

const (
	URL_POISEARCH_KEYWORD = `http://restapi.amap.com/v3/place/text`
	URL_POISEARCH_AROUND  = `http://restapi.amap.com/v3/place/around`
)

type PoiSearchResponse struct {
	*ApiResponse
	Pois []*Poi

	Suggestion struct {
		Cities []City `json:"cities"`
	} `json:"suggestion"`
}

type Poi struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	TypeCode     string      `json:"typecode"`
	BizType      string      `json:"biz_type"`
	Address      string      `json:"address"`
	Location     string      `json:"location"`
	Distance     string      `json:"distance"`
	Tel          string      `json:"tel"`
	PCode        string      `json:"pcode"`
	PName        string      `json:"pname"`
	CCode        string      `json:"ccode"`
	CityCode     string      `json:"citycode"`
	CityName     string      `json:"cityname"`
	AdCode       string      `json:"adcode"`
	AdName       string      `json:"adname"`
	EntrLocation string      `json:"entr_location"`
	ExitLocation string      `json:"exit_location"`
	NaviPoiid    string      `json:"navi_poiid"`
	GridCode     string      `json:"gridcode"`
	Alias        string      `json:"alias"`
	BusinessArea string      `json:"business_area"`
	ParkingType  string      `json:"parking_type"`
	IndoorMap    string      `json:"indoor_map"`
	IndoorData   *IndoorData `json:"indoor_data,omitempty"`
}

type City struct {
	Name     string `json:"name"`
	Num      string `json:"num"`
	CityCode string `json:"citycode"`
	AdCode   string `json:"adcode"`
}

type IndoorData struct {
	Cpid      string `json:"cpid"`
	Floor     string `json:"floor"`
	TrueFloor string `json:"truefloor"`
}

type PoiSearchRequest struct {
	*ApiRequest
	url        string
	KeyWords   string `json:"keywords"`
	Types      string `json:"types"`
	City       string `json:"city"`
	Offset     string `json:"offset"`
	Page       string `json:"page"`
	Building   string `json:"building"`
	Floor      string `json:"floor"`
	Output     string `json:"output"`
	Extensions string `json:"extensions"`
	Location   string `json:"location"`
}

func NewPoiSearchRequest(c *Client, keyword string) *PoiSearchRequest {
	return &PoiSearchRequest{
		ApiRequest: &ApiRequest{client: c},
		KeyWords:   keyword,
		Extensions: "all",
		url:        URL_POISEARCH_KEYWORD,
	}
}

func (p *PoiSearchRequest) GetUrlParas() string {
	return GetUrlParas(p.client.key, p)
}

func (p *PoiSearchRequest) Do() (*PoiSearchResponse, error) {
	respobj := &PoiSearchResponse{}

	err := p.do(respobj)
	if err != nil {
		return nil, err
	}

	//未搜索到结果根据建议城市查询
	if len(respobj.Pois) == 0 && len(respobj.Suggestion.Cities) > 0 {
		p.SetCity(respobj.Suggestion.Cities[0].AdCode)
		err := p.do(respobj)
		if err != nil {
			return nil, err
		}
	}

	return respobj, nil
}

func (p *PoiSearchRequest) do(respobj *PoiSearchResponse) error {
	murl := p.url + "?" + p.GetUrlParas()

	err := p.HttpGet(murl, respobj)
	if err != nil {
		return err
	}

	//增加city adcode
	for _, poi := range respobj.Pois {
		if poi.AdCode != "" {
			if location.IsZhixia(poi.AdCode) {
				poi.CCode = poi.AdCode[0:2] + "0000"
			} else {
				poi.CCode = poi.AdCode[0:4] + "00"
			}

		}
	}

	return nil
}

func (p *PoiSearchRequest) SetCity(city string) *PoiSearchRequest {
	p.City = city
	return p
}

func (p *PoiSearchRequest) SetType(poitype string) *PoiSearchRequest {
	p.Types = poitype
	return p
}

func (p *PoiSearchRequest) AddKeword(keywords ...string) *PoiSearchRequest {
	for i, keyword := range keywords {
		if i == 0 && p.KeyWords == "" {
			p.KeyWords = keyword
		} else {
			p.KeyWords += "|" + keyword
		}
	}
	return p
}

func (p *PoiSearchRequest) SetTypes(poitypes []string) *PoiSearchRequest {
	p.Types = strings.Join(poitypes, "|")
	return p
}

func (p *PoiSearchRequest) SetPageSize(size int) *PoiSearchRequest {
	p.Offset = fmt.Sprintf("%d", size)
	return p
}

func (p *PoiSearchRequest) SetPage(page string) *PoiSearchRequest {
	p.Page = page
	return p
}

func (p *PoiSearchRequest) SetBuilding(building string) *PoiSearchRequest {
	p.Building = building
	return p
}

func (p *PoiSearchRequest) SetFloor(floor string) *PoiSearchRequest {
	p.Floor = floor
	return p
}

func (p *PoiSearchRequest) SetOutputJson() *PoiSearchRequest {
	p.Output = "JSON"
	return p
}

func (p *PoiSearchRequest) SetOutputXml() *PoiSearchRequest {
	p.Output = "XML"
	return p
}

func (p *PoiSearchRequest) SetAroundSearch(loc string) *PoiSearchRequest {
	p.url = URL_POISEARCH_AROUND
	p.Location = loc
	return p
}
