package example

import (
	"net/http"

	"github.com/rur/good/baseline/routes_test/service"
	"github.com/rur/treetop"
)

// -------------------------
// example Handlers
// -------------------------

// Ref: placeholder-form
// Block: form
// Method: POST
// Doc: Placeholder form
func placeholderFormHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
		FormError   interface{}
	}{
		HandlerInfo: "example Page placeholderFormHandler",
		FormError:   rsp.HandleSubView("form-error", req),
	}
	return data
}

// Ref: basic-form-error
// Block: form-error
// Doc: Format and display a form error message
func basicFormErrorHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page basicFormErrorHandler",
	}
	return data
}

// Ref: alternative-content
// Block: content
// Doc: Alaternative Content Page
func alternativeContentHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page alternativeContentHandler",
	}
	return data
}

// Ref: settings-layout
// Block: content
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
// Doc: Advanced settings area
func advancedSettingsHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo  string
		SettingsForm interface{}
	}{
		HandlerInfo:  "example Page advancedSettingsHandler",
		SettingsForm: rsp.HandleSubView("settings-form", req),
	}
	return data
}

// Ref: update-advanced-settings
// Block: settings-form
// Method: POST
// Doc: Accept update to advanced settings and show result
func updateAdvancedSettingsHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
		FormError   interface{}
	}{
		HandlerInfo: "example Page updateAdvancedSettingsHandler",
		FormError:   rsp.HandleSubView("form-error", req),
	}
	return data
}
