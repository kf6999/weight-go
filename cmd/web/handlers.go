package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"weight.kenfan.org/internal/models"
	"weight.kenfan.org/internal/validator"
)

type weightCreateForm struct {
	Weight              string `form:"weight"`
	Notes               string `form:"notes"`
	validator.Validator `form:"-"`
}

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

func (app *application) weightCreatePost(w http.ResponseWriter, r *http.Request) {
	var form weightCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Weight), "weight", "This field cannot be blank")
	form.CheckField(validator.IsInt(form.Weight), "weight", "This field must be a number")
	form.CheckField(validator.MaxCharacters(form.Notes, 300), "notes", "This field cannot be more than 300 characters long")

	if !form.Valid() {
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
