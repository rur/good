# Scaffold Overview

The scaffolded code is a plain and self-contained setup for a Golang HTML web app with several named pages, each of which can have multiple endpoints.


## Code layout

The generated code is intended to be plain and readable so that it is easy to find your way around the code and make the modifications that suit your needs.

> __tip__
> 
> You can toggle through the areas of the code important
> for customization by searching for the `// EDITME:` comments.

### Site packages & internal dependency relations

It is important to understand the relationship between the site Go packages
when adding your code to the scaffold.

```  
[main]
    +----------------> [./page/{*name}] --+
    |                         |           |
    +----> [./page] <---------+           |
    |           |                         |
    |           V                         |
    +----> [./service] <------------------+
    |
    +----> [./static/{js styles public}](embedded)
```
Note that cyclical dependencies are strictly prohibited by Golang. For example,
the `[site]/page` pkg cannot refer to code inside any of the named pages, by design.

### Site files

This is the file structure set up by the `good scaffold` command.

| Location                  | Note 
|---------------------------|--------------
| main.go                   | Initiailize and start the web server
| pages.go                  | (generated file) static link to page routes during init
| static/{js styles public} | embedded browser assets
| service/*.go              | Place your IO & wrapper code here (Auth, Postgres, S3, etc...) 
| page/*.go                 | Handlers and utilities available to all pages
| page/templates/**.tmpl    | Template files available to all pages
| page/{*name}/             | A named page package, eg. name=mydashboard (see details below)


### Page files

This is the file structure set up by the `good page {name}` command. 

| Location                  | Note 
|---------------------------|--------------
| page/{name}/routemap.toml | configuration of route, template and handler mappings
| page/{name}/routes.go     | (generated file) endpoint plumbing generated from the routemap
| page/{name}/resources.go  | Request-scoped handler resources, implements `bindResources(myHandler)`
| page/{name}/*.go          | Local request handlers, referenced in routes.go 
| page/{name}/templates/\*\*/*.tmpl | Template files for this page, organized by block name


