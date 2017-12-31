package logs

import (
	"github.com/astaxie/beego/logs"
	"io/ioutil"
)

var (
	DEBUG  bool
	logger *logs.BeeLogger
)

type Blogger struct {
	*logs.BeeLogger
}

func init() {
	logger = logs.NewLogger(1)
	defautInit()
}

func GetBlogger() *Blogger {
	return &Blogger{logger}
}

func defautInit() {
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(2)

	if DEBUG {
		logger.SetLevel(logs.LevelDebug)
		logger.SetLogger("console", ``)
	} else {
		logger.SetLevel(logs.LevelInfo)
		logger.DelLogger("console")
	}
}

func SetDebug(debug bool) {
	DEBUG = debug
	defautInit()
}

func Init(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return InitByString(string(content))
}

func InitByString(content string) error {

	defautInit()
	logger.SetLogger("file", content)
	logger.Info("startint...")

	return nil
}
