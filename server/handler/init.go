package handler

import (
	"github.com/bingbaba/util/logs"
	"github.com/xuebing1110/rtbus/api"
)

var (
	BusTool *api.BusPool
	logger  *logs.Blogger
)

func init() {
	logger = logs.GetBlogger()
	BusTool = api.NewBusPoolAsync()
}

type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}
