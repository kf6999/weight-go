package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"weight.kenfan.org/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	weights, err := app.weights.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Weights = weights

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) weightView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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
	data := app.newTemplateData(r)
	data.Weight = weight

	app.render(w, http.StatusOK, "view.tmpl", data)
}
func (app *application) weightCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) weightCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	weight := r.PostForm.Get("weight")
	notes := r.PostForm.Get("notes")

	id, err := app.weights.Insert(weight, notes)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/weight/view/%d", id), http.StatusSeeOther)
}
