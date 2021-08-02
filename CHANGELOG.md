# Change Log

All notable changes to this project will be documented in this file.

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
