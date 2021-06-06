package example

import (
	"net/http"

	"github.com/rur/good/baseline/routes_test/service"
	"github.com/rur/treetop"
)

// -------------------------
// example Handlers
// -------------------------

// Ref: settings-layout
// Block: content
// Method:
// Doc: Settings page layout
func settingsLayoutHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
		Settings    interface{}
		Tabs        interface{}
	}{
		HandlerInfo: "example Page settingsLayoutHandler",
		Settings:    rsp.HandleSubView("settings", req),
		Tabs:        rsp.HandleSubView("tabs", req),
	}
	return data
}

// Ref: general-settings
// Block: settings
// Method:
// Doc: General settings area
func generalSettingsHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page generalSettingsHandler",
	}
	return data
}

// Ref: advanced-settings
// Block: settings
// Method:
// Doc: Advanced settings area
func advancedSettingsHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page advancedSettingsHandler",
	}
	return data
}

// Ref: settings-tabs
// Block: tabs
// Method:
// Doc: Tabs for the settings page content
func settingsTabsHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page settingsTabsHandler",
	}
	return data
}
