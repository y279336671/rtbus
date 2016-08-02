package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func Run() {
	//martini
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/rtbus/bj/direction/:linenum", LineNumHandler)
	m.Get("/rtbus/bj/station/:linenum/:direction", LineStationHandler)
	m.Get("/rtbus/bj/bus/:linenum/:direction", LineBusHandler)

	m.RunOnAddr(":1315")
}
