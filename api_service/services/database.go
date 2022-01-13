package services

import (
	"db_tp/db"
	"db_tp/models"
)

var DatabaseSrv *DatabaseService

type DatabaseService struct {
	db *db.PostgresDbEngine
}

func NewDatabaseService(db *db.PostgresDbEngine) *DatabaseService {
	srv := new(DatabaseService)
	srv.db = db
	return srv
}

func (forum *DatabaseService) GetStatus() *models.Status {

}

func (forum *DatabaseService) ClearDatabase() error {

}
