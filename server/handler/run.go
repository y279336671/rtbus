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
	m.Get("/rtbus/bj/info/:linenum/:direction", LineInfoHandler)

	m.RunOnAddr(":1315")
}
