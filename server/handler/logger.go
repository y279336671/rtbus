package handler

import (
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"os"
	"time"
)

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
func Logger() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		logger.Info("Started %s %s for %s", req.Method, req.URL.Path, addr)

		rw := res.(martini.ResponseWriter)
		c.Next()

		logger.Info("Completed %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}

func GetNilLogger() *log.Logger {
	f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if e == nil {
		return log.New(f, "", 0)
	} else {
		return log.New(os.Stdout, "prefix", 0)
	}
}

func SetNilLogger(l *log.Logger) {
	f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if e == nil {
		l.SetOutput(f)
	}
}
