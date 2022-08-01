package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
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
	data.Form = weightCreateForm{}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

type weightCreateForm struct {
	Weight      string
	Notes       string
	FieldErrors map[string]string
}

func (app *application) weightCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := weightCreateForm{
		Weight:      r.PostForm.Get("weight"),
		Notes:       r.PostForm.Get("notes"),
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Weight) == "" {
		form.FieldErrors["weight"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Weight) > 100 {
		form.FieldErrors["weight"] = "This field cannot be more than 100 characters long"
	} else if _, err := strconv.Atoi(form.Weight); err != nil {
		form.FieldErrors["weight"] = "Enter number"
	}

	if utf8.RuneCountInString(form.Notes) > 300 {
		form.FieldErrors["notes"] = "This field cannot be more than 300 characters long"
	}
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.weights.Insert(form.Weight, form.Notes)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/weight/view/%d", id), http.StatusSeeOther)
}
