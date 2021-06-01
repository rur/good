package service

import (
	"sort"
	"strings"
)

// Sitemap is a utility to help handler to generate URLs based on
// routing configuration
type Sitemap map[string]Routemap

// Routemap is the hierarchical parsed from all page route map files
type Routemap struct {
	Block  string
	Path   string
	URI    string
	Routes map[string]Routemap
}

// Pages list the pages in the site in lex' sorted order
func (sm Sitemap) Pages() (pages []string) {
	for page := range sm {
		pages = append(pages, page)
	}
	sort.Strings(pages)
	return
}

// GetPath will return a routing path for a given route view name path
// starting with the page, delimited by " > "
//
// eg. The following will return the path value for the specified views
//
// 		sm.GetPath("example > placeholder-content > my-form")
//
func (sm Sitemap) GetPath(path string) (string, bool) {
	if path == "" {
		return "", false
	}
	parts := strings.Split(path, " > ")
	cursor, ok := sm[parts[0]]
	for i := 1; ok && i < len(parts); i++ {
		cursor, ok = cursor.Routes[parts[i]]
	}
	if !ok {
		return "", false
	} else if cursor.Path != "" {
		return cursor.Path, true
	}
	return "", false
}
