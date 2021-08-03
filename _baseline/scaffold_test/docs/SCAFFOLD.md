# Scaffold Docs

#### Index

- [Code Layout](#code-layout)
- [Site Env](#site-env)
- [Managing Pages](#managing-pages)
- [Site Scaffold Files](#site-scaffold-files)
- [Page Scaffold Files](#page-scaffold-files)

### Intro

The generated code is intended to be plain and readable. If you are familiar with Golang,
it should be easy to find your way around the code and make the modifications that suit your needs.

> **tip**
>
> You can toggle through the areas of the code important
> for customization by searching for the `// EDITME:` comments.

## Code Layout

### Packages & internal dependencies

It is important to understand the dependency relationship within the scaffold Go
packages when adding code to your site.

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

Cyclical imports are prohibited in Golang. For example, the `{site}/page` pkg cannot refer to code
inside any of the named pages, this is by design.

## Site Env

TODO: write this

## Managing Pages

TODO: write this

## HTTP Router

TODO: write this

## Files

### Site Scaffold Files

These are the key files created by the `good scaffold` command.

| Location                  | Note                                                                      |
| ------------------------- | ------------------------------------------------------------------------- |
| main.go                   | Initiailize and start the web server                                      |
| pages.go                  | (generated file) static link to page routes during init                   |
| static/{js styles public} | embedded browser assets                                                   |
| service/env.go            | Site-wide services and config passed to handlers                          |
| service/\*.go             | Place your IO & wrapper code in this package (Auth, Postgres, S3, etc...) |
| page/handlers.go          | Handlers and utilities available to ALL pages                             |
| page/templates/\*\*.tmpl  | Template files available to ALL pages                                     |
| page/{\*name}/            | A named page package (see details below)                                  |
| doc/\*.md                 | Scaffold documentation                                                    |

### Page Scaffold files

These are the key files created by the `good page {name}` command.

| Location                           | Note                                                                    |
| ---------------------------------- | ----------------------------------------------------------------------- |
| page/{name}/routemap.toml          | configuration of route, template and handler mappings                   |
| page/{name}/routes.go              | (generated file) endpoint plumbing generated from the routemap          |
| page/{name}/resources.go           | Request-scoped handler resources, implements `bindResources(myHandler)` |
| page/{name}/handlers.go            | Local request handlers, referenced in routes.go                         |
| page/{name}/templates/\*\*/\*.tmpl | Template files for this page, organized by block name                   |
