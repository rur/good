package routemap

import (
	"fmt"
	"path"
	"regexp"
	"sort"

	toml "github.com/pelletier/go-toml"
)

var (
	// view '_ref' must be strictly lower, kebab-case string
	refRegex = regexp.MustCompile("^[a-z][a-z0-9]*(-[a-z][a-z0-9]*)*$")

	// TODO: consider making the values regexp for validation purposes
	knownKeys = map[string]bool{
		"_ref":      true,
		"_uri":      true,
		"_default":  true,
		"_doc":      true,
		"_path":     true,
		"_template": true,
		"_handler":  true,
		"_method":   true,
		"_fragment": true,
		"_partial":  true,
		"_merge":    true,
		"_includes": true,
	}
)

// TemplateBlock is a named child template slot within in a
// go template. Views is a list of RouteViews directly defined for that
// slot
type TemplateBlock struct {
	Name  string
	Views []RouteView

	// TODO: byRef map[string]*RouteView
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
	Partial  bool     `toml:"_partial"`
	Includes []string `toml:"_includes"`
	Blocks   []TemplateBlock
}

// PageRoutes is the top level view for a site page, it includes
// a URI and golang package namespace
type PageRoutes struct {
	RouteView
	URI string `toml:"_uri"`
}

// Missing is a missing value in the TOML file tree, it pairs the route view data
// with the position reference in the source TOML file
type Missing struct {
	Ref      string
	Position toml.Position
}

// ProcessRoutemap will unmarshal and validate a TOML routemap.
// This will keep track of missing template and handler values.
func ProcessRoutemap(tree *toml.Tree, templatePath string) (routes PageRoutes, templates []Missing, handlers []Missing, err error) {
	parser := routeParser{
		usedRefs:     make(map[string]toml.Position),
		templateBase: templatePath,
	}
	routes.RouteView, err = parser.unmarshalView(tree)
	if err != nil {
		return
	}
	if uri, ok := tree.Get("_uri").(string); ok {
		routes.URI = uri
	}
	templates = parser.missingTemplates
	handlers = parser.missingHandlers
	return
}

// parser state during unmarshal/validate process
type routeParser struct {
	usedRefs         map[string]toml.Position
	missingTemplates []Missing
	missingHandlers  []Missing
	blockPath        []string
	templateBase     string
}

// unmarshalView will parse fields, validate and descend to unmarshal sub views routes
func (parser *routeParser) unmarshalView(tree *toml.Tree) (view RouteView, err error) {
	// Popupate struct with known '_*' properties
	err = tree.Unmarshal(&view)
	if err != nil {
		pos := tree.Position()
		err = fmt.Errorf(":%d:%d: unmarshal error, %s", pos.Line, pos.Col, err)
		return
	}
	// validate reference value and check for uniqueness
	if !refRegex.MatchString(view.Ref) {
		pos := tree.Position()
		err = fmt.Errorf(
			":%d:%d: Unknown or invalid _ref '%s', references be all lowercase joined by a dash '-'",
			pos.Line, pos.Col, view.Ref,
		)
		return
	}
	if usePos, ok := parser.usedRefs[view.Ref]; ok {
		pos := tree.Position()
		err = fmt.Errorf(
			":%d:%d: duplicate _ref '%s', already used in routemap file at line %d, column %d",
			pos.Line, pos.Col, view.Ref, usePos.Line, usePos.Col,
		)
		return
	} else {
		// stash the location in the input file where this reference name appears
		parser.usedRefs[view.Ref] = tree.GetPosition("_ref")
	}
	// fill template field if missing
	if view.Template == "" {
		view.Template = path.Join(append(append([]string{parser.templateBase}, parser.blockPath...), view.Ref+".html.tmpl")...)
		parser.missingTemplates = append(parser.missingTemplates, Missing{
			Position: tree.Position(),
			Ref:      view.Ref,
		})
	}
	// fill handler field if missing
	if view.Handler == "" {
		view.Handler = fmt.Sprintf("%sHandler", kebabToCamel(view.Ref))
		parser.missingHandlers = append(parser.missingTemplates, Missing{
			Position: tree.Position(),
			Ref:      view.Ref,
		})
	}

	// Now descend into sub view blocks by scanning for non-underscore keys
	keys := tree.Keys()
	sort.Strings(keys) // make output stable (we don't care about the original order)
	currentBlockPath := parser.blockPath
	for _, key := range keys {
		if knownKeys[key] {
			continue
		}
		pos := tree.GetPositionPath([]string{key})
		// this is a block name, validate format
		if !refRegex.MatchString(key) {
			err = fmt.Errorf(
				":%d:%d: Unknown or invalid key '%s', block names must be all lowercase joined by a dash '-'",
				pos.Line, pos.Col, key,
			)
			return
		}
		parser.blockPath = append(currentBlockPath, key)
		block := TemplateBlock{
			Name: key,
		}
		val := tree.GetArray(key)
		if subtrees, ok := val.([]*toml.Tree); ok {
			for _, sTree := range subtrees {
				var sView RouteView
				// recursive call
				sView, err = parser.unmarshalView(sTree)
				if err != nil {
					return
				}
				block.Views = append(block.Views, sView)
			}
		} else {
			err = fmt.Errorf(
				":%d:%d: invalid value for key '%s', expecting an array of tables, got %#v",
				pos.Line, pos.Col, key, val,
			)
			return
		}
		view.Blocks = append(view.Blocks, block)
	}
	parser.blockPath = currentBlockPath
	return
}
