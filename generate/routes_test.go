package generate

import (
	"reflect"
	"testing"

	"github.com/rur/good/routemap"
)

func TestTemplateDataFromRoutes(t *testing.T) {
	def := routemap.PageRoutes{
		Namespace: "github.com/rur/example/site/page/mypage",
		URI:       "/my-page",
		RouteView: routemap.RouteView{
			Ref:      "mypage",
			Template: "page/mypage/templates/mypage.html.tmpl",
			Handler:  "hlp.BindEnv(bindResources(mypageHandler))",
			Doc:      "Test page docs",
		},
	}

	gotHandlers, gotEntries, gotRoutes, err := TemplateDataFromRoutes(def)

	wantHandlers := []Handler{
		{
			Ref:        "mypage",
			Type:       "PageView",
			Doc:        "Test page docs",
			Identifier: "mypageHandler",
		},
	}
	wantEntries := []Entry{}
	wantRoutes := []Route{}

	if err != nil {
		t.Errorf("TemlateDataFromRoutes() error = %v", err)
		return
	}

	if !reflect.DeepEqual(gotHandlers, wantHandlers) {
		t.Errorf("TemlateDataFromRoutes() gotHandlers = %v, want %v", gotHandlers, wantHandlers)
	}
	if !reflect.DeepEqual(gotEntries, wantEntries) {
		t.Errorf("TemlateDataFromRoutes() gotEntries = %v, want %v", gotEntries, wantEntries)
	}
	if !reflect.DeepEqual(gotRoutes, wantRoutes) {
		t.Errorf("TemlateDataFromRoutes() gotRoutes = %v, want %v", gotRoutes, wantRoutes)
	}
}
