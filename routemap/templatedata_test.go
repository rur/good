package routemap

import (
	"reflect"
	"strings"
	"testing"

	"github.com/rur/good/generate"
)

func TestTemplateDataFromRoutes(t *testing.T) {
	def := PageRoutes{
		EntryPoint: "/my-page",
		RouteView: RouteView{
			Ref:      "mypage",
			Template: "page/mypage/templates/mypage.html.tmpl",
			Handler:  "mypageHandler",
			Doc:      "Test page docs",
			Blocks: []TemplateBlock{
				{
					Name: "content",
					Views: []RouteView{
						{
							Ref:      "my-content",
							Template: "page/mypage/templates/content/my-content.html.tmpl",
							Handler:  "myContentHandler",
							Doc:      "The default content",
							Default:  true,
							Path:     "/my-page",
							Partial:  true,
							Method:   "GET",
							Includes: []string{"my-nav"},
							Blocks: []TemplateBlock{
								{
									Name: "form",
									Views: []RouteView{
										{
											Ref:      "my-form",
											Template: "page/mypage/templates/content/form/my-form.html.tmpl",
											Handler:  "myFormHandler",
											Doc:      "A content form",
											Path:     "/my-page/form",
											Fragment: true,
										},
									},
								},
							},
						},
						{
							Ref:      "other-content",
							Template: "page/mypage/templates/content/other-content.html.tmpl",
							Handler:  "otherContentHandler",
							Doc:      "The other content",
							Path:     "/my-page/other",
							Method:   "POST",
						},
					},
				},
				{
					Name: "nav",
					Views: []RouteView{
						{
							Ref:      "my-nav",
							Template: "page/mypage/templates/nav/my-nav.html.tmpl",
							Handler:  "myNavHandler",
							Doc:      "The default nav",
							Default:  true,
						},
					},
				},
			},
		},
	}

	// TODO: test missing templates and handlers
	gotEntries, gotRoutes, _, _, err := TemplateDataForRoutes(def, nil, nil)

	wantEntries := []generate.Entry{
		{
			Assignment: "mypage",
			Block:      "",
			Extends:    "",
			Template:   "page/mypage/templates/mypage.html.tmpl",
			Handler:    "mypageHandler",
			Type:       "PageView",
		},
		{
			Type:    "Spacer",
			Comment: "[[content]]",
		},
		{
			Assignment: "myContent",
			Extends:    "mypage",
			Block:      "content",
			Template:   "page/mypage/templates/content/my-content.html.tmpl",
			Handler:    "myContentHandler",
			Type:       "DefaultSubView",
		},
		{
			Type:    "Spacer",
			Comment: "[[content.form]]",
		},
		{
			Assignment: "myForm",
			Extends:    "myContent",
			Block:      "form",
			Template:   "page/mypage/templates/content/form/my-form.html.tmpl",
			Handler:    "myFormHandler",
			Type:       "SubView",
		},
		{
			Type:    "Spacer",
			Comment: "[[content]]",
		},
		{
			Assignment: "otherContent",
			Extends:    "mypage",
			Block:      "content",
			Template:   "page/mypage/templates/content/other-content.html.tmpl",
			Handler:    "otherContentHandler",
			Type:       "SubView",
		},
		{
			Type:    "Spacer",
			Comment: "[[nav]]",
		},
		{
			Assignment: "myNav",
			Extends:    "mypage",
			Block:      "nav",
			Template:   "page/mypage/templates/nav/my-nav.html.tmpl",
			Handler:    "myNavHandler",
			Type:       "DefaultSubView",
		},
	}
	wantRoutes := []generate.Route{
		{
			Method:    "GET",
			Path:      "/my-page",
			Includes:  []string{"myNav"},
			Reference: "myContent",
		},
		{
			Path:         "/my-page/form",
			Reference:    "myForm",
			FragmentOnly: true,
		},
		{
			Method:    "POST",
			Path:      "/my-page/other",
			Reference: "otherContent",
			PageOnly:  true,
		},
	}

	if err != nil {
		t.Errorf("TemlateDataFromRoutes() error = %v", err)
		return
	}

	if !reflect.DeepEqual(gotEntries, wantEntries) {
		t.Errorf("TemlateDataFromRoutes() gotEntries = %v,\n\n want %v\n\n", gotEntries, wantEntries)
	}
	if !reflect.DeepEqual(gotRoutes, wantRoutes) {
		t.Errorf("TemlateDataFromRoutes() gotRoutes = %v,\n\n want %v\n\n", gotRoutes, wantRoutes)
	}
}

