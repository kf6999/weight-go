package main

import (
	"html/template"
	"path/filepath"
	"time"
	"weight.kenfan.org/internal/models"
)

type templateData struct {
	Weight      *models.Weight
	Weights     []*models.Weight
	CurrentYear int
}

func humanDate(t time.Time) string {
	return t.Format("Jan 02 2006 at 15:04")
}

func normDate(t time.Time) string {
	return t.Format("01-02-2006")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"normDate":  normDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the  page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map as normal...
		cache[name] = ts
	}

	return cache, nil
}
