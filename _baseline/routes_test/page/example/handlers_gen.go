package example

import (
	"net/http"

	"github.com/rur/good/baseline/routes_test/service"
	"github.com/rur/treetop"
)

// -------------------------
// example Handlers
// -------------------------

// Ref: settings
// Block: content
// Method:
// Doc: other Content Page
func settingsHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
		Settings    interface{}
		Tabs        interface{}
	}{
		HandlerInfo: "example Page settingsHandler",
		Settings:    rsp.HandleSubView("settings", req),
		Tabs:        rsp.HandleSubView("tabs", req),
	}
	return data
}

// Ref: general-settings-tab
// Block: settings
// Method:
// Doc: General settings area
func generalSettingsTabHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page generalSettingsTabHandler",
	}
	return data
}

// Ref: advanced-settings-tab
// Block: settings
// Method:
// Doc: Advanced settings area
func advancedSettingsTabHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page advancedSettingsTabHandler",
	}
	return data
}

// Ref: other-tabs
// Block: tabs
// Method:
// Doc: Tabs for the contentent in the content page
func otherTabsHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page otherTabsHandler",
	}
	return data
}
