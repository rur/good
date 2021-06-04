package generate

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/rur/good/routemap"
)

// Handler is data for a handler function which should be created
type Handler struct {
	Ref        string
	Extends    string
	Method     string
	Doc        string
	Identifier string
	Blocks     []HandleBlock
}

// HandleBlock is the details of sub-views which should
// be delegated to in the handler
type HandleBlock struct {
	FieldName string
	Name      string
}

// Entry for the routes.go file
type Entry struct {
	Type       string // "SubView" "DefaultSubView" "Spacer"
	Assignment string
	Extends    string
	Block      string
	Template   string
	Handler    string
	Comment    string
}

// Route is a path mapped to a view definition
type Route struct {
	Method       string
	Path         string
	Includes     []string
	Reference    string
	PageOnly     bool
	FragmentOnly bool
}

type stackData struct {
	view      routemap.RouteView
	extends   string
	blockPath []string
}

func popStack(stack *[]stackData) stackData {
	sLen := len(*stack)
	d := (*stack)[sLen-1]
	*stack = (*stack)[:sLen-1]
	return d
}

// TemplateDataFromRoutes will take hierarchical definition of views and flatten to
// data for rendering in the templates
func TemplateDataFromRoutes(def routemap.PageRoutes) (entries []Entry, routes []Route, err error) {
	stack := []stackData{
		{view: def.RouteView},
	}
	var spacer string

	// emitting entries using a pre-order traversal will ensure that all view variable are declared
	// before they are used to create sub views
	for len(stack) > 0 {
		sData := popStack(&stack)
		view := sData.view

		if sp := fmtSpacer(sData.blockPath); sp != "" && sp != spacer {
			// add a separator to make the routemap code easier to follow
			spacer = sp
			entries = append(entries, Entry{
				Type:    "Spacer",
				Comment: spacer,
			})
		}

		entry := Entry{
			Block:    safeLast(sData.blockPath),
			Extends:  sData.extends,
			Template: view.Template,
			Handler:  view.Handler,
		}

		if len(sData.blockPath) == 0 {
			entry.Type = "PageView"
		} else if view.Default {
			entry.Type = "DefaultSubView"
		} else {
			entry.Type = "SubView"
		}

		if view.Path != "" {
			entry.Assignment = kebabToCamel(view.Ref)
			route := Route{
				Reference:    entry.Assignment,
				Path:         view.Path,
				Method:       view.Method,
				PageOnly:     view.Page,
				FragmentOnly: view.Fragment,
			}
			for _, incl := range view.Includes {
				route.Includes = append(route.Includes, kebabToCamel(incl))
			}
			routes = append(routes, route)
		}

		if len(view.Blocks) > 0 {
			entry.Assignment = kebabToCamel(view.Ref)

			for i := len(view.Blocks) - 1; i >= 0; i-- {
				blockName := view.Blocks[i].Name
				nBlock := append(sData.blockPath, blockName)

				for j := len(view.Blocks[i].Views) - 1; j >= 0; j-- {
					stack = append(stack, stackData{
						view:      view.Blocks[i].Views[j],
						extends:   entry.Assignment,
						blockPath: nBlock,
					})
				}
			}
		}
		entries = append(entries, entry)
	}
	return
}

func fmtSpacer(blocks []string) string {
	if len(blocks) == 0 {
		return ""
	}
	return fmt.Sprintf("[[%s]]", strings.Join(blocks, "."))
}

// kebabToCamel coverts an kebab-case string to camelCase
func kebabToCamel(str string) string {
	parts := strings.Split(str, "-")
	out := []byte(parts[0])
	for i := 1; i < len(parts); i++ {
		out = append(out, strings.Title(parts[i])...)
	}
	return string(out)
}

// safeLast returns the last string in a slice or empty string if the slice is empty
func safeLast(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	return arr[len(arr)-1]
}

// RoutesScaffold will generate all files for the good routes command
func RoutesScaffold(pageName string, config routemap.PageRoutes, scaffold fs.FS) (files []File, err error) {
	entries, routes, err := TemplateDataFromRoutes(config)
	if err != nil {
		return
	}
	data := struct {
		Name      string
		Namespace string
		Entries   []Entry
		Routes    []Route
	}{
		Name:      pageName,
		Namespace: config.Namespace,
		Entries:   entries,
		Routes:    routes,
	}
	// page/name/routes.go
	files = append(files, File{
		Dir:      filepath.Join("page", pageName),
		Name:     "routes.go",
		Contents: mustExecute("scaffold/page/name/routes.go.tmpl", data, scaffold),
	})
	return
}
