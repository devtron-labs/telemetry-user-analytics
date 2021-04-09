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
	gocron.Every(5).Minute().Do(serviceImpl.Process)
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
		//modelHistory.InstallCount = 1
		//modelHistory.SuccessCount = 1
		//modelHistory.FailCount = 0
		//modelHistory.ActivePlatform = 1
		modelHistory, err = impl.telemetryInstallHistoryRepository.CreatePlatformHistory(modelHistory)
		if err != nil {
			impl.logger.Errorw("error while fetching telemetry from db", "error", err)
			return
		}
	} else {
		installCount := 0
		failureCount := 0
		activeCount := 0

		platforms, err := impl.telemetryPlatformRepository.GetAll()
		if err != nil && err != pg.ErrNoRows {
			impl.logger.Errorw("error while fetching telemetry from db", "error", err)
			return
		}
		installCount = len(platforms)
		now := time.Now()
		dateGap := now.Add(-30 * time.Minute)
		for _, platform := range platforms {
			if platform.CreatedOn.Before(dateGap) && platform.ModifiedOn.IsZero() {
				failureCount = failureCount + 1
			}

			if platform.ModifiedOn.After(dateGap) && platform.ModifiedOn.IsZero() {
				activeCount = activeCount + 1
			}
		}

		// find out total platform, set into install count
		// find out total platform which are created by last 1 hour but not having modified date, to set into failure count
		// find out total platform which are modified date, set into success count
		// find out total platform which are modified since last 1 day, set into active account

		modelHistory.InstallCount = installCount
		modelHistory.ActivePlatform = activeCount
		modelHistory.FailCount = failureCount
		modelHistory.SuccessCount = installCount - failureCount
		modelHistory, err = impl.telemetryInstallHistoryRepository.UpdatePlatformHistory(modelHistory)
		if err != nil {
			impl.logger.Errorw("error while fetching telemetry from db", "error", err)
			return
		}
	}
}
