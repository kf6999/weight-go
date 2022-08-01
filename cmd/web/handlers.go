package main

import (
	"errors"
	"fmt"
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

	app.render(w, http.StatusOK, "home.tmpl", &templateData{
		Weights: weights,
	})
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

	app.render(w, http.StatusOK, "view.tmpl", &templateData{
		Weight: weight,
	})
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
