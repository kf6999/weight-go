package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"weight.kenfan.org/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	weights, err := app.weights.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, weight := range weights {
		fmt.Fprintf(w, "%+v\n", weight)
	}

	//files := []string{
	//	"./ui/html/base.tmpl",
	//	"./ui/html/pages/home.tmpl",
	//	"./ui/html/partials/nav.tmpl",
	//}
	//
	//// Use template.ParseFiles to read template file into template set
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//	http.Error(w, "Internal Server Error", 500)
	//	return
	//}
	//
	//// Use Execute() to write the content of the "base" template as the response body.
	//// Last parameter to Execute() represents dynamic data we want to pass in
	//err = ts.ExecuteTemplate(w, "base", nil)
	//if err != nil {
	//	app.serverError(w, err)
	//}
}

func (app *application) weightView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	weight, err := app.weights.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Weight: weight,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}
func (app *application) weightCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	weight := 100
	note := "test"

	id, err := app.weights.Insert(weight, note)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/weight/view?id=%d", id), http.StatusSeeOther)
}
