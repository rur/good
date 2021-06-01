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

func GetRoutes(routemap string) (*PageRoutes, error) {
	tree, err := toml.Load(routemap)
	if err != nil {
		return nil, err
	}
	var rts PageRoutes

	err = tree.Unmarshal(&rts)
	if err != nil {
		return nil, err
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
			for _, stree := range subtrees {
				var sView RouteView
				err := stree.Unmarshal(&sView)
				if err != nil {
					return nil, err
				}
				block.Views = append(block.Views, sView)
			}
		} else {
			return nil, fmt.Errorf("unknown value: %#v", val)
		}
		rts.Blocks = append(rts.Blocks, block)
	}

	return &rts, nil
}
