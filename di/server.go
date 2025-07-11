package di

import (
	"go.uber.org/fx"
	"html/template"
	"net/http"

	"golang_test_task1/application"
	"golang_test_task1/config"
	apphttp "golang_test_task1/userinterface/http"
	"golang_test_task1/userinterface/http/handlers"
)

func ServerModule() fx.Option {
	return fx.Options(
		applicationModule(),
		templateModule(),
		formHandlerModule(),
		fx.Provide(
			config.NewServerConfig,
			apphttp.NewRouter,
		),
		fx.Invoke(
			func(cfg *config.ServerConfig, router http.Handler) error {
				return http.ListenAndServe(cfg.Address(), router)
			},
		),
	)
}

func templateModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func() *template.Template {
				return template.Must(template.ParseGlob("userinterface/http/views/*.html"))
			},
		),
	)
}

func formHandlerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(s *application.UrlScrapper) handlers.UrlScrapper {
				return s
			},
			handlers.NewFormHandler,
		),
	)
}
