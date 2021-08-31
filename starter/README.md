# Good Starter

Built-in page boilerplate templates

## Usage

These templates should be used with the `good page` command like so

```bash
good page ./my_site newpagename --starter :bootstrap5/layout
```

The leading colon ':' character tells the tool to look at the embedded templates.

## Available Templates

### Basic (default)

A good staring point when you want to build a custom layout.

```
... --starter :basic
```

### Bootstrap v5

The [Bootstrap toolkit](https://getbootstrap.com/docs/5.0) is a great choice for a backend dashboard,
it is mature and well documented. You have a few different setups to choose from:

* `:bootstrap5/layout` A robust layout for dashboards, built using the Bootstrap toolkit.
* `:bootstrap5/examples` A set of functioning demo apps using the BSv5 layout
* `:bootstrap5/datatable` __TODO:__ Bootstrap data table with pagination and sorting
* `:bootstrap5/login` __TODO:__ Login and registration page flow

#### Screenshot of BSv5 Examples

![Alt text](../docs/BSv5_examples.png?raw=true "Bootstrap")


### Info

The introduction that is the default landing page for each new scaffold.

```
... --starter :intro
```

### Minimum

A bare bones page setup, if you truly want to start from scratch.

```
... --starter :minimum
```
