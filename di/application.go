package di

import (
	"go.uber.org/fx"
	"net/http"

	"golang_test_task1/application"
)

func applicationModule() fx.Option {
	return fx.Options(
		httpClientModule(),
		urlScrapperModule(),
	)
}

func httpClientModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func() *http.Client {
				return http.DefaultClient
			},
		),
	)
}

func urlScrapperModule() fx.Option {
	return fx.Options(
		fx.Provide(
			application.NewUrlScrapper,
		),
	)
}
