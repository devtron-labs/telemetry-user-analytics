/*
	@author: vikram@github.com/devtron-labs
	@description: telemetry crud
*/
package repository

import (
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"time"
)

type TelemetryPlatformRepository interface {
	CreatePlatform(model *Platform) (*Platform, error)
	UpdatePlatform(model *Platform) (*Platform, error)
	GetById(id int) (*Platform, error)
	GetByUPID(upid string) (*Platform, error)
	GetAll() ([]Platform, error)
	GetConnection() (dbConnection *pg.DB)
}

type TelemetryPlatformRepositoryImpl struct {
	dbConnection *pg.DB
	Logger       *zap.SugaredLogger
}

func NewTelemetryPlatformRepositoryImpl(dbConnection *pg.DB) *TelemetryPlatformRepositoryImpl {
	return &TelemetryPlatformRepositoryImpl{dbConnection: dbConnection}
}

type Platform struct {
	TableName      struct{}  `sql:"platform"`
	Id             int32     `sql:"id,pk"`
	UPID           string    `sql:"upid,notnull"`
	CreatedOn      time.Time `sql:"created_on"`
	ModifiedOn     time.Time `sql:"modified_on"`
	ServerVersion  string    `sql:"server_version"`
	DevtronVersion string    `sql:"devtron_version"`
}

func (impl *TelemetryPlatformRepositoryImpl) CreatePlatform(model *Platform) (*Platform, error) {
	err := impl.dbConnection.Insert(model)
	if err != nil {
		return model, err
	}
	return model, nil
}
func (impl *TelemetryPlatformRepositoryImpl) UpdatePlatform(model *Platform) (*Platform, error) {
	err := impl.dbConnection.Update(model)
	if err != nil {
		return model, err
	}
	return model, nil
}
func (impl *TelemetryPlatformRepositoryImpl) GetById(id int) (*Platform, error) {
	var model Platform
	err := impl.dbConnection.Model(&model).Where("id = ?", id).Select()
	return &model, err
}
func (impl *TelemetryPlatformRepositoryImpl) GetByUPID(upid string) (*Platform, error) {
	var model Platform
	err := impl.dbConnection.Model(&model).Where("upid = ?", upid).Select()
	return &model, err
}
func (impl *TelemetryPlatformRepositoryImpl) GetAll() ([]Platform, error) {
	var model []Platform
	err := impl.dbConnection.Model(&model).Order("updated_on desc").Select()
	return model, err
}
func (impl *TelemetryPlatformRepositoryImpl) GetConnection() (dbConnection *pg.DB) {
	return impl.dbConnection
}
