package page

import (
	"net/http"
	"strings"

	"github.com/rur/good/_baseline/site/service"
	"github.com/rur/treetop"
)

// ViewHandlerWithEnv is the signature of a treetop view handler with an additional
// Env type parameter, used
type ViewHandlerWithEnv func(*service.Env, treetop.Response, *http.Request) interface{}

// Helper is used by the page Routes function to perform the following tasks
type Helper interface {
	BindEnv(ViewHandlerWithEnv) treetop.ViewHandlerFunc
	HandleGET(string, http.Handler)
	HandlePOST(string, http.Handler)
	Handle(string, http.Handler)
}

// DefaultHelper implements the Helper interface using
// the standard server mux, you may wish to add your own depending
// on your requirements. If so, you can simply modify the following
// code to suit your needs
type DefaultHelper struct {
	Env *service.Env
	// EDITME: if you wish to use a different router library
	Mux *http.ServeMux
}

// BindEnv will inject the sever environment to a treetop request handler using closures
func (hlp *DefaultHelper) BindEnv(fn ViewHandlerWithEnv) treetop.ViewHandlerFunc {
	return func(rsp treetop.Response, req *http.Request) interface{} {
		return fn(hlp.Env, rsp, req)
	}
}

// HandleGET registers a suppled handler with a pattern accepting GET/HEAD requests only,
// the semantics of how the pattern are handled depends on your router.
func (hlp *DefaultHelper) HandleGET(pattern string, hdr http.Handler) {
	hlp.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		switch strings.ToUpper(r.Method) {
		case "GET", "HEAD":
			hdr.ServeHTTP(w, r)

		case "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE":
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	})
}

// HandlePOST registers a suppled handler with a pattern accepting POST requests only,
// the semantics of how the pattern are handled depends on your router.
func (hlp *DefaultHelper) HandlePOST(pattern string, hdl http.Handler) {
	hlp.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		switch strings.ToUpper(r.Method) {
		case "POST":
			hdl.ServeHTTP(w, r)

		case "GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE":
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	})
}

// Handle will register a pattern without filtering requests by method, the
// request handler is responsible for differentiating
func (hlp *DefaultHelper) Handle(pattern string, hdl http.Handler) {
	hlp.Mux.Handle(pattern, hdl)
}
