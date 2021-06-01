package generate

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
)

type TemplateBlock struct {
	Name  string
	Views []RouteView
}

type RouteView struct {
	Name     string `toml:"_name"`
	Default  bool   `toml:"_default"`
	Doc      string `toml:"_doc"`
	Path     string `toml:"_path"`
	Template string `toml:"_template"`
	Handler  string `toml:"_handler"`
	Method   string `toml:"_method"`
	Fragment bool   `toml:"_fragment"`
	Page     bool   `toml:"_page"`
	Blocks   []TemplateBlock
}

type PageRoutes struct {
	RouteView
	Namespace string `toml:"_namespace"`
	URI       string `toml:"_uri"`
}

func LoadRouteRoutemap(routemap string) (*PageRoutes, error) {
	tree, err := toml.Load(routemap)
	if err != nil {
		return nil, err
	}
	var rts PageRoutes

	err = rts.RouteView.UnmarshalFrom(tree)
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
