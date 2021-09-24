package example

import (
	"github.com/rur/good/baseline/routes_test/page"
	"github.com/rur/treetop"
)

// Routes is the plumbing code for page endpoints, templates and handlers
func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	// Code created by go generate. You should edit the routemap.toml file; DO NOT EDIT.

	example := treetop.NewView(
		"page/example/templates/example.html.tmpl",
		hlp.BindEnv(bindResources(exampleHandler)),
	)

	// [[content]]
	examplePlaceholder := example.NewDefaultSubView(
		"content",
		"page/example/templates/content/example-placeholder.html.tmpl",
		treetop.Constant(struct {
			Form interface{}
		}{
			Form: struct{ FormError string }{
				FormError: "This is a test",
			},
		}),
	)

	// [[content.form]]
	placeholderForm := examplePlaceholder.NewDefaultSubView(
		"form",
		"page/example/templates/content/form/placeholder-form.html.tmpl",
		hlp.BindEnv(bindResources(placeholderFormHandler)),
	)

	// [[content.form.form-error]]
	basicFormError := placeholderForm.NewDefaultSubView(
		"form-error",
		"page/example/templates/content/form/form-error/basic-form-error.html.tmpl",
		hlp.BindEnv(bindResources(basicFormErrorHandler)),
	)

	// [[content]]
	alternativeContent := example.NewSubView(
		"content",
		"page/example/templates/content/alternative-content.html.tmpl",
		hlp.BindEnv(bindResources(alternativeContentHandler)),
	)
	settingsLayout := example.NewSubView(
		"content",
		"page/example/templates/content/settings-layout.html.tmpl",
		hlp.BindEnv(bindResources(settingsLayoutHandler)),
	)

	// [[content.settings]]
	generalSettings := settingsLayout.NewSubView(
		"settings",
		"page/example/templates/content/settings/general-settings.html.tmpl",
		hlp.BindEnv(bindResources(generalSettingsHandler)),
	)
	advancedSettings := settingsLayout.NewSubView(
		"settings",
		"page/example/templates/content/settings/advanced-settings.html.tmpl",
		hlp.BindEnv(bindResources(advancedSettingsHandler)),
	)

	// [[content.settings.settings-form]]
	updateAdvancedSettings := advancedSettings.NewSubView(
		"settings-form",
		"page/example/templates/content/settings/settings-form/update-advanced-settings.html.tmpl",
		hlp.BindEnv(bindResources(updateAdvancedSettingsHandler)),
	)
	updateAdvancedSettings.HasSubView("form-error")

	// [[content.tabs]]
	settingsLayout.NewDefaultSubView(
		"tabs",
		"page/example/templates/content/tabs/settings-tabs.html.tmpl",
		treetop.Constant("Hello World"),
	)

	// [[nav]]
	example.NewDefaultSubView(
		"nav",
		"page/example/templates/nav/main-nav.html.tmpl",
		hlp.BindEnv(bindResources(mainNavHandler)),
	)

	// [[scripts]]
	example.NewDefaultSubView(
		"scripts",
		"page/example/templates/scripts/site-scripts.html.tmpl",
		treetop.Noop,
	)

	hlp.Handle("/example",
		exec.NewViewHandler(examplePlaceholder))
	hlp.HandlePOST("/example/form",
		exec.NewViewHandler(placeholderForm).FragmentOnly())
	hlp.Handle("/example/alt",
		exec.NewViewHandler(alternativeContent).PageOnly())
	hlp.Handle("/example/settings",
		exec.NewViewHandler(generalSettings))
	hlp.Handle("/example/advanced-settings",
		exec.NewViewHandler(advancedSettings))
	hlp.HandlePOST("/example/advanced-settings/submit",
		exec.NewViewHandler(
			updateAdvancedSettings,
			basicFormError,
		).FragmentOnly())

}
