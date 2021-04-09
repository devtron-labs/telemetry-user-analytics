/*
	@author: vikram@github.com/devtron-labs
	@description: telemetry crud
*/
package repository

import (
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type TelemetryInstallHistoryRepository interface {
	CreatePlatformHistory(model *PlatformInstallHistory) (*PlatformInstallHistory, error)
	UpdatePlatformHistory(model *PlatformInstallHistory) (*PlatformInstallHistory, error)
	GetById(id int) (*PlatformInstallHistory, error)
	GetByUPID(upid string) (*PlatformInstallHistory, error)
	GetAll() ([]PlatformInstallHistory, error)
	GetConnection() (dbConnection *pg.DB)
}

type TelemetryInstallHistoryRepositoryImpl struct {
	dbConnection *pg.DB
	Logger       *zap.SugaredLogger
}

func NewTelemetryInstallHistoryRepositoryImpl(dbConnection *pg.DB) *TelemetryInstallHistoryRepositoryImpl {
	return &TelemetryInstallHistoryRepositoryImpl{dbConnection: dbConnection}
}

type PlatformInstallHistory struct {
	TableName      struct{} `sql:"platform_install_history"`
	Id             int32    `sql:"id,pk"`
	InstallCount   int      `sql:"install_count"`
	FailCount      int      `sql:"fail_count"`
	SuccessCount   int      `sql:"success_count"`
	ActivePlatform int      `sql:"active_platform"`
}

func (impl *TelemetryInstallHistoryRepositoryImpl) CreatePlatformHistory(model *PlatformInstallHistory) (*PlatformInstallHistory, error) {
	err := impl.dbConnection.Insert(model)
	if err != nil {
		return model, err
	}
	return model, nil
}
func (impl *TelemetryInstallHistoryRepositoryImpl) UpdatePlatformHistory(model *PlatformInstallHistory) (*PlatformInstallHistory, error) {
	err := impl.dbConnection.Update(model)
	if err != nil {
		return model, err
	}
	return model, nil
}
func (impl *TelemetryInstallHistoryRepositoryImpl) GetById(id int) (*PlatformInstallHistory, error) {
	var model PlatformInstallHistory
	err := impl.dbConnection.Model(&model).Where("id = ?", id).Select()
	return &model, err
}
func (impl *TelemetryInstallHistoryRepositoryImpl) GetByUPID(upid string) (*PlatformInstallHistory, error) {
	var model PlatformInstallHistory
	err := impl.dbConnection.Model(&model).Where("upid = ?", upid).Select()
	return &model, err
}
func (impl *TelemetryInstallHistoryRepositoryImpl) GetAll() ([]PlatformInstallHistory, error) {
	var model []PlatformInstallHistory
	err := impl.dbConnection.Model(&model).Order("updated_on desc").Select()
	return model, err
}
func (impl *TelemetryInstallHistoryRepositoryImpl) GetConnection() (dbConnection *pg.DB) {
	return impl.dbConnection
}
