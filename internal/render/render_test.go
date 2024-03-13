package render

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jdonahue135/golf-league-app/internal/models"
)

func TestAdd(t *testing.T) {
	res := Add(1, 2)
	if res != 3 {
		t.Error(fmt.Sprintf("Expected 3 but got %d", res))
	}
}

func TestIterate(t *testing.T) {
	res := Iterate(5)
	if len(res) != 5 {
		t.Error(fmt.Sprintf("Expected slice of length 5 but got %d", len(res)))
	}
}

func TestHumanDate(t *testing.T) {
	date := time.Date(1991, time.February, 21, 0, 0, 0, 0, time.UTC)
	res := HumanDate(date)
	if res != "1991-02-21" {
		t.Error(fmt.Sprintf("Expected 1991-02-21 but got %s", res))
	}
}

func TestFormatDate(t *testing.T) {
	date := time.Date(1991, time.February, 21, 0, 0, 0, 0, time.UTC)
	format := "2006-01-02"
	res := FormatDate(date, format)
	if res != "1991-02-21" {
		t.Error(fmt.Sprintf("Expected 1991-02-21 but got %s", res))
	}
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
