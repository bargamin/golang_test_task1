package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"golang_test_task1/di"
)

func serverStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := fx.New(di.ServerModule())

			ctx := context.Background()

			if err := app.Start(ctx); err != nil {
				return err
			}

			defer func() { _ = app.Stop(ctx) }()

			return nil
		},
	}
}
