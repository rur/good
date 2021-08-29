package newpage

import (
	"net/http"

	"github.com/rur/good/baseline/page_test/service"
	"github.com/rur/treetop"
)

// Ref: base
// Doc: Base HTML template for newpage page
func baseHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		PageTitle string
		Content   interface{}
		Nav       interface{}
		Scripts   interface{}
		Styles    interface{}
		Modal     interface{}
	}{
		PageTitle: "newpage Page",
		Content:   rsp.HandleSubView("content", req),
		Nav:       rsp.HandleSubView("nav", req),
		Scripts:   rsp.HandleSubView("scripts", req),
		Styles:    rsp.HandleSubView("styles", req),
		Modal:     rsp.HandleSubView("modal", req),
	}
	return data
}

// Ref: content-wrapper
// Block: content
// Method: GET
// Doc: Landing content view for the newpage page URI
func contentWrapperHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		SubSection      interface{}
		QuickSearchMenu interface{}
	}{
		SubSection:      rsp.HandleSubView("subsection", req),
		QuickSearchMenu: rsp.HandleSubView("quick-search-menu", req),
	}
	return data
}

// Ref: sidebar-nav
// Block: nav
func sidebarNavHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		User     service.User
		NavLinks interface{}
	}{
		User:     rsc.user,
		NavLinks: rsp.HandleSubView("nav-links", req),
	}
	return data
}

// Ref: sidebar-nav-links
// Block: nav-links
func sidebarNavLinksHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	type NavItem struct {
		Title    string
		Href     string
		Selected bool
		Partial  bool
	}
	return struct {
		Items []NavItem
	}{
		Items: []NavItem{
			{
				Title: "Intro",
				Href:  "/",
			},
			{
				Title:    "newpage",
				Href:     "/newpage",
				Selected: true,
			},
			// EDITME: To reduce duplication between pages, you might want to move main navigation handlers
			//         and template files to the shared [site]/page folder
		},
	}
}
