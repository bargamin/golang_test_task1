package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"html/template"
	"net/http"

	"golang_test_task1/application"
	"golang_test_task1/domain"
)

const layoutHtml = "layout.html"

type UrlScrapper interface {
	GetInfoByURL(ctx context.Context, rawURL string) (*domain.WebsiteInfo, error)
}

type FormHandler struct {
	scrapper  UrlScrapper
	templates *template.Template
	validate  *validator.Validate
}

func NewFormHandler(s UrlScrapper, templates *template.Template) *FormHandler {
	return &FormHandler{
		scrapper:  s,
		templates: templates,
		validate:  validator.New(),
	}
}

func (h *FormHandler) Form(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, layoutHtml, map[string]any{
		"Form": UrlForm{},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *FormHandler) ProcessForm(w http.ResponseWriter, r *http.Request) {
	urlForm := UrlForm{
		Url: r.FormValue("url"),
	}

	pageData := map[string]any{
		"Form": urlForm,
	}

	if validationErr := h.validate.Struct(urlForm); validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		pageData["Flash"] = validationErr.Error()
	} else if info, err := h.scrapper.GetInfoByURL(r.Context(), urlForm.Url); err != nil {
		httpErr := &application.HTTPError{}
		if errors.As(err, &httpErr) {
			pageData["Flash"] = fmt.Sprintf("HTTP error: %s with http code %d", httpErr.Description, httpErr.Status)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		pageData["Info"] = *info
	}

	if err := h.templates.ExecuteTemplate(w, layoutHtml, pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
