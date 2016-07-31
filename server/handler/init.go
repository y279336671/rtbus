package handler

import (
	"github.com/bingbaba/util/logs"
	"github.com/xuebing1110/rtbus/api"
)

var (
	BusSess *api.BJBusSess
	logger  *logs.Blogger
)

func init() {
	logger = logs.GetBlogger()

	var err error
	BusSess, err = api.NewBJBusSess()
	if err != nil {
		panic(err)
	}
}

type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}
