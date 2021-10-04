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

It is important to understand the dependency relationships within the Go
packages of the scaffold.

```
# Go packages:

 ┌─────────────────┐
 │ {scaffold}/main │
 └────────┬────────┘
          │
          │      ┌────────────────────────┐
          ├─────►│ {scaffold}/page/{name} ├─┬───┐
          │      └────────────────────────┘ │   │
          │                                 │   │
          │      ┌─────────────────┐        │   │
          ├─────►│ {scaffold}/page │◄───────┘   │
          │      └───────┬─────────┘            │
          │              │                      │
          │              │                      │
          │      ┌───────▼─────────┐            │
          ├─────►│ {scaffold}/site │◄───────────┘
          │      └─────────────────┘
          │
          │      ┌───────────────────┐
          └─────►│ {scaffold}/static │ (embedded)
                 └───────────────────┘
```

> **note**
>
> Cyclical imports are prohibited in Golang. Therefore, the `{scaffold}/page` pkg cannot refer to code
> inside any of the named pages, this is by design.

## Site Env

This singleton _struct_ is initialized at startup and passed to handers using the `hlp.BindEnv` function.

Put your site-wide stuff here like static config, connection pools, etc...

## Managing Pages

TODO: write details about page/pages commands and add a sub header about custom starter pages.

## HTTP Router

The standard `net/http` server Mux is fine for basic needs. If you have a preferred routing library with
more advanced features (CSRF etc..), modify the `[site]/page/helper.go` file and correct any compiler errors;
that should do it!

## Files

### Site Scaffold Files

These are the key files created by the `good scaffold` command.

| Location                  | Note                                                                      |
| ------------------------- | ------------------------------------------------------------------------- |
| main.go                   | Initialize and start the web server                                       |
| pages.go                  | (generated file) static link to page routes during init                   |
| static/{js styles public} | embedded browser assets                                                   |
| site/env.go               | Site-wide services and config passed to handlers                          |
| site/\*.go                | Place your standard types, IO & wrapper code here (eg. Auth, Postgres, S3, etc...) |
| page/handlers.go          | Handler code available to ALL pages                                       |
| page/helper.go            | Helper type passed to page routes.go during initialization                |
| page/keyed.go             | Key map of hard coded raw template strings (comes in handy)               |
| page/\*.go                | Misc shared page utilities                                                |
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
| page/{name}/handlers_\*.go            | Local request handlers, referenced in routes.go                         |
| page/{name}/templates/\*\*/\*.tmpl | Template files for this page, usually organized by block name                   |
