package handler

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func Run(port int) {
	//martini
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get(`/rtbus/:city/:linenum`, BusLineHandler)
	m.Get(`/rtbus/:city/:linenum/:direction`, BusDirHandler)
	m.Get(`/rtbus/:city/:linenum/:direction/bus`, RunningBusHandler)

	m.Use(Logger())
	m.Map(GetNilLogger())

	m.RunOnAddr(fmt.Sprintf(":%d", port))
}
