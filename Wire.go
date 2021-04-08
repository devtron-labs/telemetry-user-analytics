//+build wireinject

package main

import (
	"github.com/devtron-labs/telemetry-user-analytics/api"
	"github.com/devtron-labs/telemetry-user-analytics/client"
	"github.com/devtron-labs/telemetry-user-analytics/internal/logger"
	"github.com/devtron-labs/telemetry-user-analytics/internal/sql"
	"github.com/devtron-labs/telemetry-user-analytics/internal/sql/repository"
	"github.com/devtron-labs/telemetry-user-analytics/pkg/telemetry"
	"github.com/devtron-labs/telemetry-user-analytics/pubsub"
	"github.com/google/wire"
)

func InitializeApp() (*App, error) {
	wire.Build(
		NewApp,
		api.NewMuxRouter,
		logger.NewSugardLogger,
		logger.NewHttpClient,
		sql.GetConfig,
		sql.NewDbConnection,
		api.NewRestHandlerImpl,
		wire.Bind(new(api.RestHandler), new(*api.RestHandlerImpl)),
		client.NewPubSubClient,
		pubsub.NewNatSubscription,

		telemetry.NewTelemetryEventServiceImpl,
		wire.Bind(new(telemetry.TelemetryEventService), new(*telemetry.TelemetryEventServiceImpl)),
		repository.NewTelemetryEventRepositoryImpl,
		wire.Bind(new(repository.TelemetryEventRepository), new(*repository.TelemetryEventRepositoryImpl)),
	)
	return &App{}, nil
}
