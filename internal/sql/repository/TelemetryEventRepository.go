/*
	@author: vikram@github.com/devtron-labs
	@description: telemetry crud
*/
package repository

import (
	"fmt"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"time"
)

type TelemetryEventRepository interface {
	CreatePlatform(model *Platform, tx *pg.Tx) (*Platform, error)
	UpdatePlatform(model *Platform, tx *pg.Tx) (*Platform, error)
	GetById(id int32) (*Platform, error)
	GetAll() ([]Platform, error)
	GetConnection() (dbConnection *pg.DB)
}

type TelemetryEventRepositoryImpl struct {
	dbConnection *pg.DB
	Logger       *zap.SugaredLogger
}

func NewTelemetryEventRepositoryImpl(dbConnection *pg.DB) *TelemetryEventRepositoryImpl {
	return &TelemetryEventRepositoryImpl{dbConnection: dbConnection}
}

type Platform struct {
	TableName  struct{}  `sql:"platform"`
	Id         int32     `sql:"id,pk"`
	UPID       string    `sql:"upid,notnull"`
	CreatedOn  time.Time `sql:"created_on"`
	ModifiedOn time.Time `sql:"modified_on"`
}

type PlatformInstallHistory struct {
	TableName      struct{} `sql:"platform_install_history"`
	Id             int32    `sql:"id,pk"`
	InstallCount   int      `sql:"install_count"`
	FailCount      int      `sql:"fail_count"`
	SuccessCount   int      `sql:"success_count"`
	ActivePlatform int      `sql:"active_platform"`
}

func (impl TelemetryEventRepositoryImpl) CreatePlatform(model *Platform, tx *pg.Tx) (*Platform, error) {
	err := tx.Insert(model)
	if err != nil {
		fmt.Println("Exception;", err)
		return model, err
	}
	//TODO - Create Entry In UserRole With Default Role for User
	return model, nil
}
func (impl TelemetryEventRepositoryImpl) UpdatePlatform(model *Platform, tx *pg.Tx) (*Platform, error) {
	err := tx.Update(model)
	if err != nil {
		fmt.Println("Exception;", err)
		return model, err
	}

	//TODO - Create Entry In UserRole With Default Role for User

	return model, nil
}
func (impl TelemetryEventRepositoryImpl) GetById(id int32) (*Platform, error) {
	var model Platform
	err := impl.dbConnection.Model(&model).Where("id = ?", id).Where("active = ?", true).Select()
	return &model, err
}
func (impl TelemetryEventRepositoryImpl) GetAll() ([]Platform, error) {
	var model []Platform
	err := impl.dbConnection.Model(&model).Where("active = ?", true).Order("updated_on desc").Select()
	return model, err
}
func (impl *TelemetryEventRepositoryImpl) GetConnection() (dbConnection *pg.DB) {
	return impl.dbConnection
}
