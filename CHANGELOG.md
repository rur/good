# Change Log

All notable changes to this project will be documented in this file.

## [v0.1.5] - 2025-03-16

- Update Semaphore CI README tag
- Allow good pages command to adopt the package name for the pages template
- Update, remove use of deprecated ioutil std library

## [v0.1.4] - 2025-03-16

Update development version of golang to 1.23.


## [v0.1.3] - 2023-02-13

Remove the RWMutex from the resources.go template.

### Improve

- Use the request context to keep resource caching local to the handler goroutine.

## [v0.1.2] - 2022-04-23

Bugfix to correct slice-abuse in generating entries from the route map for the
routes.go template.

## [v0.1.1] - skipped (not released)

Tighten up some code and make page resources optional

### Added

- Added `--no-resources` flag to the page generation command

### Remove

- Trim unused resources and handlers in some of the page starter templates

### Improve

- Make site main file generally easier to read
- Env singleton instance is scopes to the main function
- Initializing Env resources in main func allows the use of defer for teardown

## [v0.1.0] - 2021-12-11

Stamp out an initial release build with a version number.

### Added

- Add the `good version` command to print the current version of the CLI

## [Unreleased] - 2021-10-17

Add interactive confirmation to page command and improve reporting of CLI errors to the
user.

### Breaking Change

- `good page` command will now trigger a wizard in terminal mode, this might break scripts
  - use the `-y` flag to prevent blocking on user input

### Changes

- All commands attempt to explain any expected CLI errors to the user
  - Unexpected errors continue to cause a panic and spew a stack trace
- CLI commands that read interactive input now have a 30 second time limit

### Bugfix

- Recursive `go list ./...`, searching for a go package when _go.mod_ is not found cause CLI to hang
  - Use a 5 second timeout context when calling out to the Go tools.

## [Unreleased] - 2021-10-03

Rename scaffold sub-package to "site" and add a new starter template with
a functioning sign-in and sign-up user flow (mock in-memory database).

### Breaking Change

- Renamed scaffold _service_ package to _site_, since that is more appropriate

### Added

- Added new starter template `:bootstrap5/login`
- Add an `env.HTTPS` flag to the site Env, because it is useful
- The server setup assumes that HTTPS is enabled in non-dev mode
- Added README.md.tmpl to default good starter output
- Add README to the `:intro` starter template
  - same info as the intro index.html file

### Fixed

- Upgrade version of Treetop to **v0.4.1** to get a patch fix for `Vary` header handling

## [Unreleased] - 2021-09-23

Minor improvements, main fix is to allow pages to be added even when a scaffold
site package has a compile error.

### Added

- Add panic recovery to page Env binding helper

### Fixed

- Good Page command will still work when the target scaffold has a compile error.

### Changed

- Minor improvements to scaffold markdown docs/ROUTEMAP.md guide
- Add a function docstring when generating the routes.go file

## [Unreleased] - 2021-08-30

Add Bootstrap v5 built-in starter templates and improve documentation

### Added

- Add Bootstrap v5 starter templates: layout and examples
- root static files are listed during server setup
- Building GUIs, and starter template sections on README
- _page/keyed.go_ file to scaffold page package
-

### Changed

- Rename `--starter-template` to `--starter`
- Scaffold creates a single intro page with instructions on getting started
- Reorg scaffold page templates between _default_ and _named_

## [Unreleased] - 2021-08-02

First major milestone for enhancements from the POC functionality.

### Added

- `--starter-template` for the page scaffold command
- `good starter` command to help create a customized page code template
- README files for the site and pages to serve as documentation
- Integration tests in the form of a 'baseline' file diff
- Routes command now generates missing handlers and templates

### Changed

- Scaffold command will no longer create a number of pages
- Scaffold creates a single intro page with instructions on getting started
- The default starter page has been extracted into folder, embedded with the good binary

### Fix

- Incomplete and buggy processing of routemap for refcount
- Write changes back out to routemap.toml

## [Unreleased] - 2021-06-20

First fully functioning version of the good commands capable of embedding a web-GUI, adding pages and updating/generating routes.

### Added

- `good` CLI with the following commands: `scaffold`, `page`, `pages` and `routes`
