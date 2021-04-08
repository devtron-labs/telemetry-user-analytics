package telemetry

import (
	"bytes"
	"fmt"
	"github.com/devtron-labs/telemetry-user-analytics/common"
	"github.com/devtron-labs/telemetry-user-analytics/internal/sql/repository"
	"go.uber.org/zap"
	"net/http"
)

type TelemetryEventService interface {
	GetAll() ([]common.TelemetryEvent, error)
}

type TelemetryEventServiceImpl struct {
	logger                   *zap.SugaredLogger
	telemetryEventRepository repository.TelemetryEventRepository
	client                   *http.Client
}

func NewTelemetryEventServiceImpl(logger *zap.SugaredLogger, telemetryEventRepository repository.TelemetryEventRepository,
	client *http.Client) *TelemetryEventServiceImpl {
	serviceImpl := &TelemetryEventServiceImpl{
		logger:                   logger,
		telemetryEventRepository: telemetryEventRepository,
		client:                   client,
	}
	return serviceImpl
}

func (impl TelemetryEventServiceImpl) GetAll() ([]common.TelemetryEvent, error) {
	model, err := impl.telemetryEventRepository.GetAll()
	if err != nil {
		impl.logger.Errorw("error while fetching telemetry from db", "error", err)
		return nil, err
	}
	var response []common.TelemetryEvent
	for _, m := range model {
		response = append(response, common.TelemetryEvent{
			Id:   m.Id,
			UPID: m.UPID,
		})
	}
	if response == nil || len(response) == 0 {
		response = make([]common.TelemetryEvent, 0)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/triggers", ""), bytes.NewBuffer([]byte("")))
	if err != nil {
		impl.logger.Errorw("error while writing test suites", "err", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = impl.client.Do(req)
	if err != nil {
		impl.logger.Errorw("error while UpdateJiraTransition request ", "err", err)
		return nil, err
	}

	return response, nil
}
