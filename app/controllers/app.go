package controllers

import (
	"github.com/revel/revel"
)

// App - struct
type App struct {
	*revel.Controller
}

// Index - ctrl route to Index.html
func (c App) Index() revel.Result {
	return c.Render()
}
