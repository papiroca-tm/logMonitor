package controllers

import (
	"github.com/revel/revel"
	logit "logMonitor/app/services/logMonitor"
)

// LogMonitor ...
type LogMonitor struct {
	*revel.Controller
}

// Index ...
func (c LogMonitor) Index() revel.Result {
	logit.INFO("get-запрос на Index монитора логов", "просмотр логов", "")
	return c.Render()
}

// GetLogs ...
func (c LogMonitor) GetLogs() revel.Result {
	logit.INFO("get-запрос на GetLogs монитора логов", "просмотр логов", "")
	return c.RenderJson(logit.Get())
}