package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var pageData = []struct {
	name          string
	renderer      string
	template      string
	expectedError bool
	errorMessage  string
}{
	{"go_page", "go", "home", false, "error rendering go template"},
	{"go_page_no_template", "go", "no_template", true, "no error rendering non existent go template, when one is expected"},
	{"jet_page", "jet", "home", false, "error rendering jet template"},
	{"jet_page_no_template", "jet", "no_template", true, "no error rendering non existent jet template, when one is expected"},
	{"invalid_render_engine", "fish", "home", true, "no error rendering with an invalid template engine"},
}

func TestRender_Page(t *testing.T) {
	for _, data := range pageData {
		r, err := http.NewRequest("GET", "/some-url", nil)
		if err != nil {
			t.Error(err)
		}
		w := httptest.NewRecorder()
		testRenderer.Renderer = data.renderer
		testRenderer.RootPath = "./testData"
		err = testRenderer.Page(w, r, data.template, nil, nil)
		if data.expectedError {
			if err == nil {
				t.Errorf("%s: %s\n", data.name, data.errorMessage)
			}
		} else {
			if err != nil {
				t.Errorf("%s: %s: %s\n", data.name, data.errorMessage, err.Error())
			}
		}
	}
}

func TestRender_GoPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/any", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testData"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error Rendering Page", err)
	}
}

func TestRender_JetPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/any", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "jet"
	testRenderer.RootPath = "./testData"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error Rendering Page", err)
	}
}
