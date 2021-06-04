package generate

import (
	"reflect"
	"testing"
)

func TestSiteFromPagePackage(t *testing.T) {
	tests := []struct {
		name        string
		pkg         GoPackage
		wantSitePkg GoPackage
		wantErr     bool
	}{
		{
			name: "typical",
			pkg: GoPackage{
				Dir:        "/some/path/to/module/site/page/test",
				ImportPath: "github.com/rur/example/site/page/test",
				Module: GoModule{
					Path: "github.com/rur/example",
					Dir:  "/some/path/to/module",
				},
			},
			wantSitePkg: GoPackage{
				Dir:        "/some/path/to/module/site",
				ImportPath: "github.com/rur/example/site",
				Module: GoModule{
					Path: "github.com/rur/example",
					Dir:  "/some/path/to/module",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSitePkg, err := SiteFromPagePackage(tt.pkg)
			if (err != nil) != tt.wantErr {
				t.Errorf("SiteFromPagePackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSitePkg, tt.wantSitePkg) {
				t.Errorf("SiteFromPagePackage() = %v, want %v", gotSitePkg, tt.wantSitePkg)
			}
		})
	}
}
