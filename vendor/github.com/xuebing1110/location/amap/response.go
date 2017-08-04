package amap

type Response interface {
	GetStatus() string
	GetInfo() string
}

type ApiResponse struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Count  string `json:"count"`
}

func (resp *ApiResponse) GetStatus() string {
	return resp.Status
}

func (resp *ApiResponse) GetInfo() string {
	return resp.Info
}
