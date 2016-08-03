package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func Run() {
	//martini
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/rtbus/bj/direction/:linenum", BJBusLineHandler)
	m.Get("/rtbus/bj/station/:linenum/:direction", BJBusSnHandler)
	m.Get("/rtbus/bj/bus/:linenum/:direction", BJBusSnBusHandler)

	m.Get(`/rtbus/:city/direction/:linenum`, CllBusLineHandler)
	m.Get(`/rtbus/:city/station/:linenum/:direction`, CllBusSnHandler)
	m.Get(`/rtbus/:city/bus/:linenum/:direction`, CllBusSnBusHandler)

	m.RunOnAddr(":1315")
}
