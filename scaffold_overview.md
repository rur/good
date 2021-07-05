# Scaffold Overview

The scaffolded code is a self-contained setup for a Golang HTML web app.
It contains a series of named pages, each of which can have a variety of endpoints
and functionality.

## Code layout

The generated code is intended to be plain and readable. If you are familiar with Golang,
it should be easy to find your way around the code and make the modifications that suit your needs.

> **tip**
>
> You can toggle through the areas of the code important
> for customization by searching for the `// EDITME:` comments.

### Site packages & internal dependency relations

It is important to understand the _imports from_ relationship between the scaffold Go packages
when adding code to your site, since cyclical imports are prohibited by Golang.

For example the `{site}/page` Go pkg cannot refer to code inside any of the named pages, this is by design.

```
[main]
    +----------------> [{site}/page/{*name}] --+
    |                              |           |
    +----> [{site}/page] <---------+           |
    |           |                              |
    |           V                              |
    +----> [{site}/service] <------------------+
    |
    +----> [{site}/static](embedded)
```

### Site files

These are the key files set up by the `good scaffold` command.

| Location                  | Note                                                                      |
| ------------------------- | ------------------------------------------------------------------------- |
| main.go                   | Initiailize and start the web server                                      |
| pages.go                  | (generated file) static link to page routes during init                   |
| static/{js styles public} | embedded browser assets                                                   |
| service/\*.go             | Place your IO & wrapper code in this package (Auth, Postgres, S3, etc...) |
| page/handlers.go          | Handlers and utilities available to all pages                             |
| page/templates/\*\*.tmpl  | Template files available to all pages                                     |
| page/{\*name}/            | A named page package (see details below)                                  |

### Page files

These are the key files set up by the `good page {name}` command.

| Location                           | Note                                                                    |
| ---------------------------------- | ----------------------------------------------------------------------- |
| page/{name}/routemap.toml          | configuration of route, template and handler mappings                   |
| page/{name}/routes.go              | (generated file) endpoint plumbing generated from the routemap          |
| page/{name}/resources.go           | Request-scoped handler resources, implements `bindResources(myHandler)` |
| page/{name}/handlers.go            | Local request handlers, referenced in routes.go                         |
| page/{name}/templates/\*\*/\*.tmpl | Template files for this page, organized by block name                   |
