package controllers

import (
	"github.com/revel/revel"
	logit "logMonitor/app/services/logMonitor"
)

// App - struct
type App struct {
	*revel.Controller
}

// Index - ctrl route to Index.html
func (c App) Index() revel.Result {
	logit.INFO("get-запрос на Index", "основной поток", "")
	return c.Render()
}
