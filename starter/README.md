# Good Starter

Built-in boilerplate page setup templates.

## Usage

These templates can be used with the `good page` command like so

```bash
good page ./my_site newpagename --starter :bootstrap5/layout
```

The leading colon ':' character denotes that this is a built-in starter
template. Without that, the command will attempt to load a path from the
file system.

## Available Templates

### Basic (default)

A good staring point when you want to build a custom page, or use a specific GUI
toolkit.

```
... --starter :basic
```

### Bootstrap v5

The standard Bootstrap toolkit is a great choice for a backend dashboard,
it is mature and well documented. You have a few examples of Bootstrap v5 based
layouts to choose from:

* `:bootstrap5/layout` A robust layout for dashboards, built using the Bootstrap toolkit.
* `:bootstrap5/examples` A set of functioning demo apps using the BSv5 layout
* `:bootstrap5/datatable` __TODO:__ Bootstrap data table with pagination and sorting
* `:bootstrap5/login` __TODO:__ Login and registration page flow

#### Screenshot of Layout

![Alt text](../docs/BSv5_layout.png?raw=true "Bootstrap")


### Info

The introduction page created for each new scaffold.

```
... --starter :intro
```

### Minimum

A bare bones page setup, if you truly want to start from scratch.

```
... --starter :minimum
```
