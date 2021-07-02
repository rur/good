
## Routemap Overview

A `routemap.toml`* file is included inside each site page package. This configuration file combines classic HTML template composition with HTTP routes.
 
- It is capable of binding different request endpoints to variations of the page layout
- Acts as the _source of truth_ for a page, mapping template files to handlers and paths
- The `good routes ...` command uses this routemap to generate the _routes.go_ code for a page

_* See [TOML docs](https://toml.io/) for an intro to the file format_

#### Basic Composition
The sample below specifies two endpoints: `/mypage` and `/mypage/details`, each with a different template for the _my-content_ section. There must be a `"my-content"` template block positioned inside the parent template (_base.html.tmpl_) for this to work.

Example of a routemap located at `[site]/page/mypage/routemap.toml`
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

#### Advanced Composition

The next sample shows more advanced layout features suitable for modern web apps. 
This routemap is for a settings page with separate sections for basic and advanced controls. 
The advanced section has additional sidebar navigation. You will find nested _sidebar_ and _tabcontent_
template blocks inside the `_ref = "advanded-settings"` view in order to support this.

A good tip to grok a page config to to find path endpoints and understand the layout variations from there.
In this sample you will find the following paths: `/settings`, `/advanced-settings/user`
and `/advanced-settings/organization`. 

Note that the `_default` flag is used to include a given template in the parent layout by default.

Sample file `[site]/page/settings/routemap.toml`
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
   
If you wish to understand the plumbing that makes this work, read a [routes.go](https://github.com/rur/good/blob/main/_baseline/routes_test/page/example/routes.go) file and
see the [Treetop Library](https://github.com/rur/treetop) for details.
   
