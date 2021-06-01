package generate

import (
	toml "github.com/pelletier/go-toml"
)

type RouteBlock struct {
	Name string
}

type RouteView struct {
	Name     string `toml:"_name"`
	Default  bool   `toml:"_default"`
	Doc      string `toml:"_doc"`
	Path     string `toml:"_path"`
	Template string `toml:"_template"`
	Handler  string `toml:"_handler"`
	Blocks   []Block
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

	return &rts, nil
}
