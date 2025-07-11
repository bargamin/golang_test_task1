package handlers

type UrlForm struct {
	Url string `validate:"required,url"`
}