func TestTemplateDataFromRoutesValidation(t *testing.T) {
	def := PageRoutes{
		EntryPoint: "/my-page",
		RouteView: RouteView{
			Ref:      "mypage",
			Template: "page/mypage/templates/mypage.html.tmpl",
			Handler:  "mypageHandler",
			Doc:      "Test page docs",
			Blocks: []TemplateBlock{
				{
					Name: "content",
					Views: []RouteView{
						{
							Ref:      "my-content",
							Template: "page/mypage/templates/content/my-content.html.tmpl",
							Handler:  "myContentHandler",
							Doc:      "The default content",
							Includes: []string{"my-nav"},
						},
						{
							Ref:      "other-content",
							Template: "page/mypage/templates/content/other-content.html.tmpl",
							Handler:  "otherContentHandler",
							Doc:      "The other content",
						},
					},
				},
			},
		},
	}

	_, _, _, _, err := TemplateDataForRoutes(def, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "no paths were found in this routemap") {
		t.Errorf("TemplateDataForRoutes() expecting to complain about zero routes, got: %s", err)
	}
}

func TestTemplateDataFromRoutes_EmptyBlock(t *testing.T) {
	def := PageRoutes{
		EntryPoint: "/my-page",
		RouteView: RouteView{
			Ref:      "mypage",
			Template: "page/mypage/templates/mypage.html.tmpl",
			Handler:  "mypageHandler",
			Path:     "/my-page",
			Doc:      "Test page docs",
			Blocks: []TemplateBlock{
				{
					Name:  "content",
					Views: nil,
				},
			},
		},
	}
	// TODO: test missing templates and handlers
	gotEntries, gotRoutes, _, _, err := TemplateDataForRoutes(def, nil, nil)
	wantEntries := []generate.Entry{
		{
			Assignment: "mypage",
			Block:      "",
			Extends:    "",
			Template:   "page/mypage/templates/mypage.html.tmpl",
			Handler:    "mypageHandler",
			Type:       "PageView",
		},
		{
			Type:    "HasSubView",
			Extends: "mypage",
			Block:   "content",
		},
	}
	wantRoutes := []generate.Route{
		{
			Path:      "/my-page",
			Reference: "mypage",
			PageOnly:  true,
		},
	}

	if err != nil {
		t.Errorf("TemlateDataFromRoutes() unexpected error = %v", err)
		return
	}

	if !reflect.DeepEqual(gotEntries, wantEntries) {
		t.Errorf("TemlateDataFromRoutes() gotEntries = %v,\n\n want %v\n\n", gotEntries, wantEntries)
	}

	if !reflect.DeepEqual(gotRoutes, wantRoutes) {
		t.Errorf("TemlateDataFromRoutes() gotRoutes = %v,\n\n want %v\n\n", gotRoutes, wantRoutes)
	}
}

