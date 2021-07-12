_Site Package: `github.com/rur/good/baseline/starter_test`_

# `./baseline/starter_test` Web App

This is a self-contained setup for a Golang HTML web app. It contains a series of named pages,
each of which can have a variety of endpoints and functionality.

## Setup

### Generated Code

To ensure that the site code is up to date run the recursive generate command.

```
go generate github.com/rur/good/baseline/starter_test/...
```

### Run the server

Start the server with optional port and dev-mode flags

```
go run github.com/rur/good/baseline/starter_test --port 8080 --dev
```

## Scaffold Overview

This code was originally set up by the [good scaffold](https://www.github.com/rur/good), like so...

```
$ good scaffold ./baseline/starter_test
```

### Code layout

The generated code is intended to be plain and readable. If you are familiar with Golang,
it should be easy to find your way around the code and make the modifications that suit your needs.

> **tip**
>
> You can toggle through the areas of the code important
> for customization by searching for the `// EDITME:` comments.

### Packages & internal import relations

It is important to understand the _imports from_ relationship between the scaffold Go packages
when adding code to your site.

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

Cyclical imports are prohibited by Golang, so for example, the `{site}/page` pkg cannot refer to code
inside any of the named pages, this is by design.

### Site files

These are the key files set up by the `good scaffold` command.

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

### Page files

These are the key files set up by the `good page {name}` command.

| Location                           | Note                                                                    |
| ---------------------------------- | ----------------------------------------------------------------------- |
| page/{name}/routemap.toml          | configuration of route, template and handler mappings                   |
| page/{name}/routes.go              | (generated file) endpoint plumbing generated from the routemap          |
| page/{name}/resources.go           | Request-scoped handler resources, implements `bindResources(myHandler)` |
| page/{name}/handlers.go            | Local request handlers, referenced in routes.go                         |
| page/{name}/templates/\*\*/\*.tmpl | Template files for this page, organized by block name                   |
