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

	var err error
	BusTool, err = api.NewBusPool()
	if err != nil {
		Logger.Error("%v", err)
	}
}

type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}