func TestTemplateDataFromRoutes_DeeplyNestedBlocksDebugging(t *testing.T) {
	// Troubleshooting an issues with how block names are assigned to entries for siblinings
	// at blocks nested four levels deep specifically.
	def := PageRoutes{
		EntryPoint: "/my-page",
		RouteView: RouteView{
			Ref:      "mypage",
			Template: "page/mypage/templates/mypage.html.tmpl",
			Handler:  "mypageHandler",
			Doc:      "Test page docs",
			Blocks: []TemplateBlock{
				{
					Name: "aaaa",
					Views: []RouteView{
						{
							Ref:      "aaaa-content",
							Template: "page/mypage/templates/aaaa-content.html.tmpl",
							Handler:  "aaaaContentHandler",
							Doc:      "The aaa content",
							Blocks: []TemplateBlock{
								{
									Name: "bbbb",
									Views: []RouteView{
										{
											Ref:      "bbbb-content",
											Template: "page/mypage/templates/bbbb-content.html.tmpl",
											Handler:  "bbbbContentHandler",
											Doc:      "The bbbb content",
											Blocks: []TemplateBlock{
												{
													Name: "cccc",
													Views: []RouteView{
														{
															Ref:      "cccc-content",
															Template: "page/mypage/templates/cccc-content.html.tmpl",
															Handler:  "ccccContentHandler",
															Doc:      "The cccc content",
															Blocks: []TemplateBlock{
																{
																	Name: "dddd",
																	Views: []RouteView{
																		{
																			Ref:      "dddd-content",
																			Template: "page/mypage/templates/dddd-content.html.tmpl",
																			Handler:  "ddddContentHandler",
																			Doc:      "The dddd content",
																			Path:     "/dddd",
																		},
																	},
																},
																{
																	Name: "d2d2d2d2",
																	Views: []RouteView{
																		{
																			Ref:      "d2d2d2d2-content",
																			Template: "page/mypage/templates/d2d2d2d2-content.html.tmpl",
																			Handler:  "d2d2d2d2ContentHandler",
																			Doc:      "The d2d2d2d2 content",
																			Path:     "/d2d2d2d2",
																			Blocks: []TemplateBlock{
																				{
																					Name: "eeee",
																					Views: []RouteView{
																						{
																							Ref:      "eeee-content",
																							Template: "page/mypage/templates/eeee-content.html.tmpl",
																							Handler:  "eeeeContentHandler",
																							Doc:      "The eeee content",
																							Path:     "/eeee",
																						},
																					},
																				},
																			},
																		},
																	},
																},
																{
																	Name: "d3d3d3d3",
																	Views: []RouteView{
																		{
																			Ref:      "d3d3d3d3-content",
																			Template: "page/mypage/templates/d3d3d3d3-content.html.tmpl",
																			Handler:  "d3d3d3d3ContentHandler",
																			Doc:      "The d3d3d3d3 content",
																			Path:     "/d3d3d3d3",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	gotEntries, gotRoutes, _, _, err := TemplateDataForRoutes(def, nil, nil)
	wantEntries := []generate.Entry{
		{
			Assignment: "mypage",
			Block:      "",
			Extends:    "",
			Template:   "page/mypage/templates/mypage.html.tmpl",
			Handler:    "mypageHandler",
			Type:       "PageView",
		},
		{
			Type:    "Spacer",
			Comment: "[[aaaa]]",
		},
		{
			Type:       "SubView",
			Assignment: "aaaaContent",
			Extends:    "mypage",
			Block:      "aaaa",
			Template:   "page/mypage/templates/aaaa-content.html.tmpl",
			Handler:    "aaaaContentHandler",
		},
		{
			Type:    "Spacer",
			Comment: "[[aaaa.bbbb]]",
		},
		{
			Type:       "SubView",
			Assignment: "bbbbContent",
			Extends:    "aaaaContent",
			Block:      "bbbb",
			Template:   "page/mypage/templates/bbbb-content.html.tmpl",
			Handler:    "bbbbContentHandler",
		},
		{
			Type:    "Spacer",
			Comment: "[[aaaa.bbbb.cccc]]",
		},
		{
			Type:       "SubView",
			Assignment: "ccccContent",
			Extends:    "bbbbContent",
			Block:      "cccc",
			Template:   "page/mypage/templates/cccc-content.html.tmpl",
			Handler:    "ccccContentHandler",
		},
		{
			Type:    "Spacer",
			Comment: "[[aaaa.bbbb.cccc.dddd]]",
		},
		{
			Type:       "SubView",
			Assignment: "ddddContent",
			Extends:    "ccccContent",
			Block:      "dddd",
			Template:   "page/mypage/templates/dddd-content.html.tmpl",
			Handler:    "ddddContentHandler",
		},
		{
			Type:    "Spacer",
			Comment: "[[aaaa.bbbb.cccc.d2d2d2d2]]",
		},
		{
			Type:       "SubView",
			Assignment: "d2d2d2d2Content",
			Extends:    "ccccContent",
			Block:      "d2d2d2d2",
			Template:   "page/mypage/templates/d2d2d2d2-content.html.tmpl",
			Handler:    "d2d2d2d2ContentHandler",
		},
		{
			Type:    "Spacer",
			Comment: "[[aaaa.bbbb.cccc.d2d2d2d2.eeee]]",
		},
		{
			Type:       "SubView",
			Assignment: "eeeeContent",
			Extends:    "d2d2d2d2Content",
			Block:      "eeee",
			Template:   "page/mypage/templates/eeee-content.html.tmpl",
			Handler:    "eeeeContentHandler",
		},
		{
			Type:    "Spacer",
			Comment: "[[aaaa.bbbb.cccc.d3d3d3d3]]",
		},
		{
			Type:       "SubView",
			Assignment: "d3d3d3d3Content",
			Extends:    "ccccContent",
			Block:      "d3d3d3d3",
			Template:   "page/mypage/templates/d3d3d3d3-content.html.tmpl",
			Handler:    "d3d3d3d3ContentHandler",
		},
	}
	wantRoutes := []generate.Route{
		{
			Path:      "/dddd",
			Reference: "ddddContent",
			PageOnly:  true,
		},
		{
			Path:      "/d2d2d2d2",
			Reference: "d2d2d2d2Content",
			PageOnly:  true,
		},
		{
			Path:      "/eeee",
			Reference: "eeeeContent",
			PageOnly:  true,
		},
		{
			Path:      "/d3d3d3d3",
			Reference: "d3d3d3d3Content",
			PageOnly:  true,
		},
	}

	if err != nil {
		t.Errorf("TemlateDataFromRoutes() unexpected error = %v", err)
		return
	}

	if !reflect.DeepEqual(gotEntries, wantEntries) {
		t.Errorf("TemlateDataFromRoutes() gotEntries = %v,\n\n want %v\n\n", gotEntries, wantEntries)
	}

	if !reflect.DeepEqual(gotRoutes, wantRoutes) {
		t.Errorf("TemlateDataFromRoutes() gotRoutes = %v,\n\n want %v\n\n", gotRoutes, wantRoutes)
	}
}
