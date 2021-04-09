package telemetry

import (
	"fmt"
	"github.com/devtron-labs/telemetry-user-analytics/internal/sql/repository"
	"github.com/go-pg/pg"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CronService interface {
	Process()
}

type CronServiceImpl struct {
	logger                            *zap.SugaredLogger
	telemetryPlatformRepository       repository.TelemetryPlatformRepository
	telemetryInstallHistoryRepository repository.TelemetryInstallHistoryRepository
	client                            *http.Client
}

func NewCronServiceImpl(logger *zap.SugaredLogger, telemetryPlatformRepository repository.TelemetryPlatformRepository,
	telemetryInstallHistoryRepository repository.TelemetryInstallHistoryRepository, client *http.Client) *CronServiceImpl {
	serviceImpl := &CronServiceImpl{
		logger:                            logger,
		telemetryPlatformRepository:       telemetryPlatformRepository,
		telemetryInstallHistoryRepository: telemetryInstallHistoryRepository,
		client:                            client,
	}
	gocron.Every(3).Minute().Do(serviceImpl.Process)
	<-gocron.Start()
	return serviceImpl
}

func (impl *CronServiceImpl) Process() {
	fmt.Println(">>>>>>>>>>>>> cron process", time.Now())
	modelHistory, err := impl.telemetryInstallHistoryRepository.GetById(1)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error while fetching telemetry from db", "error", err)
		return
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
			return
		}
	} else {
		modelHistory.InstallCount = modelHistory.InstallCount + 1
		modelHistory.ActivePlatform = modelHistory.ActivePlatform + 1
		modelHistory.SuccessCount = modelHistory.SuccessCount + 1
		modelHistory, err = impl.telemetryInstallHistoryRepository.UpdatePlatformHistory(modelHistory)
		if err != nil {
			impl.logger.Errorw("error while fetching telemetry from db", "error", err)
			return
		}
	}
}
