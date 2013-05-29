package config

import (
	"appengine"
)

type Config struct {
	IsDevApp      bool
	Title         string
	AnalyticsId   string
	AnalyticsName string
	LivefyreId    string
	TwitterId     string
	TwitterWidget string
}

var Default = &Config{}

func InitDefault() {

	Default.IsDevApp = appengine.IsDevAppServer()

	Default.Title = "GoPaper"

	//Default.AnalyticsId = "UA-xyz-1"
	//Default.AnalyticsName = ".com"

	/*if !Default.IsDevApp {
		Default.LivefyreId = "xyz"
	}*/

	//Default.TwitterId = "akhansari"
	//Default.TwitterWidget = "xyz"
}
