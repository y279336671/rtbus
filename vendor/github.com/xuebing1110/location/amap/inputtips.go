package amap

import (
	// "github.com/xuebing1110/location"
	"fmt"
	"strings"
)

const (
	URL_INPUTTIPS = `http://restapi.amap.com/v3/assistant/inputtips`
)

type InputtipsResponse struct {
	*ApiResponse
	Tips []*Tip

	Suggestion struct {
		Cities []City `json:"cities"`
	} `json:"suggestion"`
}

type Tip struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	District string `json:"district"`
	AdCode   string `json:"adcode"`
	Location string `json:"location"`
	Address  string `json:"address"`
}

type InputtipsRequest struct {
	*ApiRequest
	url       string
	KeyWords  string `json:"keywords"`
	Types     string `json:"types"`
	Location  string `json:"location"`
	City      string `json:"city"`
	CityLimit string `json:"citylimit"`
	DataType  string `json:"datatype"`
	Output    string `json:"output"`
	CityName  string `json:"-"`
}

func NewInputtipsRequest(c *Client, keyword string) *InputtipsRequest {
	return &InputtipsRequest{
		ApiRequest: &ApiRequest{client: c},
		KeyWords:   keyword,
		url:        URL_INPUTTIPS,
	}
}

func (p *InputtipsRequest) GetUrlParas() string {
	return GetUrlParas(p.client.key, p)
}

func (p *InputtipsRequest) Do() (*InputtipsResponse, error) {
	if p.CityLimit == "true" && p.City == "" && p.Location != "" {
		req := NewPoiSearchRequest(p.client, "").SetAroundSearch(p.Location).SetPageSize(2)
		resp, err := req.Do()
		if err != nil {
			return nil, err
		}
		if len(resp.Pois) > 0 {
			p.City = resp.Pois[0].CityCode
			p.CityName = resp.Pois[0].CityName
		}
	}

	respobj := &InputtipsResponse{}
	err := p.do(respobj)
	if err != nil {
		return nil, err
	}

	//remove unuse data
	if len(respobj.Tips) > 0 {
		if respobj.Tips[0].ID == "" {
			respobj.Tips = respobj.Tips[1:]
		}
	}

	return respobj, nil
}

func (p *InputtipsRequest) do(respobj *InputtipsResponse) error {
	murl := p.url + "?" + p.GetUrlParas()

	err := p.HttpGet(murl, respobj)
	if err != nil {
		return err
	}

	return nil
}

func (p *InputtipsRequest) SetCity(city string) *InputtipsRequest {
	p.City = city
	p.CityName = city
	return p
}

func (p *InputtipsRequest) AddKeword(keywords ...string) *InputtipsRequest {
	for i, keyword := range keywords {
		if i == 0 && p.KeyWords == "" {
			p.KeyWords = keyword
		} else {
			p.KeyWords += "|" + keyword
		}
	}
	return p
}

func (p *InputtipsRequest) SetTypes(poitypes []string) *InputtipsRequest {
	p.Types = strings.Join(poitypes, "|")
	return p
}

func (p *InputtipsRequest) SetOutputJson() *InputtipsRequest {
	p.Output = "JSON"
	return p
}

func (p *InputtipsRequest) SetOutputXml() *InputtipsRequest {
	p.Output = "XML"
	return p
}

func (p *InputtipsRequest) SetCityLimit(limit bool) *InputtipsRequest {
	if limit {
		p.CityLimit = "true"
	} else {
		p.CityLimit = "false"
	}
	return p
}
func (p *InputtipsRequest) SetLocation(loc string) *InputtipsRequest {
	p.Location = loc
	return p
}
func (p *InputtipsRequest) SetLatLon(lat, lon string) *InputtipsRequest {
	p.Location = fmt.Sprintf("%s,%s", lon, lat)
	return p
}
