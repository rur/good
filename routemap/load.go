package routemap

import (
	"fmt"
	"sort"

	toml "github.com/pelletier/go-toml"
)

// TemplateBlock is a named child template slot within in a
// go template. Views is a list of RouteViews directly defined for that
// slot
type TemplateBlock struct {
	Name  string
	Views []RouteView
}

// RouteView is a handler + template pair corresponding to a single
// partial, it may contain a number of slots available for extention
// by sub views
type RouteView struct {
	Ref      string   `toml:"_ref"`
	Default  bool     `toml:"_default"`
	Doc      string   `toml:"_doc"`
	Path     string   `toml:"_path"`
	Template string   `toml:"_template"`
	Handler  string   `toml:"_handler"`
	Method   string   `toml:"_method"`
	Fragment bool     `toml:"_fragment"`
	Page     bool     `toml:"_page"`
	Includes []string `toml:"_includes"`
	Blocks   []TemplateBlock
}

// PageRoutes is the top level view for a site page, it includes
// a URI and golang package namespace
type PageRoutes struct {
	RouteView
	Namespace string `toml:"_namespace"`
	URI       string `toml:"_uri"`
}

// GetFrom will attempt to unmarshal routes from a loaded TOML tree
func GetFrom(tree *toml.Tree) (*PageRoutes, error) {
	var rts PageRoutes

	err := rts.RouteView.UnmarshalFrom(tree)
	if err != nil {
		return nil, err
	}

	if ns, ok := tree.Get("_namespace").(string); ok {
		rts.Namespace = ns
	}
	if uri, ok := tree.Get("_uri").(string); ok {
		rts.URI = uri
	}
	return &rts, nil
}

// UnmarshalFrom will populate the view definition recursively from a
// loaded toml Tree
func (rv *RouteView) UnmarshalFrom(tree *toml.Tree) error {
	err := tree.Unmarshal(rv)
	if err != nil {
		return err
	}

	keys := tree.Keys()
	sort.Strings(keys)
	for _, key := range keys {
		if key[0] == '_' {
			continue
		}
		block := TemplateBlock{
			Name: key,
		}
		val := tree.GetArray(key)
		if subtrees, ok := val.([]*toml.Tree); ok {
			for _, sTree := range subtrees {
				var sView RouteView
				// recursive call
				err := sView.UnmarshalFrom(sTree)
				if err != nil {
					return err
				}
				block.Views = append(block.Views, sView)
			}
		} else {
			return fmt.Errorf("unknown value: %#v", val)
		}
		rv.Blocks = append(rv.Blocks, block)
	}
	return nil
}
