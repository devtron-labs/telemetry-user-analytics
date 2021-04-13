package telemetry

import (
	"github.com/devtron-labs/telemetry-user-analytics/common"
	"github.com/devtron-labs/telemetry-user-analytics/internal/sql/repository"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"net/http"
)

type TelemetryEventService interface {
	CreatePlatform(dto *common.TelemetryUserAnalyticsDto) (*common.TelemetryUserAnalyticsDto, error)
	GetByUPID(upid string) (*common.TelemetryUserAnalyticsDto, error)
	GetAll() ([]*common.TelemetryUserAnalyticsDto, error)
}

type TelemetryEventServiceImpl struct {
	logger                            *zap.SugaredLogger
	telemetryPlatformRepository       repository.TelemetryPlatformRepository
	telemetryInstallHistoryRepository repository.TelemetryInstallHistoryRepository
	client                            *http.Client
}

func NewTelemetryEventServiceImpl(logger *zap.SugaredLogger, telemetryPlatformRepository repository.TelemetryPlatformRepository,
	telemetryInstallHistoryRepository repository.TelemetryInstallHistoryRepository, client *http.Client) *TelemetryEventServiceImpl {
	serviceImpl := &TelemetryEventServiceImpl{
		logger:                            logger,
		telemetryPlatformRepository:       telemetryPlatformRepository,
		telemetryInstallHistoryRepository: telemetryInstallHistoryRepository,
		client:                            client,
	}
	return serviceImpl
}

func (impl *TelemetryEventServiceImpl) CreatePlatform(dto *common.TelemetryUserAnalyticsDto) (*common.TelemetryUserAnalyticsDto, error) {
	/*
		model, err := impl.telemetryPlatformRepository.GetByUPID(dto.UPID)
		if err != nil && err != pg.ErrNoRows {
			impl.logger.Errorw("error while fetching telemetry from db", "error", err)
			return nil, err
		}
		if err == pg.ErrNoRows {
	*/
	model := &repository.Platform{}
	model.UPID = dto.UPID
	model.DevtronVersion = dto.DevtronVersion
	model.ServerVersion = dto.ServerVersion
	model.CreatedOn = dto.Timestamp
	model.Clusters = dto.Clusters
	model.Environments = dto.Environments
	model.NoOfProdApps = dto.NoOfProdApps
	model.NoOfNonProdApps = dto.NoOfNonProdApps
	model.Users = dto.Users
	model.EventType = dto.EventType
	model, err := impl.telemetryPlatformRepository.CreatePlatform(model)
	if err != nil {
		impl.logger.Errorw("error while fetching telemetry from db", "error", err)
		return nil, err
	}
	dto.Id = model.Id

	// total install count counter logic
	_, err = impl.telemetryInstallHistoryRepository.GetById(1)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error while fetching telemetry from db", "error", err)
		return nil, err
	}
	if err == pg.ErrNoRows {
		modelHistory := &repository.PlatformInstallHistory{}
		modelHistory.InstallCount = 1
		modelHistory.SuccessCount = 1
		modelHistory.FailCount = 0
		modelHistory.ActivePlatform = 1
		modelHistory, err = impl.telemetryInstallHistoryRepository.CreatePlatformHistory(modelHistory)
		if err != nil {
			impl.logger.Errorw("error while fetching telemetry from db", "error", err)
			return nil, err
		}
	} else {
		// todo-  call cron service
	}
	/*
		} else {
			model.DevtronVersion = dto.DevtronVersion
			model.ServerVersion = dto.ServerVersion
			model.ModifiedOn = dto.Timestamp
			model, err := impl.telemetryPlatformRepository.UpdatePlatform(model)
			if err != nil {
				impl.logger.Errorw("error while fetching telemetry from db", "error", err)
				return nil, err
			}
			dto.Id = model.Id
		}
	*/
	dto.Id = model.Id
	return dto, nil
}

func (impl *TelemetryEventServiceImpl) GetByUPID(upid string) (*common.TelemetryUserAnalyticsDto, error) {
	data := &common.TelemetryUserAnalyticsDto{}
	model, err := impl.telemetryPlatformRepository.GetByUPID(upid)
	if err != nil {
		impl.logger.Errorw("error while fetching telemetry from db", "error", err)
		return nil, err
	}
	data.UPID = model.UPID
	data.ServerVersion = model.ServerVersion
	data.DevtronVersion = model.DevtronVersion
	return data, nil
}

func (impl *TelemetryEventServiceImpl) GetAll() ([]*common.TelemetryUserAnalyticsDto, error) {
	model, err := impl.telemetryPlatformRepository.GetAll()
	if err != nil {
		impl.logger.Errorw("error while fetching telemetry from db", "error", err)
		return nil, err
	}
	var response []*common.TelemetryUserAnalyticsDto
	for _, m := range model {
		response = append(response, &common.TelemetryUserAnalyticsDto{
			Id:   m.Id,
			UPID: m.UPID,
		})
	}
	if response == nil || len(response) == 0 {
		response = make([]*common.TelemetryUserAnalyticsDto, 0)
	}

	return response, nil
}
