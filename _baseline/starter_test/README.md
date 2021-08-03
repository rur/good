# `github.com/rur/good/baseline/starter_test`

#### Web App

This is a self-contained setup for a Golang HTML web app. Refer to the
`./page/*/routemap.toml` files to explore endpoints.

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

## Add a page

Use the good page command like so to add pages to the scaffold site.

```bash
$ good page ./baseline/starter_test yourpagename
```

#### Page names

Page names must contain lowercase `a` to `z` characters only, and be no shorter than
three characters long.

## Scaffold Overview

This code was originally set up by the [good scaffold](https://www.github.com/rur/good), like so...

```
$ good scaffold ./baseline/starter_test
```

Refer to the [Scaffold Docs](./docs/SCAFFOLD.md) for details
