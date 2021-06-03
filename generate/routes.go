package generate

import (
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
	Name       string
}

// Route is a path mapped to a view definition
type Route struct {
	Method    string
	Path      string
	Includes  []string
	Reference string
	Type      string // "Page" "Fragment" ""
}

// TemplateDataFromRoutes will take hierarchical definition of views and flatten to
// data for rendering in the templates
func TemplateDataFromRoutes(def routemap.PageRoutes) (handlers []Handler, entries []Entry, routes []Route, err error) {
	viewStack := []routemap.RouteView{def.RouteView}
	extendsStack := [][]string{nil}

	// traverse route definitions using a pre-order traversal
	for len(viewStack) > 0 {
		view := popView(&viewStack)
		extends := popStr(&extendsStack)
		hlr := Handler{
			Ref:        view.Ref,
			Extends:    safeLast(extends),
			Method:     view.Method,
			Doc:        view.Doc,
			Identifier: kebabToCamel(view.Ref) + "Handler",
		}
		handlers = append(handlers, hlr)

		for i := len(view.Blocks) - 1; i >= 0; i-- {
			// Note: we can add a sentinel value to the extends stack for spacer entries
			nExt := append(extends, view.Blocks[i].Name)
			for j := len(view.Blocks[i].Views) - 1; j >= 0; j-- {
				viewStack = append(viewStack, view.Blocks[i].Views[j])
				extendsStack = append(extendsStack, nExt)
			}
		}
	}
	return
}

// popStr will return the last element of the slice and shorten it by one
func popStr(stack *[][]string) []string {
	len := len(*stack)
	str := (*stack)[len-1]
	*stack = (*stack)[:len-1]
	return str
}

// popView will return the last element of the slice and shorten it by one
func popView(stack *[]routemap.RouteView) routemap.RouteView {
	len := len(*stack)
	view := (*stack)[len-1]
	*stack = (*stack)[:len-1]
	return view
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
