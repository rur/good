# Page Routemap Docs

#### Index:

- [Overview](#overview)
- [Examples](#examples)
- [Handlers](#handlers)
- [Templates](#templates)
- [Route Map TOML Reference](#route-map-toml-reference)

## Overview

The `routemap.toml` file is the configuration for template composition and
mappings to HTTP routes and handlers in this page.

- It is the _source of truth_ for how templates, handlers and paths are mapped within a page.
- It is used to generate the _routes.go_ code with the `good routes` command
- It is capable of describing many variations of the page layout
- HTTP endpoints can be bound to a specific layouts

#### Modifying the Routemap

Add or change endpoints by editing the TOML file and running the go generate command
to update the page routes code.

```
go generate github.com/rur/good/baseline/starter_test/...
```

##### Auto-gen handlers and template files

It can be tedious to create and configure template files and handlers when adding new views
so, as a convenience, if the `_template` or `_handler` fields are missing from the routemap when the
_routes.go_ file is generated, the missing items will be generated for you.

## Examples

### Example of basic composition

The sample below specifies two endpoints: `/mypage` and `/mypage/details`, each with a different
template for the _my-content_ section. In this example, the _base.html.tmpl_ file
will contain a `"my-content"` template block.

```TOML
_ref = "mypage"
_template = "page/templates/base.html.tmpl"
_handler = "hlp.BindEnv(myPageBaseHandler)"

    [my-content]
    _ref = "welcome"
    _path = "/mypage"
    _handler = "hlp.BindEnv(bindResources(welcomeHandler))"
    _template = "page/mypage/templates/content/welcome.html.tmpl"

    [my-content]
    _ref = "details"
    _path = "/mypage/details"
    _handler = "hlp.BindEnv(bindResources(detailsHandler))"
    _template = "page/mypage/templates/content/details.html.tmpl"

```

Preview of the `base.html.tmpl` template with a _my-content_ block positioned in the layout.

```html
<!DOCTYPE html>
<html>
  ... {{ template "my-content" .MyContent }} ...
</html>
```

### Example of advanced composition

The next sample shows more advanced layout features suitable for modern web apps.
The example is a settings page with sections for basic and advanced controls.
The _advanced_ section includes nested sidebar navigation.

In this sample you will find the following paths: `/settings`, `/advanced-settings/user`
and `/advanced-settings/organization`.

Note that the `_default` flag is used to include a given template in the parent
template layout by default.

```TOML
_ref = "settings"
_template = "page/templates/base.html.tmpl"
_handler = "hlp.BindEnv(settingsBaseHandler)"
_entry = "/settings"

    [sitenav]
    _ref = "pagenav"
    _default = true
    _handler = "treetop.Noop"
    _template = "page/settings/templates/sitenav/pagenav.html.tmpl"
    _doc = "Default settings navigation"

    [content]
    _ref = "basic-settings"
    _path = "/settings"
    _handler = "hlp.BindEnv(bindResources(basicSettingsHandler))"
    _template = "page/settings/templates/content/basic-settings.html.tmpl"
    _doc = "Standard user settings, nothing fancy"

    [content]
    _ref = "advanced-section"
    _handler = "hlp.BindEnv(bindResources(advancedSectionHandler))"
    _template = "page/settings/templates/content/advanced-section.html.tmpl"
    _doc = "Layout for advanced settings requiring sub navigation"

        [content.sidebar]
        _ref = "adv-tabs"
        _default = true
        _handler = "hlp.BindEnv(advancedTabsHandler)"
        _template = "page/settings/templates/content/sidebar/adv-tabs.html.tmpl"

        [content.tabcontent]
        _ref = "advanced-user"
        _path = "/advanced-settings/user"
        _template = "page/settings/templates/content/tabcontent/advanced-user.html.tmpl"
        _handler = "hlp.BindEnv(bindResources(advancedUserHandler))"

        [content.tabcontent]
        _ref = "advanced-org"
        _path = "/advanced-settings/organization"
        _template = "page/settings/templates/content/tabcontent/advanced-org.html.tmpl"
        _handler = "hlp.BindEnv(bindResources(advancedOrgHandler))"

```

Indentation is optional, if you wish to understand the plumbing, try scanning the generated routes.go file.
It might be helpful to use the [Treetop Library docs](https://pkg.go.dev/github.com/rur/treetop) as a reference.

## Handlers

The job of a handler is to load data for a template, as outlined in the route map.

There are some predefined handlers
for trivial tasks like constants `treetop.Constant(123)` or no-op `treetop.Noop`, but typically you will pass a custom
function that loads template data for the layout.

The minimal, top level signature that is required is this

```
func myHandler(resp treetop.Response, req *http.Request) interface{} {
  return "data for the template"
}
```

### Child Handlers

The _routemap_ defines a layout hierarchy. Since child templates are embedded in a parent, child handlers
are also controlled by a parent handler. The child handler will only be called when the parent
dispatches to it. This happens at arms-length using the _HandleSubView_ method. See an example below.

```
func parentHandler(resp treetop.Response, req *http.Request) interface{} {
  return {
    SomeChildBlock interface{}
  }{
    // synchronous dispatch to child handler, returns child template data
    SomeChildBlock: resp.HandleSubView("some-child-block", req),
  }
}
```

The top level handler is the entry point for the request handling lifecycle.

### Binding Resources

Site config and resources are bound to handlers using type-safe function closures.

#### Site Env

If your handler needs to access site-wite services or configuration, you can bind an extra parameter
in your route map TOML like so,

`_handler = "hlp.BindEnv(myHandlerWithEnv)"`

Your handler must have the following function signature.

```
func myHandlerWithEnv(env service.Env, resp treetop.Response, req *http.Request) interface{} {
  return "data for the template"
}
```

#### Request Resources

Within the layout page package the `resources.go` file load a set of per-request resources, like user details
or DB data based on the request. You can bind an additional resources parameter to you handlers in the
route map like so

`_handler = "hlp.BindEnv(bindResources(myHandlerWithResources))"`

Your handler must have the following function signature.

```
func myHandlerWithResources(rsc resources, env service.Env, resp treetop.Response, req *http.Request) interface{} {
  return "data for the template"
}
```

## Templates

This set up uses standard HTML Golang template files. The files are organized within
the page _templates_ directory with a `.html.tmpl` extension (by convention).

#### Template Embedding

The file byte embedding features of the Go compiler<sup>[[go1.16](https://golang.org/doc/go1.16#library-embed)]</sup> are used to include assets in you binary at compile time.
The scaffold site is configured to embed the following template directories:

- shared page templates `[site]/page/templates` and
- template files for each named page `[site]/page/{name}/templates`

Note templates are read from disk in development mode using the `--dev` flag.

## Route Map TOML Reference

A Route Map is a tree of views, each node on the tree can have the following
fields:

- **\_ref** `string` _(required)_ a reference unique within the TOML file (alpha-only with dash sep)
- **\_doc** `string` optional docstring for this node
- **\_path** `string` Attach a HTTP endpoint to this view, a view hierarchy will be used to construct a `http.Handler` instance.
- **\_template** `string` the template path for this view
- **\_handler** `string` go code referencing template handler function
- **\_method** `string` HTTP method to restrict routing too semantics depend on your chosen router
- **\_default** `bool` flag to indicate that this view should be included in the parent layout by default
- **\_fragment** `bool` This view can be loaded independently using an XHR template request
- **\_partial** `bool` This view can be loaded independently, but can also be loaded as a regular full page using the same path
- **\_merge** `string` (documentation) The treetop-merge method ascribed to the top level element
- **\_includes** `list<string>` Other views in this file to include in this layout, named by _\_ref_
- **\_entrypoint** `string` (root only) document the entry point path for this page

### Routemap Hierarchy

Each view node can declare a set of named blocks. These serve as slots that a child
template can fit into.

Each block has zero or more alternative views declared. Consider this example
translated to JSON; the top level view has two alternative sub-views for the
`"content"` block.

View `_ref: "first-content"` has a further sub block of `"sub-content"`.

#### TOML Version

```TOML
_ref = "mypage"
_template = "some/template.html"
_handler = "someHandlerFunc"

  [[content]]
  _ref = "first-content"
  _template = "content/first.html"
  _handler = "firstContentHandler"

    [[content.sub-content]]
    _ref = "nested-content"
    _template = "content/nested.html"
    _handler = "nestedContentHandler"

  [[content]]
  _ref = "second-content"
  _template = "content/seconds.html"
  _handler = "secondContentHandler"

```

#### JSON Version

```JSON
{
  "_ref": "mypage",
  "_template": "some/template.html",
  "_handler": "someHandlerFunc",

  "content": [{
    "_ref": "first-content",
    "_template": "content/first.html",
    "_handler": "firstContentHandler",

    "sub-content": [{
      "_ref": "nested-content",
      "_template": "content/nested.html",
      "_handler": "nestedHandler"
    }]
  }, {
    "_ref": "second-content",
    "_template": "content/seconds.html",
    "_handler": "secondContentHandler"
  }]
}
```

The TOML version is easier on the eyes!
