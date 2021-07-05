# Routemap Overview

The `routemap.toml` file in a site page is the configuration for HTML template composition and
mappings to HTTP routes and handlers.

- The routemap is used to generate the _routes.go_ code for a page (see the `good routes` command)
- It is the _source of truth_ for how templates, handlers and paths are mapped within a page.
- It is capable of describing variations of the page layout
- HTTP endpoints can be bound to a specific layouts

_See [TOML docs](https://toml.io/) for info on the format_

### Example - Basic Composition

The sample below specifies two endpoints: `/mypage` and `/mypage/details`, each with a different
template for the _my-content_ section. In this example, the _base.html.tmpl_ file
will contain a `"my-content"` template block.

##### Example of a routemap file located at `[site]/page/mypage/routemap.toml`

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

### Example - Advanced Composition

The next sample shows more advanced layout features suitable for modern web apps.
The example is a settings page with sections for basic and advanced controls.
The _advanced_ section includes nested sidebar navigation.

In this sample you will find the following paths: `/settings`, `/advanced-settings/user`
and `/advanced-settings/organization`.

Note that the `_default` flag is used to include a given template in the parent
template layout by default.

##### Example of a routemap located at `[site]/page/settings/routemap.toml`

```TOML
_ref = "settings"
_template = "page/templates/base.html.tmpl"
_handler = "hlp.BindEnv(settingsBaseHandler)"

    [sitenav]
    _ref = "pagenav"
    _default = true
    _handler = "treetop.Noop"
    _template = "page/settings/templates/sitenav/pagenav.html.tmpl"

    [content]
    _ref = "basic-settings"
    _path = "/settings"
    _handler = "hlp.BindEnv(bindResources(basicSettingsHandler))"
    _template = "page/settings/templates/content/basic-settings.html.tmpl"

    [content]
    _ref = "advanced-settings"
    _handler = "hlp.BindEnv(bindResources(detailsHandler))"
    _template = "page/settings/templates/content/advanced.html.tmpl"

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

If you wish to understand the plumbing that makes this work, try scanning a generated
[routes.go](https://github.com/rur/good/blob/main/_baseline/routes_test/page/example/routes.go) file and
use the [Treetop Library docs](https://github.com/rur/treetop) as a reference.
