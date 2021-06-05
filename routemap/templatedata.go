package routemap

import (
	"fmt"
	"strings"

	"github.com/rur/good/generate"
)

type stackData struct {
	view      RouteView
	extends   string
	blockPath []string
}

func popStack(stack *[]stackData) stackData {
	sLen := len(*stack)
	d := (*stack)[sLen-1]
	*stack = (*stack)[:sLen-1]
	return d
}

// TemplateDataForRoutes will take hierarchical definition of views and flatten to
// data for rendering in the templates
func TemplateDataForRoutes(page PageRoutes, missTpl []Missing, missHlr []Missing) (
	entries []generate.Entry,
	routes []generate.Route,
	templates []generate.HTMLTemplate,
	handlers []generate.Handler,
	err error,
) {
	stack := []stackData{
		{view: page.RouteView},
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
			entries = append(entries, generate.Entry{
				Type:    "Spacer",
				Comment: spacer,
			})
		}

		entry := generate.Entry{
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
			route := generate.Route{
				Reference:    entry.Assignment,
				Path:         view.Path,
				Method:       view.Method,
				PageOnly:     !view.Fragment && !view.Partial,
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
