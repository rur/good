package generate

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestGetRoutes(t *testing.T) {
	toml, err := ioutil.ReadFile("./testdata/routemap.toml")
	if err != nil {
		t.Fatal("failed ot load test data", err)
	}
	tests := []struct {
		name    string
		toml    string
		want    *PageRoutes
		wantErr bool
	}{
		{
			name: "basic",
			toml: string(toml),
			want: &PageRoutes{
				Namespace: "github.com/rur/good/admin/site",
				URI:       "/example",
				RouteView: RouteView{
					Name:     "example",
					Doc:      "Base HTML template for example page",
					Template: "page/templates/base.html.tmpl",
					Handler:  "hlp.BindEnv(page.BaseHandler)",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRoutes(tt.toml)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}
