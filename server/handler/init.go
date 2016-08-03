package handler

import (
	"github.com/bingbaba/util/logs"
	"github.com/xuebing1110/rtbus/api"
)

var (
	BjBusSess   *api.BJBusSess
	DCllBusPool *api.CllBusPool
	logger      *logs.Blogger
)

func init() {
	logger = logs.GetBlogger()

	var err error
	BjBusSess, err = api.NewBJBusSess()
	if err != nil {
		panic(err)
	}

	DCllBusPool = api.NewCllBusPool()
}

type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}
